package helpers

import (
	"bringeee-capstone/configs"
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadFileToS3(filename string, file multipart.File) (string, error) {
	// Connect AWS
	awsConfig, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			configs.Get().AwsS3.AccessKey,
			configs.Get().AwsS3.SecretKey,
			"",
		)),
		awsConfig.WithRegion(configs.Get().AwsS3.Region),
	)
	if err != nil {
		return "", err
	}

	// s3 Client
	client := s3.NewFromConfig(awsConfig)
	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(configs.Get().AwsS3.Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}

func DeleteFromS3(filename string) error {

	// Connect AWS
	awsConfig, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			configs.Get().AwsS3.AccessKey,
			configs.Get().AwsS3.SecretKey,
			"",
		)),
		awsConfig.WithRegion(configs.Get().AwsS3.Region),
	)
	if err != nil {
		return err
	}

	// s3 Client
	client := s3.NewFromConfig(awsConfig)
	_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(configs.Get().AwsS3.Bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return err
	}
	return nil
}
