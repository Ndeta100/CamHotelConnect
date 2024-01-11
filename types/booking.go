package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID    primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NumPerson int                `bson:"numPerson" json:"numPerson"`
	FromDate  time.Time          `bson:"fromDate" json:"fromDate"`
	TillDate  time.Time          `bson:"tillDate" json:"tillDate"`
	Canceled  bool               `bson:"canceled" json:"canceled"`
}
