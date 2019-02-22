package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/beeceej/inflight"
	ext "github.com/beeceej/posts/pipeline/shared/external"
	"github.com/beeceej/posts/pipeline/shared/post"
)

type handler struct {
	Inflight  *inflight.Inflight
	mediumAPI *ext.MediumAPI
}

func (h *handler) Handle(ref inflight.Ref) (*inflight.Ref, error) {
	b, err := h.Inflight.Get(ref.Object)
	if err != nil {
		return nil, err
	}
	postIndex := &post.PostIndex{
		Posts: []post.Post{},
	}

	if err = json.Unmarshal(b, &postIndex.Posts); err != nil {
		return nil, err
	}

	authorID, err := h.mediumAPI.GetAuthorID()
	if err != nil {
		return nil, err
	}
	fmt.Println(authorID)

	fmt.Println(len(postIndex.Posts))
	for _, p := range postIndex.Posts {
		err = h.mediumAPI.CreatePost(p, authorID)
	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func main() {
	var (
		mediumToken        string
		inflightBucketName string
		inflightKeyPath    string
		s3svc              s3iface.S3API
		cfg                aws.Config
		err                error
	)
	mediumToken = os.Getenv("MEDIUM_INTEGRATION_TOKEN")
	inflightBucketName = os.Getenv("INFLIGHT_BUCKET_NAME")
	inflightKeyPath = os.Getenv("PIPELINE_SUB_PATH")
	mediumToken = "2d81a795e4ce471ad826cf82b029cdab1136b5d5f2a5651b495f22071426385ef"
	inflightBucketName = "beeceej-pipelines"
	inflightKeyPath = "blog-post-pipeline"

	if mediumToken == "" {
		panic("Missing MEDIUM_INTEGRATION_TOKEN env var")
	}
	if inflightBucketName == "" {
		panic("Missing INFLIGHT_BUCKET_NAME env var")
	}
	if inflightKeyPath == "" {
		panic("Missing PIPELINE_SUB_PATH env var")
	}
	if cfg, err = external.LoadDefaultAWSConfig(); err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	s3svc = s3.New(cfg)

	mediumAPI := &ext.MediumAPI{
		Client:           http.DefaultClient,
		IntegrationToken: mediumToken,
	}

	h := &handler{
		Inflight: inflight.NewInflight(
			inflight.Bucket(inflightBucketName),
			inflight.KeyPath(inflightKeyPath),
			s3svc),
		mediumAPI: mediumAPI,
	}

	fmt.Println(h.Handle(inflight.Ref{
		Bucket: inflightBucketName,
		Path:   inflightKeyPath,
		Object: "272c46f7094da5e2b71e3da571bc96f6",
	}))
	//lambda.Start(h.Handle)
}
