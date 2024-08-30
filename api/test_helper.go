package api

import (
	"context"
	"github.com/Ndeta100/CamHotelConnect/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
)

type testdb struct {
	client *mongo.Client
	store  *db.Store
}

func (tdb testdb) teardown(t *testing.T) {
	dbname := os.Getenv("MONGO_DB_NAME")
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	dburi := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,
		store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
			Hotel:   hotelStore,
		},
	}
}

func init() {
	//if err := godotenv.Load("../.env"); err != nil {
	//	log.Fatal(err)
	//}
}
