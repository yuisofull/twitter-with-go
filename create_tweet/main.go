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
	"path/filepath"
	"strings"
	"time"
)

type Tweet struct {
	UserID    string    `json:"user_id,omitempty" dynamo:"UserID,hash"`
	Text      string    `json:"text_content"`
	ID        int       `json:"-" dynamo:"TweetID,range"`
	Status    int       `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Image     string    `json:"image,omitempty"`
}

const TableName = "tweet"

var (
	bucketName = "go-food-delivery2" //change to your bucket name
	region     = "ap-southeast-1"
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

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: "Error while retrieving AWS credentials",
		}, nil
	}

	image, err := Upload(req, cfg)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: err.Error(),
		}, nil
	}

	var tweet Tweet

	r, err := awslambda.NewReaderMultipart(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: err.Error(),
		}, nil
	}

	form, err := r.ReadForm(1024)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: err.Error(),
		}, nil
	}

	tweet.ID = time.Now().Nanosecond()
	tweet.CreatedAt = time.Now()
	tweet.UpdatedAt = time.Now()
	tweet.Status = 1
	tweet.Text = form.Value["text_content"][0]
	tweet.UserID = form.Value["user_id"][0]
	tweet.Image = image

	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String("ap-south-1")})
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
	lambda.Start(create)
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
