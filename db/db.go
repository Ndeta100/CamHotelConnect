package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	DBNAME   = "hotel-reservation"
	DBURI    = "mongodb://localhost:27017"
	TestNAME = "hotel-reservation-test"
)

func toObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
