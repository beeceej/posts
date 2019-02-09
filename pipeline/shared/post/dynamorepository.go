package post

import (
	"errors"

	"github.com/cenkalti/backoff"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
)

type PostDynamoRepository struct {
	dynamodbiface.DynamoDBAPI
	TableName string
}

func (p *PostDynamoRepository) Write(posts []Post) error {
	if len(posts) == 0 {
		return errors.New("Come on, gotta give me something")
	}

	if len(posts) == 1 {
		return backoff.Retry(
			p.writeSingle(posts[0]),
			backoff.NewExponentialBackOff(),
		)
	}

	return backoff.Retry(
		p.writeBulk(posts),
		backoff.NewExponentialBackOff(),
	)
}

func (p *PostDynamoRepository) writeBulk(posts []Post) func() error {
	return func() error {
		writeRequests, err := mkBatchWriteReqs(posts)

		if err != nil {
			return backoff.Permanent(err)
		}

		batchWriteReq := p.BatchWriteItemRequest(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]dynamodb.WriteRequest{
				p.TableName: writeRequests,
			},
		})

		_, err = batchWriteReq.Send()

		if err != nil && aws.IsErrorRetryable(err) {
			return err
		} else if err != nil && !aws.IsErrorRetryable(err) {
			return backoff.Permanent(err)
		}

		return nil
	}
}

func (p *PostDynamoRepository) writeSingle(post Post) func() error {
	return func() error {
		avm, err := dynamodbattribute.MarshalMap(post)
		if err != nil {
			return backoff.Permanent(err)
		}

		putItemReq := p.PutItemRequest(&dynamodb.PutItemInput{
			TableName: aws.String(p.TableName),
			Item:      avm,
		})

		_, err = putItemReq.Send()
		if err != nil && aws.IsErrorRetryable(err) {
			return err
		} else if err != nil && !aws.IsErrorRetryable(err) {
			return backoff.Permanent(err)
		}
		return nil
	}
}

func mkBatchWriteReqs(items []Post) (writeRequests []dynamodb.WriteRequest, err error) {
	for _, post := range items {
		avm, err := dynamodbattribute.MarshalMap(post)
		if err != nil {
			return nil, err
		}
		req := dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: avm,
			}}
		writeRequests = append(writeRequests, req)
	}
	return writeRequests, nil
}
