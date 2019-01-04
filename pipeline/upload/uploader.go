package upload

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/beeceej/posts/pipeline/shared/domain"
)

type Uploader struct {
	s3iface.S3API
}

func (u Uploader) Upload(postIndex *domain.PostIndex) error {
	var (
		b                     []byte
		existingPublishedPost *domain.Post
	)

	for _, post := range postIndex.Posts {
		getObjReq := u.S3API.GetObjectRequest(
			&s3.GetObjectInput{
				Bucket: aws.String("static.beeceej.com"),
				Key:    aws.String(filepath.Join("posts", post.NormalizedTitle) + ".json"),
			},
		)
		a, err := getObjReq.Send()
		if err != nil {
			fmt.Println(post.NormalizedTitle, err.Error())
		} else {
			defer a.Body.Close()
			b, err = ioutil.ReadAll(a.Body)
			if err != nil {
				return err
			}

			existingPublishedPost := new(domain.Post)
			existingPublishedPost.FromBytes(b)
		}
		if existingPublishedPost != nil {
			fmt.Println("existingPublishedPost.MD5", existingPublishedPost.MD5, " post.MD5", post.MD5)
			if existingPublishedPost.MD5 != post.MD5 {
				putObjReq := u.PutObjectRequest(&s3.PutObjectInput{
					Bucket: aws.String("static.beeceej.com"),
					Key:    aws.String(filepath.Join("posts", post.NormalizedTitle) + ".json"),
					Body:   bytes.NewReader(post.ToBytes()),
				})
				putObjReq.Send()
			}
		} else {
			putObjReq := u.PutObjectRequest(&s3.PutObjectInput{
				Bucket: aws.String("static.beeceej.com"),
				Key:    aws.String(filepath.Join("posts", post.NormalizedTitle) + ".json"),
				Body:   bytes.NewReader(post.ToBytes()),
			})
			putObjReq.Send()
		}

	}

	return nil
}
