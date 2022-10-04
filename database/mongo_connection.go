package database

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

const timeout = 10 * time.Second

func NewClient() *mongo.Client {

	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	uri := os.Getenv("MONGO_URI")

	opts := options.Client().ApplyURI(uri)

	client, err := mongo.NewClient(opts)
	if err != nil {
		logrus.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		logrus.Fatal(err)
	}

	return client
}

var Client *mongo.Client = NewClient()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(os.Getenv("DB_NAME")).Collection(collectionName)
	return collection
}
