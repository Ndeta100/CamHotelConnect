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

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userCollection = "users"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":5000", "The address of the api server")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	//handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	app := fiber.New(config)
	apiV1 := app.Group("api/v1")
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	err = app.Listen(*listenAddr)
	if err != nil {
		return
	}
}
func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Working just find"})
}
func handleUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"me": "Ndeta"})
}
