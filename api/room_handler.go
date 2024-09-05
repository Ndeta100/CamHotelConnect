package api

import (
	"context"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/db"
	"github.com/Ndeta100/CamHotelConnect/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(time.Time(p.FromDate)) || now.After(time.Time(p.TillDate)) {
		return fmt.Errorf("cannot book a room in the past")
	}
	return nil
}

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

// HandleBookRoom TODO:
//
//	@Summary		Book a room
//	@Description	Book a room with the specified parameters. The room must be available.
//	@Tags			rooms
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Room ID"
//	@Param			booking	body		BookRoomParams	true	"Booking parameters"
//	@Success		200		{object}	types.Booking	"Room booked successfully"
//	@Failure		400		{object}	GenericResp		"Bad Request - Room already booked or invalid parameters"
//	@Failure		500		{object}	GenericResp		"Internal Server Error"
//	@Router			/rooms/{id}/book [post]
func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomID := c.Params("id")
	roomOid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(GenericResp{
			Type: "error",
			Msg:  "Internal server error",
		})
	}
	ok, err = h.isRoomAvailable(c.Context(), roomOid, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(GenericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked", c.Params("id")),
		})
	}
	booking := types.Booking{
		RoomID:    roomOid,
		UserID:    user.ID,
		FromDate:  params.FromDate,
		TillDate:  params.TillDate,
		NumPerson: params.NumPersons,
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

// HandleGetRooms TODO
//
//	@Summary		Get all rooms
//	@Description	Get a list of all rooms available.
//	@Tags			rooms
//	@Produce		json
//	@Success		200	{array}		types.Room	"List of rooms retrieved successfully"
//	@Failure		500	{object}	GenericResp	"Internal Server Error"
//	@Router			/rooms [get]
func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) isRoomAvailable(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}
