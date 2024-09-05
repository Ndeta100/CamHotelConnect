package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string               `bson:"name" json:"name"`
	Location    string               `bson:"location" json:"location"`
	Description string               `bson:"description" json:"description"`
	Amenities   []string             `bson:"amenities" json:"amenities"`
	Rooms       []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating      int                  `bson:"rating" json:"rating"`
	Images      []string             `bson:"images" json:"images"`
	UserId      primitive.ObjectID   `bson:"user_id" json:"user_id"`
}

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Seaside bool               `bson:"seaside" json:"seaside"`
	//small, normal or kingSize
	Size    string             `bson:"size" json:"size"`
	Price   float64            `bson:"price" json:"price"`
	Images  []string           `bson:"images" json:"images"`
	HotelID primitive.ObjectID ` bson:"hotelID" json:"hotelID"`
}

type Image struct {
	ID  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	URL string             `bson:"url" json:"url"`
}
