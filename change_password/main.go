package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type Body struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ChangePassword(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body Body
	b64String, _ := base64.StdEncoding.DecodeString(req.Body)
	rawIn := json.RawMessage(b64String)
	bodyBytes, err := rawIn.MarshalJSON()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving AWS credentials",
		}, nil
	}

	cip := cognitoidentityprovider.NewFromConfig(cfg)
	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
		AuthParameters: map[string]string{
			"USERNAME": body.Username,
			"PASSWORD": body.OldPassword,
		},
	}
	authResp, err := cip.InitiateAuth(context.TODO(), authInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	challengeInput := &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: "NEW_PASSWORD_REQUIRED",
		ClientId:      aws.String(os.Getenv("CLIENT_ID")),
		ChallengeResponses: map[string]string{
			"USERNAME":     body.Username,
			"NEW_PASSWORD": body.NewPassword,
		},
		Session: authResp.Session,
	}
	challengeResp, err := cip.RespondToAuthChallenge(context.TODO(), challengeInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	res, _ := json.Marshal(challengeResp)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(res),
	}, nil
}

func main() {
	lambda.Start(ChangePassword)
}
