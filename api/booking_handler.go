package api

import (
	"github.com/Ndeta100/CamHotelConnect/db"
	_ "github.com/Ndeta100/CamHotelConnect/types"
	"github.com/Ndeta100/CamHotelConnect/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// HandleGetBookings TODO: This needs to be admin authorized
//
//	@Summary		Get all bookings
//	@Description	Get a list of all bookings.
//	@Tags			bookings
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		types.Booking	"List of bookings retrieved successfully"
//	@Failure		404	{object}	GenericResp		"Bookings not found"
//	@Failure		500	{object}	GenericResp		"Internal Server Error"
//	@Router			/bookings [get]
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return ErrResourceNotFound("bookings")
	}
	return c.JSON(bookings)
}

// HandleGetBooking TODO: This needs to be user authorized
//
//	@Summary		Get a booking
//	@Description	Get details of a booking by its ID. Only the user who made the booking and hotel owner can retrieve it.
//	@Tags			bookings
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Booking ID"
//	@Success		200	{object}	types.Booking	"Booking details retrieved successfully"
//	@Failure		404	{object}	GenericResp		"Booking not found"
//	@Failure		401	{object}	GenericResp		"Unauthorized"
//	@Failure		500	{object}	GenericResp		"Internal Server Error"
//	@Router			/bookings/{id} [get]
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return ErrUnauthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}
	return c.JSON(booking)
}

// HandleCancelBooking
//
//	@Summary		Cancel a booking
//	@Description	Cancel a booking by its ID. Only the user who made the booking can cancel it.
//	@Tags			bookings
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Booking ID"
//	@Success		200	{object}	GenericResp	"Booking canceled successfully"
//	@Failure		404	{object}	GenericResp	"Booking not found"
//	@Failure		401	{object}	GenericResp	"Unauthorized"
//	@Failure		500	{object}	GenericResp	"Internal Server Error"
//	@Router			/bookings/{id}/cancel [put]
func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return ErrUnauthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(GenericResp{
		Type: "msg",
		Msg:  "booking canceled",
	})
}
