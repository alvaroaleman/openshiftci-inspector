package metrics

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/janoszen/openshiftci-inspector/asset"
	"github.com/janoszen/openshiftci-inspector/asset/indexstorage"
	"github.com/janoszen/openshiftci-inspector/asset/storage"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/tsdb"
)

func NewQuery(
	index indexstorage.AssetIndex,
	storage storage.AssetStorage,
) QueryBackend {
	return &queryImpl{
		index:   index,
		storage: storage,
	}
}

type queryImpl struct {
	index   indexstorage.AssetIndex
	storage storage.AssetStorage
}

func (q *queryImpl) Query(
	ctx context.Context,
	jobID string,
	name string,
	query string,
	startTime time.Time,
	endTime time.Time,
) (QueryResponse, error) {
	tarContents, err := q.download(jobID, name)
	if err != nil {
		return QueryResponse{}, err
	}

	tarDir, err := q.extractTar(tarContents, jobID, name)
	defer func() {
		_ = os.RemoveAll(tarDir)
	}()
	if err != nil {
		return QueryResponse{}, err
	}

	result, err := q.runQuery(ctx, tarDir, query, startTime, endTime)
	if err != nil {
		return QueryResponse{}, err
	}
	return result, nil
}

func (q *queryImpl) extractTar(tarContents []byte, jobID string, name string) (string, error) {
	tmpDir := path.Join(os.TempDir(), jobID, name)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return tmpDir, fmt.Errorf("failed to create temporary directory %s (%w)", tmpDir, err)
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(tarContents))
	if err != nil {
		return tmpDir, err
	}
	tarReader := tar.NewReader(gzipReader)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return tmpDir, fmt.Errorf("failed to extract tar file (%w)", err)
		}

		fullPath := path.Join(tmpDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			stat, err := os.Stat(fullPath)
			if err != nil {
				if err := os.Mkdir(fullPath, 0755); err != nil {
					return tmpDir, fmt.Errorf("failed to create directory %s (%w)", fullPath, err)
				}
			} else {
				if !stat.IsDir() {
					return tmpDir, fmt.Errorf("unpack target already exists and is not a directory %s (%w)", fullPath, err)
				}
			}
		case tar.TypeReg:
			outFile, err := os.Create(fullPath)
			if err != nil {
				return tmpDir, fmt.Errorf("failed to create file %s (%v)", fullPath, err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				_ = outFile.Close()
				return tmpDir, fmt.Errorf("failed to copy to file %s (%v)", fullPath, err)
			}
			_ = outFile.Close()
			if err = os.Chtimes(fullPath, header.AccessTime, header.ChangeTime); err != nil {
				return tmpDir, fmt.Errorf("failed to change file times (%w)", err)
			}
			if err = os.Chmod(fullPath, os.FileMode(header.Mode)); err != nil {
				return tmpDir, fmt.Errorf("failed to change file mode (%w)", err)
			}
		case tar.TypeBlock:
			return tmpDir, fmt.Errorf("refusing to create block device %s", fullPath)
		case tar.TypeChar:
			return tmpDir, fmt.Errorf("refusing to create character device %s", fullPath)
		case tar.TypeSymlink:
			return tmpDir, fmt.Errorf("refusing to create symlink %s", fullPath)
		case tar.TypeFifo:
			return tmpDir, fmt.Errorf("refusing to create FIFO %s", fullPath)
		default:
			return tmpDir, fmt.Errorf(
				"unsupported type flag: %d in %s",
				header.Typeflag,
				fullPath,
			)
		}
	}
	return tmpDir, nil
}

func (q *queryImpl) download(jobID string, name string) ([]byte, error) {
	hasAsset, err := q.index.HasAsset(jobID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to runQuery asset index (%w)", err)
	}
	if !hasAsset {
		return nil, fmt.Errorf("job %s has no asset %s", jobID, name)
	}
	data, err := q.storage.Fetch(asset.Asset{
		JobID:     jobID,
		AssetName: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download asset %s for job %s (%w)", name, jobID, err)
	}
	return data, nil
}

func (q *queryImpl) runQuery(ctx context.Context, dir string, query string, startTime time.Time, endTime time.Time) (QueryResponse, error) {
	registry := prometheus.NewRegistry()
	promDB, err := tsdb.Open(
		dir,
		nil,
		registry,
		&tsdb.Options{
			NoLockfile: true,
		},
	)
	if err != nil {
		return QueryResponse{}, err
	}
	engine := promql.NewEngine(promql.EngineOpts{
		Logger:     nil,
		Reg:        nil,
		Timeout:    time.Minute,
		MaxSamples: 50000000,
	})
	duration, err := time.ParseDuration("5s")
	if err != nil {
		return QueryResponse{}, fmt.Errorf("failed to parse time duration (%w)", err)
	}
	queryResult, err := engine.NewRangeQuery(
		promDB,
		query,
		startTime,
		endTime,
		duration,
	)
	if err != nil {
		return QueryResponse{}, fmt.Errorf("failed to parse runQuery (%w)", err)
	}
	defer func() {
		queryResult.Close()
	}()
	res := queryResult.Exec(ctx)
	if res.Err != nil {
		return QueryResponse{}, fmt.Errorf("failed to execute query (%w)", res.Err)
	}

	response := QueryResponse{}
	if v, err := res.Vector(); err == nil {
		response.Vector = transformVector(v)
	}
	if m, err := res.Matrix(); err == nil {
		response.Matrix = transformMatrix(m)
	}
	if s, err := res.Scalar(); err == nil {
		response.Scalar = transformScalar(s)
	}
	return response, nil
}

func transformScalar(s promql.Scalar) QueryPoint {
	return QueryPoint{
		Timestamp: s.T,
		Value:     s.V,
	}
}

func transformMatrix(m promql.Matrix) []QuerySeries {
	var result []QuerySeries
	for _, s := range m {
		series := QuerySeries{
			Labels: transformLabels(s.Metric),
			Points: transformPoints(s.Points),
		}
		result = append(result, series)
	}
	return result
}

func transformPoints(points []promql.Point) []QueryPoint {
	var result []QueryPoint
	for _, p := range points {
		result = append(result, transformPoint(p))
	}
	return result
}

func transformVector(v promql.Vector) []QuerySample {
	var result []QuerySample
	for _, e := range v {
		result = append(result, QuerySample{
			Labels: transformLabels(e.Metric),
			Point:  transformPoint(e.Point),
		})
	}
	return result
}

func transformPoint(point promql.Point) QueryPoint {
	return QueryPoint{
		Timestamp: point.T,
		Value:     point.V,
	}
}

func transformLabels(metric labels.Labels) []QueryLabel {
	var result []QueryLabel
	for _, m := range metric {
		result = append(result, QueryLabel{
			Name:  m.Name,
			Value: m.Value,
		})
	}
	return result
}
