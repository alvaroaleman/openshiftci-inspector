package s3

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/janoszen/openshiftci-inspector/asset"
)

type s3AssetStorage struct {
	s3     *s3.S3
	bucket string
	logger *log.Logger
}

func (s *s3AssetStorage) Shutdown(ctx context.Context) {
}

func (s *s3AssetStorage) Store(asset asset.Asset, mime string, data []byte) error {
	key := "/" + asset.JobID + "/" + asset.AssetName
	s.logger.Printf("Storing asset %s...\n", key)
	_, err := s.s3.PutObject(
		&s3.PutObjectInput{
			ACL:           aws.String(s3.BucketCannedACLPublicRead),
			Body:          bytes.NewReader(data),
			Bucket:        aws.String(s.bucket),
			ContentLength: aws.Int64(int64(len(data))),
			ContentType:   aws.String(mime),
			Key:           aws.String(key),
		},
	)
	if err != nil {
		s.logger.Printf("failed to store asset %s (%v).\n", key, err)
	} else {
		s.logger.Printf("Stored asset %s.\n", key)
	}
	return err
}

func (s *s3AssetStorage) Fetch(asset asset.Asset) (data []byte, err error) {
	key := "/" + asset.JobID + "/" + asset.AssetName
	s.logger.Printf("Fetching asset %s...\n", key)
	get, err := s.s3.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return
	}
	defer func() {
		_ = get.Body.Close()
	}()
	data, err = ioutil.ReadAll(get.Body)
	if err != nil {
		return
	}
	return
}
