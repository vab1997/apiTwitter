package db

import (
	"context"
	"fmt"

	"apiTwitter/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoCN *mongo.Client
var dababaseName string

func ConnectDB(ctx context.Context) error {
	user := ctx.Value(models.Key("user")).(string)
	password := ctx.Value(models.Key("password")).(string)
	host := ctx.Value(models.Key("host")).(string)
	connString := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", user, password, host)

	var clientOptions = options.Client().ApplyURI(connString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error to connect to MongoDB " + err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error to ping MongoDB " + err.Error())
		return err
	}

	fmt.Println("Connection to MongoDB OK")
	MongoCN = client
	dababaseName = ctx.Value(models.Key("database")).(string)

	return nil
}

func CheckConnection() bool {
	err := MongoCN.Ping(context.TODO(), nil)
	return err != nil
}
