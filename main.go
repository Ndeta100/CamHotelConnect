package main

import (
	"flag"
	"github.com/Ndeta100/CamHotelConnect/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The address of the api server")
	app := fiber.New()
	apiV1 := app.Group("api/v1")
	app.Get("/foo", handleFoo)
	apiV1.Get("/user", api.HandleGetUsers)
	apiV1.Get("/user/:id", api.HandleGetUser)
	err := app.Listen(*listenAddr)
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
