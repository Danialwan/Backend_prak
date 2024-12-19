package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB establishes a connection to MongoDB
func ConnectDB() {
	log.Println("Connecting to MongoDB...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB Atlas
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://alwanhakimramadhani:danialwan1.@cluster0.vnrtv.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Connect to the specific database
	DB = client.Database("SSOUNAIRSATU")
	log.Println("Connected to MongoDB!")
}

// GetCollection retrieves a specific collection from the database
func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		ConnectDB()
	}
	return DB.Collection(collectionName)
}
