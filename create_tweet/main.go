package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"net/http"
	"time"
)

type Tweet struct {
	UserID    string    `json:"user_id,omitempty" dynamo:"UserID,hash"`
	Text      string    `json:"text_content"`
	ID        int       `json:"id,omitempty" dynamo:"TweetID,range"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	//FakeID    *UID       `json:"id" gorm:"-"`
	//Images          *common.Images     `json:"image" gorm:"column:image;" form:"-"`
}

const TableName = "tweet"

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var tweet Tweet
	err := json.Unmarshal([]byte(req.Body), &tweet)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	tweet.CreatedAt = time.Now()
	tweet.UpdatedAt = time.Now()

	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String("ap-south-1")})
	client := db.Table(TableName)

	err = client.Put(tweet).Run()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error while creating tweet: %s", err.Error()),
		}, nil
	}

	var newTweet Tweet

	err = client.Get("UserID", tweet.UserID).Range("TweetID", dynamo.Equal, tweet.ID).One(&newTweet)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error while getting tweet: %s", err.Error()),
		}, nil
	}

	res, err := json.Marshal(newTweet)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(res),
	}, nil
}

func main() {
	lambda.Start(create)
}
