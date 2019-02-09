package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/beeceej/inflight"
	"github.com/beeceej/posts/pipeline/upload"
)

var (
	s3svc   s3iface.S3API
	cfg     aws.Config
	handler *upload.Handler
	err     error
)

func init() {
	if cfg, err = external.LoadDefaultAWSConfig(); err != nil {
		panic(err.Error())
	}
	s3svc = s3.New(cfg)

	inflightBucketName := os.Getenv("INFLIGHT_BUCKET_NAME")
	pipelineSubPath := os.Getenv("PIPELINE_SUB_PATH")
	if inflightBucketName == "" {
		panic("Missing env var INFLIGHT_BUCKET_NAME")
	}
	if pipelineSubPath == "" {
		panic("Missing env var PIPELINE_SUB_PATH")
	}
	handler = &upload.Handler{
		Inflight: inflight.NewInflight(
			inflight.Bucket(inflightBucketName),
			inflight.KeyPath(pipelineSubPath),
			s3svc),
		Uploader: &upload.Uploader{
			S3API: s3svc,
		},
	}
}

func main() {
	lambda.Start(handler.Handle)
}
