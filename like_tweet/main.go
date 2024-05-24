package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"os"
	"strconv"
	"time"
)

var FavoriteTableName = os.Getenv("FAVORITE_TABLE_NAME")
var TweetTableName = os.Getenv("TWEET_TABLE_NAME")
var region = os.Getenv("REGION")

type Favorite struct {
	UserID  string `json:"UserID" dynamo:"UserID,hash"`
	TweetID int    `json:"TweetID" dynamo:"TweetID,range"`
}

type Tweet struct {
	UserID    string    `json:"UserID" dynamo:"UserID,hash"`
	Text      string    `json:"text_content"`
	ID        int       `json:"TweetID" dynamo:"TweetID,range"`
	Status    int       `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Image     string    `json:"image,omitempty"`
	Likes     int       `json:"likes"`
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

	tweetID := req.PathParameters["tweetID"]
	if tweetID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "TweetID is required",
		}, nil
	}

	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	client := db.Table(FavoriteTableName)

	tweetIDInt, err := strconv.Atoi(tweetID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Error converting tweetID to int: %s", err.Error()),
		}, nil
	}
	favorite := Favorite{
		UserID:  userID,
		TweetID: tweetIDInt,
	}

	err = client.Put(favorite).Run()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error liking tweet: %s", err.Error()),
		}, nil
	}

	tweetClient := db.Table(TweetTableName)
	var tweet Tweet
	err = tweetClient.Get("UserID", tweet.UserID).Range("TweetID", dynamo.Equal, tweetIDInt).One(&tweet)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error retrieving tweet: %s", err.Error()),
		}, nil
	}

	tweet.Likes++
	err = tweetClient.Put(tweet).Run()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error updating tweet likes: %s", err.Error()),
		}, nil
	}

	// Here you would send a notification to the user who created the tweet

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "Liked successfully",
	}, nil
}

func main() {
	lambda.Start(lambdaHandler)
}
