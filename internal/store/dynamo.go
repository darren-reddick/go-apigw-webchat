package store

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/darren-reddick/go-apigw-webchat/internal/session"
)

type ConnectionStoreDynamo struct {
	TableName string
	Session   dynamodbiface.DynamoDBAPI
}

type ConnectionId struct {
	Id string `json:"Id"`
}

func NewConnectionStoreDynamo(tableName string) *ConnectionStoreDynamo {
	sess, err := session.NewDynamoDBSession()
	if err != nil {
		panic(fmt.Errorf("error connecting to dynamodb table: %s", err))
	}
	return &ConnectionStoreDynamo{
		TableName: tableName,
		Session:   sess,
	}
}

func (c *ConnectionStoreDynamo) List() []string {
	input := &dynamodb.ScanInput{
		TableName: &c.TableName,
	}

	ids, err := c.Session.Scan(input)
	if err != nil {
		fmt.Println("Failed to scan database")
	}
	values := make([]string, 0, len(ids.Items))
	for _, v := range ids.Items {
		item := ConnectionId{}
		_ = dynamodbattribute.UnmarshalMap(v, &item)
		values = append(values, item.Id)
	}
	return values
}

func (c *ConnectionStoreDynamo) Add(id string) error {

	connectionItem := ConnectionId{
		Id: id,
	}
	attributeValues, _ := dynamodbattribute.MarshalMap(connectionItem)
	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: &c.TableName,
	}
	_, err := c.Session.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConnectionStoreDynamo) Remove(id string) error {

	connectionItem := ConnectionId{
		Id: id,
	}
	attributeValues, _ := dynamodbattribute.MarshalMap(connectionItem)
	input := &dynamodb.DeleteItemInput{
		Key:       attributeValues,
		TableName: aws.String(c.TableName),
	}
	_, err := c.Session.DeleteItem(input)
	if err != nil {
		return err
	}
	return nil
}
