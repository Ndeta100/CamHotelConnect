package main

import (
	"context"
	"flag"
	"github.com/Ndeta100/CamHotelConnect/api"
	"github.com/Ndeta100/CamHotelConnect/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollection = "users"

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {

	listenAddr := flag.String("listenAddr", ":5000", "The address of the api server")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	//handlers initialization
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler   = api.NewUserHandler(userStore)
		hotelHandler  = api.NewHotelHandler(store)
		authHandler   = api.NewAuthHandler(userStore)
		roomHandler   = api.NewRoomHandler(store)
		bookingHander = api.NewBookingHandler(store)
		app           = fiber.New(config)
		auth          = app.Group("/api")
		apiV1         = app.Group("api/v1", api.JWTAuthentication(userStore))
		admin         = apiV1.Group("/admin", api.AdminAuth)
	)

	//auth handle
	auth.Post("/auth", authHandler.HandleAuth)
	//Versioned api routes
	//User handlers
	apiV1.Post("user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Put("user/:id", userHandler.HandlePutUser)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	//Hotel handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("hotel/:id/rooms", hotelHandler.HandleGetRooms)

	//room handlers
	apiV1.Get("/room", roomHandler.HandleGetRooms)
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	//TODO: cancel a booking
	//booking handlers
	apiV1.Get("/booking/:id", bookingHander.HandleGetBooking)
	apiV1.Get("booking/:id/cancel", bookingHander.HandleCancelBooking)
	//admin handlers
	admin.Get("/booking", bookingHander.HandleGetBookings)
	err = app.Listen(*listenAddr)
	if err != nil {
		return
	}
}
