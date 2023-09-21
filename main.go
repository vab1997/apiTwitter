package main

import (
	"context"
	"os"
	"strings"

	"apiTwitter/awsConfig"
	"apiTwitter/db"
	"apiTwitter/handlers"
	"apiTwitter/models"
	"apiTwitter/secretManager"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(ExecuteLambda)
}

func ExecuteLambda(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsConfig.InitAWS()

	var validate = ValidateParams()
	if !validate {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error environment variables. Should be set: SecretName, BucketName, UrlPrefix",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		return res, nil
	}

	SecretModel, err := secretManager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error to get secret " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		return res, nil
	}

	path := strings.Replace(req.PathParameters["api-twitter-go"], os.Getenv("UrlPrefix"), "", -1)

	// create a map with values for context
	options := map[string]string{
		"path":       path,
		"method":     req.HTTPMethod,
		"user":       SecretModel.Username,
		"password":   SecretModel.Password,
		"host":       SecretModel.Host,
		"database":   SecretModel.DataBase,
		"jwtSign":    SecretModel.JWTSign,
		"body":       req.Body,
		"bucketName": os.Getenv("BucketName"),
	}

	for k, v := range options {
		awsConfig.Ctx = context.WithValue(awsConfig.Ctx, models.Key(k), v)
	}

	// awsConfig.Ctx = context.WithValue(awsConfig.Ctx, models.Key("path"), path)
	// awsConfig.Ctx = context.WithValue(awsConfig.Ctx, models.Key("method"), req.HTTPMethod)
	// awsConfig.Ctx = context.WithValue(awsConfig.Ctx, models.Key("user"), SecretModel.Username)
	// awsConfig.Ctx = context.WithValue(awsConfig.Ctx, models.Key("password"), SecretModel.Password)
	// awsConfig.Ctx = context.WithValue(awsConfig.Ctx, models.Key("host"), SecretModel.Host)
	// awsConfig.Ctx = context.WithValue(awsConfig.Ctx, models.Key("database"), SecretModel.DataBase)
	// awsConfig.Ctx = context.WithValue(awsConfig.Ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	// awsConfig.Ctx = context.WithValue(awsConfig.Ctx, models.Key("body"), req.Body)
	// awsConfig.Ctx = context.WithValue(awsConfig.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	// check connection to MongoDB
	err = db.ConnectDB(awsConfig.Ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error to connect to MongoDB " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		return res, nil
	}

	resApi := handlers.Handlers(awsConfig.Ctx, req)
	if resApi.CustomResp == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: resApi.Status,
			Body:       resApi.Message,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	} else {
		return resApi.CustomResp, nil
	}
}

func ValidateParams() bool {
	_, params := os.LookupEnv("SecretName")
	if !params {
		return params
	}

	_, params = os.LookupEnv("BucketName")
	if !params {
		return params
	}

	_, params = os.LookupEnv("UrlPrefix")
	if !params {
		return params
	}

	return params
}
