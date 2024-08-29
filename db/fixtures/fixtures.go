package fixtures

import (
	"context"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/db"
	"github.com/Ndeta100/CamHotelConnect/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func AddUser(store db.Store, fname, lname string, userRole types.UserRole) *types.User {
	user, err := types.NewUserFromParams(
		types.CreateUserParams{
			Email:     fmt.Sprintf("%s@%s.com", fname, lname),
			FirstName: fname,
			LastName:  lname,
			Password:  fmt.Sprintf("%s_%s", fname, lname),
		})

	if err != nil {
		log.Fatal(err)
	}
	user.Role = userRole
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddHotel(store *db.Store, name string, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	roomIDS := rooms
	if roomIDS == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIDS,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddRoom(store *db.Store, size string, seaSide bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaSide,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.InsertRoom(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(store db.Store, roomID, userID primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
