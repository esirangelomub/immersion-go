package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/esirangelomub/immersion-go/internal/bucket"
	"github.com/esirangelomub/immersion-go/internal/queue"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	// rabbitmq config
	qcfg := queue.RabbitMQConfig{
		URL:       os.Getenv("RABBITMQ_URL"),
		TopicName: os.Getenv("RABBITMQ_TOPIC_NAME"),
		Timeout:   time.Second * 30,
	}

	// create a new queue connection
	qc, err := queue.New(queue.RabbitMQ, qcfg)
	if err != nil {
		panic(err)
	}

	// consume messages
	c := make(chan queue.QueueDto)
	qc.Consume(c)

	// bucket config
	bcfg := bucket.AwsConfig{
		Config: &aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		},
		BucketDownload: "aprenda-golang-drive-raw",
		BucketUpload:   "aprenda-golang-drive-gzip",
	}

	b, err := bucket.New(bucket.AwsProvider, bcfg)
	if err != nil {
		panic(err)
	}

	for msg := range c {
		// download file
		src := fmt.Sprintf("%s/%s", msg.Path, msg.Filename)
		dst := fmt.Sprintf("%d_%s", msg.ID, msg.Filename)
		file, err := b.Download(src, dst)
		if err != nil {
			log.Printf("error downloading file: %s", err)
			continue
		}

		body, err := io.ReadAll(file)
		if err != nil {
			log.Printf("error reading file: %s", err)
			continue
		}

		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err = zw.Write(body)
		if err != nil {
			log.Printf("error compressing file: %s", err)
			continue
		}

		err = zw.Close()
		if err != nil {
			log.Printf("error closing gzip writer: %s", err)
			continue
		}

		// upload file
		zr, err := gzip.NewReader(&buf)
		if err != nil {
			log.Printf("error creating gzip reader: %s", err)
			continue
		}

		err = b.Upload(zr, src)
		if err != nil {
			log.Printf("error uploading file: %s", err)
			continue
		}

		err = os.Remove(dst)
		if err != nil {
			log.Printf("error removing local file: %s", err)
			continue
		}
	}
}
