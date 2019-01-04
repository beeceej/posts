package inflight

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/beeceej/posts/pipeline/shared/state"
	"github.com/cenkalti/backoff"
	uuid "github.com/satori/go.uuid"
)

type Inflight struct {
	s3iface.S3API
	Bucket  string
	KeyPath string
}

func NewInflight(bucket, keypath string, s3 s3iface.S3API) *Inflight {
	return &Inflight{
		Bucket:  bucket,
		KeyPath: keypath,
		S3API:   s3,
	}
}

// Write will take the data given and attempt to put it in S3
// It then will return the S3 URI back to the caller so that the data may be passed between
// multiple step functions
func (i *Inflight) Write(data io.ReadSeeker) (*state.InflightRef, error) {
	ref := &state.InflightRef{
		Bucket: i.Bucket,
		Object: uuid.NewV4().String(),
		Path:   filepath.Join(i.KeyPath),
	}

	err := backoff.Retry(
		i.tryWriteToS3(data, ref.Object),
		backoff.NewExponentialBackOff(),
	)

	if err != nil {
		return nil, err
	}

	return ref, nil
}

// Get is
func (i *Inflight) Get(object string) ([]byte, error) {
	b := &[]byte{}
	fmt.Println("Getting data from s3:", object)
	err := backoff.Retry(
		i.tryReadFromS3(object, b),
		backoff.NewExponentialBackOff(),
	)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return *b, err
}

func (i *Inflight) tryWriteToS3(data io.ReadSeeker, uri string) func() error {
	bucket := i.Bucket
	keyPath := i.KeyPath
	return func() error {
		req := i.PutObjectRequest(&s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(filepath.Join(keyPath, uri)),
			Body:        data,
			ContentType: aws.String("application/json"),
		})

		res, err := req.Send()

		if err != nil {
			return backoff.Permanent(err)
		}

		if res.SDKResponseMetadata().Request.IsErrorRetryable() {
			return res.SDKResponseMetadata().Request.Error
		}

		return nil
	}
}

func (i *Inflight) tryReadFromS3(object string, data *[]byte) func() error {
	bucket := i.Bucket
	keyPath := i.KeyPath

	return func() error {
		req := i.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filepath.Join(keyPath, object)),
		})

		res, err := req.Send()

		if err != nil {
			return backoff.Permanent(err)
		}

		if res.SDKResponseMetadata().Request.IsErrorRetryable() {
			return res.SDKResponseMetadata().Request.Error
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return backoff.Permanent(err)
		}
		defer res.Body.Close()
		fmt.Println("bytes: ", string(b))

		*data = append(*data, b...)
		return nil
	}
}
