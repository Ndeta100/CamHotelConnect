package main

import (
	"context"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/api"
	"github.com/Ndeta100/CamHotelConnect/db"
	"github.com/Ndeta100/CamHotelConnect/db/fixtures"
	"github.com/Ndeta100/CamHotelConnect/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatal(err)
	//}
	mongo_uri := os.Getenv("MONGO_DB_URL")
	dbname := os.Getenv("MONGO_DB_NAME")
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongo_uri))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(dbname).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}
	user := fixtures.AddUser(store, "james", "foo", types.UserRoleUser)
	admin := fixtures.AddUser(store, "ndeta", "inno", types.UserRoleUser)
	fmt.Println("ADMIN-------->", api.CreateTokenFromUser(admin))
	fmt.Println("USER-------->", api.CreateTokenFromUser(user))
	hotel := fixtures.AddHotel(&store, "Hilton", "Buea", 4, nil)
	room := fixtures.AddRoom(&store, "large", true, 65.4, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("BOOKING-------->", booking.ID)
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(&store, name, location, rand.Intn(5)+1, nil)
	}
}
