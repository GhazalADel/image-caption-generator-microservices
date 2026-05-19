package objectstore

import (
	objectStorageConfig "ImageGenerationService/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/net/context"
)

func GetClient() (*s3.Client, error) {
	cfg, err := objectStorageConfig.NewConfig()
	if err != nil {
		return nil, err
	}

	region := cfg.ObjectStorage.Region

	objStorageCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	objStorageCfg.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     cfg.ObjectStorage.AccessKey,
			SecretAccessKey: cfg.ObjectStorage.SecretKey,
		}, nil
	})

	objStorageCfg.BaseEndpoint = aws.String(cfg.ObjectStorage.EndPoint)

	client := s3.NewFromConfig(objStorageCfg)

	return client, nil
}

func GetBucketName() (string, error) {
	cfg, err := objectStorageConfig.NewConfig()
	if err != nil {
		return "", err
	}

	bucketName := cfg.ObjectStorage.BucketName
	return bucketName, nil
}

func GetEndpoint() (string, error) {
	cfg, err := objectStorageConfig.NewConfig()
	if err != nil {
		return "", err
	}

	endPoint := cfg.ObjectStorage.EndPoint
	return endPoint, nil
}
