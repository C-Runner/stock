package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI     string
	MongoDBName  string
	ServerPort   string
}

var MongoClient *mongo.Client
var DB *mongo.Database

func LoadConfig() *Config {
	return &Config{
		MongoURI:    "mongodb://admin:OpenClaw2026!@43.139.206.109:27017",
		MongoDBName: "stock",
		ServerPort:  "8080",
	}
}

func ConnectMongoDB(cfg *Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	MongoClient = client
	DB = client.Database(cfg.MongoDBName)
	log.Println("Connected to MongoDB successfully!")
	return nil
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}

func DisconnectMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		MongoClient.Disconnect(ctx)
		log.Println("Disconnected from MongoDB")
	}
}
