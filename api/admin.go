package api

import (
	"github.com/Ndeta100/CamHotelConnect/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnauthorized()
	}
	if user.Role != types.UserRoleAdmin {
		return ErrUnauthorized()
	}
	return c.Next()
}
