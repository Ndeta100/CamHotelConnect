package main

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	mongoURI := os.Getenv("MONGO_DB_URL")
	dbName := os.Getenv("MONGO_DB_NAME")

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database(dbName)

	collections, err := database.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	for _, collection := range collections {
		if err := database.Collection(collection).Drop(ctx); err != nil {
			log.Printf("Failed to drop collection %s: %v", collection, err)
		} else {
			log.Printf("Dropped collection %s", collection)
		}
	}

	log.Println("Database cleaning completed.")
}
