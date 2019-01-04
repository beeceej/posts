package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/beeceej/posts/pipeline/poststojson"
	"github.com/beeceej/posts/pipeline/shared/inflight"
)

var (
	s3svc   s3iface.S3API
	cfg     aws.Config
	handler *poststojson.Handler
	err     error
)

func init() {
	if cfg, err = external.LoadDefaultAWSConfig(); err != nil {
		panic(err.Error())
	}
	s3svc = s3.New(cfg)

	inflightBucketName := os.Getenv("INFLIGHT_BUCKET_NAME")
	pipelineSubPath := os.Getenv("PIPELINE_SUB_PATH")
	postsRepositoryURL := os.Getenv("POSTS_REPO_URI")
	if inflightBucketName == "" {
		panic("Missing env var INFLIGHT_BUCKET_NAME")
	}
	if pipelineSubPath == "" {
		panic("Missing env var PIPELINE_SUB_PATH")
	}
	handler = &poststojson.Handler{
		S3API: s3svc,
		Inflight: &inflight.Inflight{
			S3API:   s3svc,
			Bucket:  inflightBucketName,
			KeyPath: pipelineSubPath,
		},
		PostsRepositoryURL: postsRepositoryURL,
	}
}

func main() {
	lambda.Start(handler.Handle)
}
