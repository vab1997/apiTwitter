package handlers

import (
	"context"
	"fmt"

	"apiTwitter/models"

	"github.com/aws/aws-lambda-go/events"
)

func Handlers(ctx context.Context, req events.APIGatewayProxyRequest) models.ResAPI {
	fmt.Println("> proccesing " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))

	var res models.ResAPI
	res.Status = 400

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		}
	case "GET":
		switch ctx.Value(models.Key("path")).(string) {
		}
	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {
		}
	case "DELETE":
		switch ctx.Value(models.Key("path")).(string) {
		}
	}

	res.Message = "Method invalid"

	return res
}
