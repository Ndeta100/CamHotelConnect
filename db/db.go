package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
)

var DBNAME = os.Getenv("MONGO_DB_NAME")

func toObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}

type Pagination struct {
	Limit int64
	Page  int64
}

type HotelFilter struct {
	Pagination
	Rating int `json:"rating"`
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
