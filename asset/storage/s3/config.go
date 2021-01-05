package s3

import (
	"errors"
)

type S3AssetStorageConfig struct {
	AccessKey            string `json:"accessKey" yaml:"accessKey" env:"AWS_ACCESS_KEY_ID"`
	SecretKey            string `json:"secretKey" yaml:"secretKey" env:"AWS_SECRET_ACCESS_KEY"`
	Bucket               string `json:"bucket" yaml:"bucket" env:"AWS_SECRET_ACCESS_KEY"`
	Region               string `json:"region" yaml:"region" env:"AWS_REGION"`
	Endpoint             string `json:"endpoint" yaml:"endpoint" env:"AWS_S3_ENDPOINT"`
	ForcePathStyleAccess bool   `json:"s3ForcePathStyle" yaml:"s3ForcePathStyle" env:"AWS_S3_PATH_STYLE_ACCESS"`
}

// Validate validates the configuration structure
func (s *S3AssetStorageConfig) Validate() error {
	if s.AccessKey == "" {
		return errors.New("access key cannot be empty")
	}
	if s.SecretKey == "" {
		return errors.New("secret key cannot be empty")
	}
	if s.Bucket == "" {
		return errors.New("bucket name cannot be empty")
	}
	if s.Region == "" {
		return errors.New("region name cannot be empty")
	}
	return nil
}
