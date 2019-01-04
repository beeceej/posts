package saveposts

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/beeceej/posts/pipeline/shared/domain"
)

type PostSaver struct {
	dynamodbiface.DynamoDBAPI
	TableName string
}

func (p *PostSaver) SavePosts(posts []domain.Post) error {
	for _, post := range posts {
		avm, err := dynamodbattribute.MarshalMap(post)
		if err != nil {
			return err
		}
		req := p.PutItemRequest(&dynamodb.PutItemInput{
			TableName: aws.String(p.TableName),
			Item:      avm,
		})
		req.Send()
	}
	return nil
}
