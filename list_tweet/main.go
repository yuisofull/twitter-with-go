package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"net/http"
	"time"
)

type Filter struct {
	//FakeUserID string `json:"-" form:"user_id"`
	UserID int    `json:"user_id,omitempty" form:"-"`
	Status []int  `json:"-"`
	Search string `json:"search,omitempty" form:"search"`
}

type Tweet struct {
	UserID    string     `json:"user_id,omitempty"`
	Text      string     `json:"text_content"`
	Id        string     `json:"id,omitempty"`
	Status    string     `json:"status,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	//FakeID    *UID       `json:"id" gorm:"-"`
	//Images          *common.Images     `json:"image" gorm:"column:image;" form:"-"`
}

const TableName = "tweet"

func list(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving AWS credentials",
		}, nil
	}

	svc := dynamodb.NewFromConfig(cfg)
	out, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(TableName),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	var tweets []Tweet
	err = attributevalue.UnmarshalListOfMaps(out.Items, &tweets)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while Unmarshal tweets",
		}, nil
	}

	res, _ := json.Marshal(tweets)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(res),
	}, nil
}

func main() {
	lambda.Start(list)
}
