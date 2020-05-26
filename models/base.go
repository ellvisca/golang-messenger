package models

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func init() {
	// Load .env
	godotenv.Load()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://mongo-0.mongo:27017/dbname_?")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	db = client.Database("messenger")
}

func GetDB() *mongo.Database {
	return db
}
