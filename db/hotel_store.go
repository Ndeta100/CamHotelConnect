package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type HotelStore interface {
	GetHotels(ctx context.Context, filter bson.M, options *HotelFilter) ([]*types.Hotel, error)
	InsertHotel(ctx context.Context, userID primitive.ObjectID, hotel *types.Hotel) (*types.Hotel, error)
	UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error
	GetHotelByID(ctx context.Context, userID, hotelID primitive.ObjectID) (*types.Hotel, error)
	DeleteHotelByID(ctx context.Context, id primitive.ObjectID) error
	HotelExists(ctx context.Context, hotel *types.Hotel) (bool, error)
}
type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	dbname := os.Getenv("MONGO_DB_NAME")
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection("hotels"),
	}
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, userID primitive.ObjectID, hotel *types.Hotel) (*types.Hotel, error) {
	exist, err := s.HotelExists(ctx, hotel)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("hotel already exists")
	}
	//add the user to the hotel
	hotel.UserId = userID
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M, page *HotelFilter) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip((page.Page - 1) * page.Limit)
	opts.SetLimit(page.Limit)
	resp, err := s.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, userID, hotelID primitive.ObjectID) (*types.Hotel, error) {
	var hotel types.Hotel
	filter := bson.M{
		"_id":     hotelID,
		"user_id": userID,
	}
	//try finding hotel
	err := s.coll.FindOne(ctx, filter).Decode(&hotel)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("hotel not found")
		}
		return nil, fmt.Errorf("failed to find hotel")
	}
	return &hotel, nil
}

func (s *MongoHotelStore) DeleteHotelByID(ctx context.Context, id primitive.ObjectID) error {
	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete hotel with ID %v: %w", id, err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf(" hotel with ID %v does not exist", id)
	}
	return nil
}

func (s *MongoHotelStore) HotelExists(ctx context.Context, hotel *types.Hotel) (bool, error) {
	filter := bson.M{"$or": []bson.M{
		{"name": hotel.Name, "location": hotel.Location},
	}}
	count, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	log.Println("Hotel exists:", count)
	return count > 0, nil
}
