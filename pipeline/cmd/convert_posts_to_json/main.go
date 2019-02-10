package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/beeceej/inflight"
	"github.com/beeceej/posts/pipeline/poststojson"
	"github.com/beeceej/posts/pipeline/shared/post"
)

func main() {
	var (
		s3svc   s3iface.S3API
		cfg     aws.Config
		handler *poststojson.Handler
		err     error
	)
	if cfg, err = external.LoadDefaultAWSConfig(); err != nil {
		panic(err.Error())
	}
	s3svc = s3.New(cfg)
	dynamosvc := dynamodb.New(cfg)
	inflightBucketName := os.Getenv("INFLIGHT_BUCKET_NAME")
	pipelineSubPath := os.Getenv("PIPELINE_SUB_PATH")
	postsGitRepositoryURL := os.Getenv("POSTS_REPO_URI")
	postsTableName := os.Getenv("POSTS_TABLE_NAME")
	if inflightBucketName == "" {
		panic("Missing env var INFLIGHT_BUCKET_NAME")
	}
	if pipelineSubPath == "" {
		panic("Missing env var PIPELINE_SUB_PATH")
	}
	if postsTableName == "" {
		panic("Missing env var POSTS_TABLE_NAME")
	}
	postConverter := &poststojson.PostConverter{
		PostGetter: &post.PostDynamoRepository{
			DynamoDBAPI: dynamosvc,
			TableName:   postsTableName,
		}}

	handler = &poststojson.Handler{
		S3API: s3svc,
		Inflight: inflight.NewInflight(
			inflight.Bucket(inflightBucketName),
			inflight.KeyPath(pipelineSubPath),
			s3svc),
		PostConverter:         postConverter,
		PostsGitRepositoryURL: postsGitRepositoryURL,
	}
	lambda.Start(handler.Handle)
}
