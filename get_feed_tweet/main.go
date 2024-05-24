package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"os"
	"time"
)

var FeedTweetTableName = os.Getenv("FEED_TWEET_TABLE_NAME")
var region = os.Getenv("REGION")

type Tweet struct {
	UserID    string    `json:"UserID" dynamo:"UserID,hash"`
	Text      string    `json:"text_content"`
	ID        int       `json:"TweetID" dynamo:"TweetID,range"`
	Status    int       `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Image     string    `json:"image,omitempty"`
}

func getUserID(req events.APIGatewayProxyRequest) (string, error) {
	if req.RequestContext.Authorizer == nil {
		return "", fmt.Errorf("no authorizer found")
	}

	if req.RequestContext.Authorizer["claims"] == nil {
		return "", fmt.Errorf("no claims found")
	}

	claims := req.RequestContext.Authorizer["claims"].(map[string]interface{})

	if claims["sub"] == nil {
		return "", fmt.Errorf("no username found")
	}

	return claims["sub"].(string), nil
}

func lambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID, err := getUserID(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Error getting user ID: %s", err.Error()),
		}, nil
	}

	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	client := db.Table(FeedTweetTableName)

	var tweets []Tweet
	err = client.Get("UserID", userID).All(&tweets)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error retrieving tweets: %s", err.Error()),
		}, nil
	}

	res, err := json.Marshal(tweets)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error marshalling tweets: %s", err.Error()),
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
	lambda.Start(lambdaHandler)
}
