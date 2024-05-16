package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	aws1 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/grokify/go-awslambda"
	"github.com/guregu/dynamo"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Tweet struct {
	UserID    string    `json:"-" dynamo:"UserID,hash"`
	Text      string    `json:"text_content"`
	ID        int       `json:"-" dynamo:"TweetID,range"`
	Status    int       `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Image     string    `json:"image,omitempty"`
}

var TableName = os.Getenv("TWEET_TABLE_NAME")

var (
	bucketName = os.Getenv("BUCKET_NAME") //change to your bucket name
	region     = os.Getenv("REGION")
)

type customStruct struct {
	Content       string `json:"content,omitempty"`
	FileName      string `json:"fileName,omitempty"`
	FileExtension string `json:"fileExtension,omitempty"`
	Link          string `json:"link,omitempty"`
}

func GetFileFromAPIGatewayProxyRequest(req events.APIGatewayProxyRequest) ([]byte, string, error) {
	r, err := awslambda.NewReaderMultipart(req)
	if err != nil {
		return []byte{}, "", err
	}

	part, err := r.NextPart()
	if err != nil {
		return []byte{}, "", err

	}

	content, err := io.ReadAll(part)
	if err != nil {
		return []byte{}, "", err
	}

	return content, part.FileName(), nil
}

func fileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func Upload(request events.APIGatewayProxyRequest, cfg aws1.Config) (image string, err error) {
	content, fileName, err := GetFileFromAPIGatewayProxyRequest(request)
	if err != nil {
		return
	}

	cfg.Region = region
	client := s3.NewFromConfig(cfg)

	fileExt := filepath.Ext(fileName)
	fileName = fmt.Sprintf("%s-%v%s", fileNameWithoutExtSliceNotation(fileName), time.Now().UnixNano(), fileExt)
	fileName = strings.Replace(fileName, " ", "-", -1)

	data := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
		Body:   bytes.NewReader(content),
	}

	_, err = client.PutObject(context.TODO(), data)
	if err != nil {
		return
	}

	image = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, fileName)

	return
}

func InvalidRequest(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: fmt.Sprintf("Error: %s", err.Error()),
	}, nil
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
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return InvalidRequest(err)
	}

	image, err := Upload(req, cfg)
	if err != nil {
		return InvalidRequest(err)
	}

	var tweet Tweet

	r, err := awslambda.NewReaderMultipart(req)
	if err != nil {
		return InvalidRequest(err)
	}

	form, err := r.ReadForm(1024)
	if err != nil {
		return InvalidRequest(err)
	}

	err = json.Unmarshal([]byte(form.Value["data"][0]), &tweet)
	if err != nil {
		return InvalidRequest(err)
	}
	tweet.ID = time.Now().Nanosecond()
	tweet.CreatedAt = time.Now()
	tweet.UpdatedAt = time.Now()
	tweet.Status = 1
	//tweet.Text = form.Value["text_content"][0]
	//tweet.UserID = form.Value["user_id"][0]
	tweet.Image = image
	tweet.UserID, err = getUserID(req)
	if err != nil {
		return InvalidRequest(err)
	}
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	client := db.Table(TableName)

	err = client.Put(tweet).Run()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error while creating tweet: %s", err.Error()),
		}, err
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
	lambda.Start(lambdaHandler)
}

//func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	var tweet Tweet
//	err := json.Unmarshal([]byte(req.Body), &tweet)
//	tweet.ID = time.Now().Nanosecond()
//	if err != nil {
//		return events.APIGatewayProxyResponse{
//			StatusCode: 400,
//			Body:       err.Error(),
//		}, nil
//	}
//
//	tweet.CreatedAt = time.Now()
//	tweet.UpdatedAt = time.Now()
//
//	sess := session.Must(session.NewSession())
//	db := dynamo.New(sess, &aws.Config{Region: aws.String("ap-south-1")})
//	client := db.Table(TableName)
//
//	err = client.Put(tweet).Run()
//	if err != nil {
//		return events.APIGatewayProxyResponse{
//			StatusCode: http.StatusInternalServerError,
//			Body:       fmt.Sprintf("Error while creating tweet: %s", err.Error()),
//		}, nil
//	}
//
//	var newTweet Tweet
//
//	err = client.Get("UserID", tweet.UserID).Range("TweetID", dynamo.Equal, tweet.ID).One(&newTweet)
//	if err != nil {
//		return events.APIGatewayProxyResponse{
//			StatusCode: http.StatusInternalServerError,
//			Body:       fmt.Sprintf("Error while getting tweet: %s", err.Error()),
//		}, nil
//	}
//
//	res, err := json.Marshal(newTweet)
//	if err != nil {
//		return events.APIGatewayProxyResponse{
//			StatusCode: 400,
//			Body:       err.Error(),
//		}, nil
//	}
//	return events.APIGatewayProxyResponse{
//		StatusCode: 200,
//		Headers: map[string]string{
//			"Content-Type": "application/json",
//		},
//		Body: string(res),
//	}, nil
//}
