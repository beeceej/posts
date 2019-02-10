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

type getSinglePost struct {
	dynamodbiface.DynamoDBAPI
	post      *Post
	tableName string
}

func (p *PostDynamoRepository) Get(id, md5 string) (post *Post, err error) {
	getter := &getSinglePost{
		DynamoDBAPI: p.DynamoDBAPI,
		tableName:   p.TableName,
	}
	err = backoff.Retry(
		getter.tryGetPost(id, md5),
		backoff.NewExponentialBackOff(),
	)
	if err != nil {
		return nil, err
	}

	return getter.post, nil
}

func (p *getSinglePost) tryGetPost(id, md5 string) func() error {
	return func() (err error) {
		getItemReq := p.GetItemRequest(&dynamodb.GetItemInput{
			Key: map[string]dynamodb.AttributeValue{
				"id": {
					S: aws.String(id),
				},
				"md5": {
					S: aws.String(md5),
				},
			},
			TableName: aws.String(p.tableName),
		})

		var (
			out     *dynamodb.GetItemOutput
			thePost Post
		)
		out, err = getItemReq.Send()
		if err != nil && aws.IsErrorRetryable(err) {
			return err
		} else if err != nil && !aws.IsErrorRetryable(err) {
			return backoff.Permanent(err)
		}

		dynamodbattribute.UnmarshalMap(out.Item, &thePost)
		if &thePost == nil || thePost.ID == "" {
			p.post = nil
		} else {
			p.post = &thePost
		}

		return nil
	}
}

// func (p *PostDynamoRepository) BatchGet(postsLike []Post) (posts []*Post, err error) {
// 	if len(posts) == 0 {
// 		return posts, errors.New("Come on, gotta give me something")
// 	}
// 	postGetter := &postGetter{
// 		tableName:   p.TableName,
// 		DynamoDBAPI: p.DynamoDBAPI,
// 	}

// 	err = backoff.Retry(
// 		postGetter.getBulk(postsLike),
// 		backoff.NewExponentialBackOff(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return postGetter.posts, nil
// }

// func (p *postGetter) getBulk(postsLike []Post) func() error {
// 	return func() error {
// 		getRequests, err := mkBatchGetItemReqs(postsLike)

// 		if err != nil {
// 			return backoff.Permanent(err)
// 		}

// 		batchGetReq := p.BatchGetItemRequest(&dynamodb.BatchGetItemInput{
// 			RequestItems: map[string]dynamodb.KeysAndAttributes{
// 				p.tableName: {
// 					Keys: getRequests, // []map[string]*dynamodb.AttributeValue
// 				},
// 			},
// 		})
// 		var response *dynamodb.BatchGetItemOutput
// 		response, err = batchGetReq.Send()
// 		if err != nil && aws.IsErrorRetryable(err) {
// 			return err
// 		} else if err != nil && !aws.IsErrorRetryable(err) {
// 			return backoff.Permanent(err)
// 		}
// 		fmt.Println(response)
// 		return nil
// 	}
// }

// func mkBatchGetItemReqs(items []Post) (getRequests []map[string]dynamodb.AttributeValue, err error) {
// 	for _, post := range items {
// 		id, hash := post.ID, post.MD5
// 		getRequests = append(getRequests, map[string]dynamodb.AttributeValue{
// 			"id": dynamodb.AttributeValue{
// 				S: aws.String(id),
// 			},
// 			"md5": dynamodb.AttributeValue{
// 				S: aws.String(hash),
// 			},
// 		})
// 	}
// 	return getRequests, nil
// }
