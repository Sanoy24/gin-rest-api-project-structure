package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri, dbName string, timeout time.Duration) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	// ping db
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	if err := createIndexes(ctx, db); err != nil {
		return nil, err
	}
	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	userCollection := db.Collection("users")

	// Create unique index on email
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// Create index on username
	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// Create indexes on created_at for sorting
	createdAtIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}},
	}

	_, err := userCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		emailIndex,
		usernameIndex,
		createdAtIndex,
	})

	return err
}

// func createIndex(ctx context.Context, db *mongo.Database, collectioName, key string) error {
// 	// Example index creation, adjust as needed
// 	collection := db.Collection(collectioName)

// 	index := mongo.IndexModel{
// 		Keys:    bson.D{{Key: key, Value: 1}},
// 		Options: options.Index().SetUnique(true),
// 	}
// 	_, err := collection.Indexes().CreateOne(ctx, index)
// 	return err
// }

func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
