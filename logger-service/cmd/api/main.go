package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "8083"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

type Config struct {
}

func main() {
	_, err := DbConnect()
	if err != nil {
		log.Panic("Can't connect to DB")
	}

}

func DbConnect() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	connect, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("Can't connect to MongoDB")
		return nil, err
	}

	err = connect.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return connect, nil
}
