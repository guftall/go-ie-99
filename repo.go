package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func initializeDatabase() {
	uri := os.Getenv("mongouri")
	log.Printf("mongodb uri: %s", uri)
	_client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	_ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = _client.Connect(_ctx)
	if err != nil {
		log.Fatal(err)
	}
	client = _client
}

func readIdentifier(key string) string {

	collection := client.Database("ieproj").Collection("rsakeys")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	type doc struct {
		Key        string
		Identifier string
	}

	var _doc doc
	err := collection.FindOne(ctx, bson.D{{"pk", key}}).Decode(&_doc)

	if err != nil {
		return ""
	}

	return _doc.Identifier
}

func isIdentifierExist(identifier string) bool {

	collection := client.Database("ieproj").Collection("rsakeys")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.D{{"identifier", identifier}})

	if err != nil {
		return false
	}

	return count > 0
}

func insertPublicKey(key, identifier string) {

	collection := client.Database("ieproj").Collection("rsakeys")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, bson.D{{"pk", key}, {"identifier", identifier}})
	if err != nil {
		log.Print(err)
		return
	}

	log.Println("inserted")
	log.Print(res.InsertedID)
}

func countRecords() int {

	collection := client.Database("ieproj").Collection("rsakeys")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Print(err)
		return -1
	}

	return int(count)

}
