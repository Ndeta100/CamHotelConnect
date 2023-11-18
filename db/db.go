package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const DBNAME = "hotel-reservation"

func toObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}
