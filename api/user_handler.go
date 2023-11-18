package api

import (
	"github.com/Ndeta100/CamHotelConnect/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Ndeta",
		LastName:  "Innocent",
	}
	return c.JSON(u)
}
func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("Ndeta")
}
