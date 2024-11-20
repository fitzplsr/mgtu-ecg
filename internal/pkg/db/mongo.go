package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoWrapper struct {
	Client *mongo.Client
	DB     string
}

func ConnectMongoDB(config Config) (*MongoWrapper, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port))

	credential := options.Credential{
		AuthSource: config.DB,
		Username:   config.User,
		Password:   config.Password,
	}
	clientOptions.SetAuth(credential)

	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connecting to MongoDB: %w", err)

	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ping MongoDB: %v", err)
	}

	return &MongoWrapper{
		Client: client,
		DB:     config.DB,
	}, nil
}
