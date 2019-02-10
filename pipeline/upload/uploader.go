package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/beeceej/posts/pipeline/shared/post"
)

type Uploader struct {
	s3iface.S3API
}

func (u Uploader) Upload(postIndex *post.PostIndex) error {
	var (
		b                     []byte
		existingPublishedPost post.Post
	)

	for _, p := range postIndex.Posts {
		getObjReq := u.S3API.GetObjectRequest(
			&s3.GetObjectInput{
				Bucket: aws.String("static.beeceej.com"),
				Key:    aws.String(filepath.Join("posts", p.NormalizedTitle) + ".json"),
			},
		)
		a, err := getObjReq.Send()
		if err != nil {
			fmt.Println(p.NormalizedTitle, err.Error())
		} else {
			defer a.Body.Close()
			b, err = ioutil.ReadAll(a.Body)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(b, &existingPublishedPost); err != nil {
				panic(err.Error())
			}
		}
		postPublishedBefore := &existingPublishedPost != nil
		hasMD5Changed := postPublishedBefore && existingPublishedPost.MD5 != p.MD5
		if postPublishedBefore && hasMD5Changed {
			fmt.Println("existingPublishedPost.MD5", existingPublishedPost.MD5, " post.MD5", p.MD5)
			putObjReq := u.PutObjectRequest(&s3.PutObjectInput{
				Bucket:      aws.String("static.beeceej.com"),
				Key:         aws.String(filepath.Join("posts", p.NormalizedTitle) + ".json"),
				Body:        bytes.NewReader(p.ToBytes()),
				ContentType: aws.String("application/json"),
			})
			putObjReq.Send()
		} else {
			putObjReq := u.PutObjectRequest(&s3.PutObjectInput{
				Bucket:      aws.String("static.beeceej.com"),
				Key:         aws.String(filepath.Join("posts", p.NormalizedTitle) + ".json"),
				Body:        bytes.NewReader(p.ToBytes()),
				ContentType: aws.String("application/json"),
			})
			putObjReq.Send()
		}
	}

	return nil
}

// UploadSiteMap is
func (u Uploader) UploadSiteMap(postIndex *post.PostIndex) error {
	var (
		b []byte
	)
	var index *post.PostIndex
	for _, p := range postIndex.Posts {
		index.Posts = append(index.Posts, post.Post{
			ID:              p.ID,
			Title:           p.Title,
			NormalizedTitle: p.NormalizedTitle,
			Author:          p.Author,
			PostedAt:        p.PostedAt,
			UpdatedAt:       p.UpdatedAt,
			Visible:         p.Visible,
			MD5:             p.MD5,
			Blurb:           p.Body[0:100],
		})
	}

	b, _ = json.Marshal(index)
	putObjReq := u.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String("static.beeceej.com"),
		Key:         aws.String("posts/all.json"),
		Body:        bytes.NewReader(b),
		ContentType: aws.String("application/json"),
	})
	putObjReq.Send()
	return nil
}
