package s3

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/janoszen/openshiftci-inspector/asset/storage"
)

// NewS3AssetStorage creates an asset storage that stores assets on an S3-compatible object storage.
func NewS3AssetStorage(config S3AssetStorageConfig, logger *log.Logger) (storage.AssetStorage, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	awsConfig := createAWSConfig(config)
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	s3connection := s3.New(sess)

	assetStorage, err := ensureBucket(config, s3connection)
	if err != nil {
		return assetStorage, err
	}

	return &s3AssetStorage{
		s3:     s3connection,
		bucket: config.Bucket,
		logger: logger,
	}, nil
}

func createAWSConfig(config S3AssetStorageConfig) *aws.Config {
	var endpoint *string
	if config.Endpoint != "" {
		endpoint = aws.String(config.Endpoint)
	}
	awsConfig := &aws.Config{
		Credentials: credentials.NewCredentials(
			&credentials.StaticProvider{
				Value: credentials.Value{
					AccessKeyID:     config.AccessKey,
					SecretAccessKey: config.SecretKey,

					SessionToken: "",
					ProviderName: "",
				},
			},
		),
		Endpoint:         endpoint,
		Region:           aws.String(config.Region),
		S3ForcePathStyle: aws.Bool(config.ForcePathStyleAccess),
	}
	return awsConfig
}

func ensureBucket(config S3AssetStorageConfig, s3connection *s3.S3) (storage.AssetStorage, error) {
	bucketLocation, err := s3connection.GetBucketLocation(
		&s3.GetBucketLocationInput{
			Bucket: aws.String(config.Bucket),
		},
	)
	if err == nil {
		if bucketLocation.LocationConstraint == nil {
			if config.Region != "us-east-1" {
				return nil, fmt.Errorf(
					"bucket %s is in us-east-1 but the region configuration specifies %s",
					config.Bucket,
					config.Region,
				)
			}
		} else if *bucketLocation.LocationConstraint != config.Region {
			return nil, fmt.Errorf(
				"bucket %s is in %s but the region configuration specifies %s",
				config.Bucket,
				*bucketLocation.LocationConstraint,
				config.Region,
			)
		}
	} else {
		_, err = s3connection.CreateBucket(
			&s3.CreateBucketInput{
				Bucket: aws.String(config.Bucket),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create S3 bucket (%w)", err)
		}
	}
	return nil, nil
}
