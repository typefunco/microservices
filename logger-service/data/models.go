package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

// Models wraps the LogEntry model to interact with MongoDB collections.
type Models struct {
	LogEntry LogEntry
}

// LogEntry represents a log entry in the MongoDB collection.
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// New initializes a new instance of Models with the provided MongoDB client.
func New(client *mongo.Client) Models {
	MongoClient = client

	return Models{
		LogEntry: LogEntry{},
	}
}

// Insert inserts a new log entry into the MongoDB collection.
// The CreatedAt and UpdatedAt fields are automatically set to the current time.
func (l *LogEntry) Insert(entry LogEntry) error {
	collection := MongoClient.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.Background(), LogEntry{
		ID:        entry.ID,
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Can't insert data into MongoDB\nError:", err)
		return err
	}

	return nil
}

// GetAll retrieves all log entries from the MongoDB collection.
// It returns a slice of LogEntry and an error if any occurs.
func (l *LogEntry) GetAll() ([]LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("logs").Collection("logs")

	filter := bson.D{}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []LogEntry
	for cursor.Next(ctx) {
		var logRecord LogEntry
		err = cursor.Decode(&logRecord)
		if err != nil {
			return nil, err
		}
		logs = append(logs, logRecord)
	}

	return logs, nil
}

// GetById retrieves a log entry by its ID from the MongoDB collection.
// It returns a pointer to LogEntry and an error if any occurs.
func (l *LogEntry) GetById(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// DropCollection drops the entire MongoDB collection.
// It deletes all documents and returns an error if any occurs.
func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("logs").Collection("logs")

	err := collection.Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing log entry in the MongoDB collection based on its ID.
// It returns an UpdateResult and an error if any occurs.
func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
