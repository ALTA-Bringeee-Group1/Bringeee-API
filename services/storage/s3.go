package storage

import (
	"bringeee-capstone/configs"
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/gommon/log"
)
type S3 struct {
	awsConfig aws.Config
}

func NewS3() *S3 {
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
		log.Warn("Cannot connect to AWS S3 service")
	}
	return &S3{
		awsConfig: awsConfig,
	}
}

/*
 * Upload from request
 * -------------------------------
 * Upload ke storage service dengan file yang bersumber dari request
 *
 * @param 	fileNamePath 	nama file beserta path yang berada pada cloud storage service
 * @param 	fileHeader	 	file yang akan di upload
 * @return 	string			fileUrl hasil kembalian dari hasil upload
 * @return 	error			error
 */
func (storage S3) UploadFromRequest(fileNamePath string, file multipart.File) (string, error) {
	// s3 Client
	client := s3.NewFromConfig(storage.awsConfig)
	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(configs.Get().AwsS3.Bucket),
		Key:    aws.String(fileNamePath),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}

/*
 * Delete file
 * -------------------------------
 * Hapus file yang berada pada cloud storage service
 *
 * @param 	fileNamePath 	nama file beserta path yang berada pada cloud storage service
 * @return 	string			fileUrl hasil kembalian dari hasil upload
 * @return 	error			error
 */
func (storage S3) Delete(fileNamePath string) error {
	// s3 Client
	client := s3.NewFromConfig(storage.awsConfig)
	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(configs.Get().AwsS3.Bucket),
		Key:    aws.String(fileNamePath),
	})
	if err != nil {
		return err
	}
	return nil
}
