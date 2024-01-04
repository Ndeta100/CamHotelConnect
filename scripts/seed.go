package main

import (
	"context"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/api"
	"github.com/Ndeta100/CamHotelConnect/db"
	"github.com/Ndeta100/CamHotelConnect/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}
	user := fixtures.AddUser(store, "james", "foo", false)
	admin := fixtures.AddUser(store, "ndeta", "inno", true)
	fmt.Println("ADMIN-------->", api.CreateTokenFromUser(admin))
	fmt.Println("USER-------->", api.CreateTokenFromUser(user))
	hotel := fixtures.AddHotel(&store, "Hilton", "Buea", 4, nil)
	room := fixtures.AddRoom(&store, "large", true, 65.4, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("BOOKING-------->", booking.ID)

}
