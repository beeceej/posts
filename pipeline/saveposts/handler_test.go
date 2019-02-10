package saveposts

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/beeceej/inflight"
	"github.com/beeceej/posts/pipeline/shared/post"
)

func TestUnmarshal(t *testing.T) {
	t.Skip()
	var (
		cfg aws.Config
		err error
	)
	if cfg, err = external.LoadDefaultAWSConfig(); err != nil {
		panic(err.Error())
	}
	s3svc := s3.New(cfg)
	dynamosvc := dynamodb.New(cfg)
	postTableName := "blog-posts"
	h := &Handler{
		Inflight: inflight.NewInflight(
			inflight.Bucket("beeceej-pipelines"),
			inflight.KeyPath("blog-post-pipeline"),
			s3svc),
		PostWriter: &post.PostDynamoRepository{
			DynamoDBAPI: dynamosvc,
			TableName:   postTableName,
		},
	}
	err = h.Handle(inflight.Ref{
		Bucket: "beeceej-pipelines",
		Path:   "blog-post-pipeline",
		Object: "f54c0cd288a5bcda9b6a376da558d10b",
	})

	if err != nil {
		panic(err.Error())
	}

	_ = h
}
