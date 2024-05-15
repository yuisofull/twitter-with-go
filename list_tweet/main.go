package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/btcsuite/btcutil/base58"
	"github.com/guregu/dynamo"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type response struct {
	Data   []Tweet `json:"data"`
	Paging Paging  `json:"paging"`
	Filter Filter  `json:"filter"`
}

type Paging struct {
	Limit      int     `json:"limit"`
	Total      int     `json:"total,omitempty"`
	Cursor     *string `json:"cursor,omitempty"`
	NextCursor *string `json:"next_cursor,omitempty"`
}

type Filter struct {
	//FakeUserID string `json:"-" form:"user_id"`
	UserID int    `json:"user_id,omitempty"`
	Status []int  `json:"-"`
	Search string `json:"search,omitempty"`
}

type Tweet struct {
	UserID    string    `json:"user_id,omitempty" dynamo:"UserID,hash"`
	Text      string    `json:"text_content"`
	ID        int       `json:"id,omitempty" dynamo:"TweetID,range"`
	Status    int       `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Image     string    `json:"image"`
}

func Encode(key dynamo.PagingKey) *string {
	hashKey := base58.Encode([]byte(*key["UserID"].S))
	rangeKey := base58.Encode([]byte(*key["TweetID"].N))
	res := fmt.Sprintf("%s.%s", hashKey, rangeKey)
	return &res
}

func Decode(encoded string) (dynamo.PagingKey, error) {
	parts := strings.Split(encoded, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid paging key")
	}
	return dynamo.PagingKey{
		"UserID":  &dynamodb.AttributeValue{S: aws.String(string(base58.Decode(parts[0])))},
		"TweetID": &dynamodb.AttributeValue{N: aws.String(string(base58.Decode(parts[1])))},
	}, nil
}

var TableName = os.Getenv("TWEET_TABLE_NAME")

func list(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var paging Paging
	//var Filter *Filter

	var results []Tweet
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(os.Getenv("REGION"))})
	table := db.Table(TableName)

	scan := table.Scan()

	paging.Limit = 1
	if req.QueryStringParameters["limit"] != "" {
		limit, err := strconv.Atoi(req.QueryStringParameters["limit"])
		if err == nil {
			paging.Limit = limit
		}
	}
	if req.QueryStringParameters["cursor"] != "" {
		jsonStr := req.QueryStringParameters["cursor"]
		paging.Cursor = &jsonStr
		key, err := Decode(*paging.Cursor)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       err.Error(),
			}, nil
		}
		scan.StartFrom(key)
	}

	key, err := scan.Limit(int64(paging.Limit)).AllWithLastEvaluatedKey(&results)
	log.Println("results: ", results)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error while getting tweet: %s", err.Error()),
		}, nil
	}
	paging.NextCursor = Encode(key)
	paging.Total = len(results)

	res, err := json.Marshal(response{Data: results, Paging: paging})
	log.Println(string(res))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*", // Required for CORS support to work
		},
		Body: string(res),
	}, nil
}

func main() {
	lambda.Start(list)
}
