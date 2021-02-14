package metrics

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/tsdb"

	"github.com/janoszen/openshiftci_inspector/job"
	"github.com/janoszen/openshiftci_inspector/preprocessor"
	"github.com/janoszen/openshiftci_inspector/widget"
)

type Axis struct {
	Label string
	Min float64 `json:"min" yaml:"min"`
	Max float64 `json:"max" yaml:"max"`
}

type Query struct {
	Query string `json:"query" yaml:"query"`
	X Axis `json:"x"`
	YLabel string `json:"ylabel"`
	LegendFormat string `json:"legend" yaml:"legend"`
}

func NewMetricsPlugin(
	tmpDir string,
	queries map[string]Query,
) (preprocessor.Plugin, error) {
	return &metricsPlugin{
		tmpDir: tmpDir,
		queries: queries,
	}, nil
}

type metricsPlugin struct {
	queries map[string]Query
	tmpDir  string
}

func (m *metricsPlugin) GetArtifacts(job job.Job) []string {
	return []string{
		"/*/prometheus.tar",
	}
}

func (m *metricsPlugin) Preprocess(job job.Job, artifacts map[string]io.Reader) (map[string]widget.Widget, error) {
	if len(artifacts) > 1 {
		return nil, fmt.Errorf("more than one prometheus.tar supplied")
	}
	for artifactName, artifactReader := range artifacts {
		if job.StartTime != nil {
			endTime := time.Now()
			if job.CompletionTime != nil {
				endTime = *job.CompletionTime
			}
			return m.run(artifactName, artifactReader, *job.StartTime, endTime)
		}
	}
	return nil, fmt.Errorf("no prometheus.tar supplied")
}

func (m *metricsPlugin) run(
	artifactName string,
	artifactReader io.Reader,
	startTime time.Time,
	endTime time.Time,
) (map[string]widget.Widget, error) {
	artifactPath := m.getUnpackPath(artifactName)
	err := m.unpack(artifactPath, artifactName, artifactReader)
	defer m.cleanup(artifactPath, artifactName)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack artifact %s (%w)", artifactName, err)
	}

	registry := prometheus.NewRegistry()
	promDB, err := tsdb.Open(
		artifactPath,
		nil,
		registry,
		&tsdb.Options{
			NoLockfile: true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start Prometheus on directory %s (%w)", artifactPath, err)
	}
	defer func() {
		_ = promDB.Close()
	}()
	engine := promql.NewEngine(promql.EngineOpts{
		Logger:     nil,
		Reg:        registry,
		Timeout:    time.Minute,
		MaxSamples: 50000000,
	})
	duration, err := time.ParseDuration("5s")
	if err != nil {
		return nil, fmt.Errorf("failed to parse time duration (%w)", err)
	}

	widgets := map[string]widget.Widget{}
	for queryName, query := range m.queries {
		w, err := m.runQuery(engine, promDB, query, startTime, endTime, duration)
		if err != nil {
			return nil, err
		}
		widgets[queryName] = w
	}
	return widgets, nil
}

func (m *metricsPlugin) runQuery(
	engine *promql.Engine,
	promDB *tsdb.DB,
	query Query,
	startTime time.Time,
	endTime time.Time,
	duration time.Duration,
) (widget.Widget, error) {
	queryResult, err := engine.NewRangeQuery(
		promDB,
		query.Query,
		startTime,
		endTime,
		duration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %s (%w)", query, err)
	}
	defer func() {
		queryResult.Close()
	}()
	res := queryResult.Exec(context.TODO())
	if res.Err != nil {
		return nil, fmt.Errorf("failed to execute query: %s (%w)", query, res.Err)
	}

	if v, err := res.Vector(); err == nil {
		return transformVector(v, query)
	}
	if m, err := res.Matrix(); err == nil {
		return transformMatrix(m, query)
	}
	if s, err := res.Scalar(); err == nil {
		return transformScalar(s)
	}
	return nil, fmt.Errorf("no valid query result")
}

func (m *metricsPlugin) getUnpackPath(artifactName string) string {
	hash := sha256.Sum256([]byte(artifactName))
	return path.Join(m.tmpDir, hex.EncodeToString(hash[0:]))
}

func (m *metricsPlugin) unpack(artifactPath string, artifactName string, artifactReader io.Reader) error {
	if err := os.MkdirAll(m.tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create temporary directory %s (%w)", artifactPath, err)
	}

	gzipReader, err := gzip.NewReader(artifactReader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader for %s (%w)", artifactName, err)
	}
	tarReader := tar.NewReader(gzipReader)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to extract tar file %s (%w)", artifactName, err)
		}
		fullPath := path.Join(artifactPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			stat, err := os.Stat(fullPath)
			if err != nil {
				if err := os.Mkdir(fullPath, 0755); err != nil {
					return fmt.Errorf("failed to create directory %s (%w)", fullPath, err)
				}
			} else {
				if !stat.IsDir() {
					return fmt.Errorf("unpack target already exists and is not a directory %s (%w)", fullPath, err)
				}
			}
		case tar.TypeReg:
			outFile, err := os.Create(fullPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s (%v)", fullPath, err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				_ = outFile.Close()
				return fmt.Errorf("failed to copy to file %s (%v)", fullPath, err)
			}
			_ = outFile.Close()
			if err = os.Chtimes(fullPath, header.AccessTime, header.ChangeTime); err != nil {
				return fmt.Errorf("failed to change file times (%w)", err)
			}
			if err = os.Chmod(fullPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to change file mode (%w)", err)
			}
		case tar.TypeBlock:
			return fmt.Errorf("refusing to create block device %s", fullPath)
		case tar.TypeChar:
			return fmt.Errorf("refusing to create character device %s", fullPath)
		case tar.TypeSymlink:
			return fmt.Errorf("refusing to create symlink %s", fullPath)
		case tar.TypeFifo:
			return fmt.Errorf("refusing to create FIFO %s", fullPath)
		default:
			return fmt.Errorf(
				"unsupported type flag: %d in %s",
				header.Typeflag,
				fullPath,
			)
		}
	}

	return nil
}

func (m *metricsPlugin) cleanup(artifactPath string, name string) {
	_ = os.RemoveAll(artifactPath)
}

