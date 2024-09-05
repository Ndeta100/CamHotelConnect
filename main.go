package main

import (
	"context"
	"github.com/Ndeta100/CamHotelConnect/api"
	"github.com/Ndeta100/CamHotelConnect/db"
	_ "github.com/Ndeta100/CamHotelConnect/docs"
	"github.com/Ndeta100/CamHotelConnect/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

const userCollection = "users"

// Configurations
// 1. MongoDB endpoint
// 2. Listen address of http server
// 3. JWT secret
// 4. MongoDBNAME
var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

// @title			CAM HOTEL CONNECT
// @version		0.0.1 beta
// @description	This project is a backend JSON API for a hotel reservation system.
// @contact.name	API Support
// @contact.email	api.support.huz@mail.com
// @license.name	Apache 2.0
func main() {
	mongoUri := os.Getenv("MONGO_DB_URL")
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	cloudKey := os.Getenv("CLOUDINARY_API_KEY")
	cloudSecret := os.Getenv("CLOUDINARY_API_SECRET")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if mongoUri == "" || cloudName == "" || cloudKey == "" || cloudSecret == "" {
		log.Fatal("Error loading environment variables")
	}

	imageUploader, err := utils.NewCloudinaryUploader(cloudName, cloudKey, cloudSecret)
	if err != nil {
		log.Fatalf("error initialing cloudinary uploader %v", err)
	}
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
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store, imageUploader)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiV1          = app.Group("api/v1")
		admin          = apiV1.Group("/admin", api.AdminAuth)
	)
	// Use CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Change this to restrict to specific origins
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Content-Type, X-Api-Token",
		AllowCredentials: true,
	}))
	// Define unauthenticated routes
	auth.Post("/auth", authHandler.HandleAuth)
	apiV1.Post("/user", userHandler.HandlePostUser)
	//Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)
	//Hotel public routes
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Use the JWT authentication middleware
	apiV1.Use(api.JWTAuthentication(userStore))
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)

	//Versioned api routes
	//User handlers
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Put("user/:id", userHandler.HandlePutUser)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	//Hotel handlers
	apiV1.Post("/hotel", hotelHandler.HandleAddHotel)

	//room handlers
	apiV1.Get("/room", roomHandler.HandleGetRooms)
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	//TODO: cancel a booking
	//booking handlers
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Get("booking/:id/cancel", bookingHandler.HandleCancelBooking)
	//admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	err = app.Listen(listenAddr)
	if err != nil {
		return
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
