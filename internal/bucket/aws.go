package bucket

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
)

type AwsConfig struct {
	Config         *aws.Config
	BucketDownload string
	BucketUpload   string
}

func newAwsSession(cfg AwsConfig) *awsSession {
	c := session.New(cfg.Config)
	return &awsSession{
		sess:           c,
		bucketDownload: cfg.BucketDownload,
		bucketUpload:   cfg.BucketUpload,
	}
}

type awsSession struct {
	sess           *session.Session
	bucketDownload string
	bucketUpload   string
}

func (a *awsSession) Download(filename string, path string) (file *os.File, err error) {
	file, err = os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(a.sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(a.bucketDownload),
			Key:    aws.String(filename),
		})

	return
}

func (a *awsSession) Upload(file io.Reader, key string) error {
	uploader := s3manager.NewUploader(a.sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.bucketUpload),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}

func (a *awsSession) Delete(key string) error {
	svc := s3.New(a.sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(a.bucketUpload),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(a.bucketUpload),
		Key:    aws.String(key),
	})
}
