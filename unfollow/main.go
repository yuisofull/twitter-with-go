package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"os"
)

var FollowTableName = os.Getenv("FOLLOW_TABLE_NAME")
var region = os.Getenv("REGION")

type Follow struct {
	FollowerID string `json:"FollowerID" dynamo:"FollowerID,hash"`
	FolloweeID string `json:"FolloweeID" dynamo:"FolloweeID,range"`
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
	followerID, err := getUserID(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Error getting user ID: %s", err.Error()),
		}, nil
	}

	followeeID := req.PathParameters["followeeID"]
	if followeeID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "FolloweeID is required",
		}, nil
	}

	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	client := db.Table(FollowTableName)

	err = client.Delete("FollowerID", followerID).Range("FolloweeID", followeeID).Run()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error unfollowing user: %s", err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "Unfollowed Successfully",
	}, nil
}

func main() {
	lambda.Start(lambdaHandler)
}
