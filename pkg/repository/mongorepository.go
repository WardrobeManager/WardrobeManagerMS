//
// mongorepository.go
//
// May 2021, Prashant Desai
//

package repository

import (
	"context"
	"fmt"

	"WardrobeManagerMS/pkg/api"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const Version = "1.0"

const (
	CONNECTIONSTRING = "mongodb://localhost:27017"
	DB               = "wardrobemanager"
	WARDS            = "wardrobe"
)

type mongoWardRepo struct {
	collection *mongo.Collection
}

func NewWardrobeRepository() (api.WardrobeRepository, error) {

	fmt.Println("Initializing Mongo User Store")
	clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	newCollection := client.Database(DB).Collection(WARDS)

	// set username as index
	_, err = newCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, err
	}
	wardRepo := &mongoWardRepo{
		collection: newCollection,
	}

	return wardRepo, nil
}

func (m *mongoWardRepo) Add(user string, wards *api.WardrobeCloset) error {

	_, err := m.collection.InsertOne(context.TODO(), wards)
	if err != nil {
		return fmt.Errorf("Error adding user %s : %w", user, err)
	}

	return nil
}

func (m *mongoWardRepo) Get(user string) (*api.WardrobeCloset, error) {
	filter := bson.M{"username": user}

	var wardCloset api.WardrobeCloset

	err := m.collection.FindOne(context.TODO(), filter).Decode(&wardCloset)
	if err != nil {
		return nil, fmt.Errorf("User %s not present", user)
	}

	return &wardCloset, nil

	return &api.WardrobeCloset{}, nil
}

func (m *mongoWardRepo) Update(user string, wards *api.WardrobeCloset) error {
	return nil
}

func (m *mongoWardRepo) Delete(user string) error {
	return nil
}
