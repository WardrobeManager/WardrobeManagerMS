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
	DB               = "wardrobemanager"
	WARDS            = "wardrobe"
)

type mongoWardRepo struct {
	collection *mongo.Collection
}

func NewWardrobeRepository(server string) (api.WardrobeRepository, error) {

	fmt.Println("Initializing Mongo User Store")
	clientOptions := options.Client().ApplyURI("mongodb://" + server)

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
			Keys:    bson.D{{Key: "user", Value: 1}},
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
	fmt.Println("FUNC START : Add")
	_, err := m.collection.InsertOne(context.TODO(), wards)
	if err != nil {
		return fmt.Errorf("Error adding user %s : %w", user, err)
	}

	return nil
}

func (m *mongoWardRepo) Get(user string) (*api.WardrobeCloset, error) {
	filter := bson.M{"user": user}

	var wardCloset api.WardrobeCloset

	err := m.collection.FindOne(context.TODO(), filter).Decode(&wardCloset)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, &api.UserNotFound{User: user}
		}

		return nil, fmt.Errorf("Internal MongoDB error : %w", err)
	}

	return &wardCloset, nil

}

func (m *mongoWardRepo) Update(user string, wards *api.WardrobeCloset) error {
	filter := bson.M{"user": user}

	_, err := m.collection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: wards}})
	if err != nil {
		return fmt.Errorf("Error updating wardrobe closet record for user %s : %w", user, err)
	}

	return nil
}

func (m *mongoWardRepo) DeleteAll(user string) error {
	filter := bson.M{"user": user}

	_, err := m.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("Error deleting user %s wardrobe closet : %w", user, err)
	}

	return nil
}
