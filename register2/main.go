package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	types2 "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"net/http"
	"os"
	"strings"
)

var UserTable = os.Getenv("USER_TABLE_NAME")

type Body struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func validate(body Body) error {
	body.Username = strings.TrimSpace(body.Username)
	body.Password = strings.TrimSpace(body.Password)
	body.Email = strings.TrimSpace(body.Email)
	if body.Username == "" {
		return errors.New("username is required")
	}
	if body.Password == "" {
		return errors.New("password is required")
	}
	if body.Email == "" {
		return errors.New("email is required")
	}
	return nil

}

func register(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	if err := validate(body); err != nil {
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
	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
		Username: aws.String(body.Username),
		Password: aws.String(body.Password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(body.Email),
			},
		},
	}

	signUpResp, err := cip.SignUp(context.TODO(), signUpInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	{
		svc := dynamodb.NewFromConfig(cfg)
		_, err := svc.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName: aws.String(UserTable),
			Key: map[string]types2.AttributeValue{
				"UserID": &types2.AttributeValueMemberS{Value: *signUpResp.UserSub},
			},
			ExpressionAttributeValues: map[string]types2.AttributeValue{
				":email":    &types2.AttributeValueMemberS{Value: body.Email},
				":username": &types2.AttributeValueMemberS{Value: body.Username},
			},
			UpdateExpression: aws.String("SET Email = :email, Username = :username"),
			ReturnValues:     types2.ReturnValueAllNew,
		})
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		}
	}

	res, _ := json.Marshal(signUpResp)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(res),
	}, nil
}

func main() {
	lambda.Start(register)
}
