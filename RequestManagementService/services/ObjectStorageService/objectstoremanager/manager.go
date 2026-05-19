package objectstoremanager

import (
	"RequestManagementService/services/ObjectStorageService/objectstore"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
)

func AddPictureToDataStorage(path string, fileName string, fileContent io.Reader) error {
	client, err := objectstore.GetClient()
	if err != nil {
		return err
	}
	destinationKey := path + fileName

	bucketName, err := objectstore.GetBucketName()
	if err != nil {
		return err
	}

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(destinationKey),
		Body:   fileContent,
	})

	return err
}

func GetPictureURL(path, fileName string) (string, error) {
	bucketName, err := objectstore.GetBucketName()
	if err != nil {
		return "", err
	}
	endpoint, err := objectstore.GetEndpoint()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", endpoint, bucketName, path+fileName), nil
}

func GetPictureFromDataStorage(path string, fileName string) ([]byte, error) {
	client, err := objectstore.GetClient()
	if err != nil {
		return nil, err
	}

	destinationKey := path + fileName

	bucketName, err := objectstore.GetBucketName()
	if err != nil {
		return nil, err
	}

	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(destinationKey),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get picture from object storage: %v", err)
	}
	defer result.Body.Close()

	bodyBytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image body: %v", err)
	}

	return bodyBytes, nil
}
