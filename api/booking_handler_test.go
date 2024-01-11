package api

import (
	"encoding/json"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/db/fixtures"
	"github.com/Ndeta100/CamHotelConnect/types"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAdminGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	user := fixtures.AddUser(*db.store, "james", "foo", false)
	hotel := fixtures.AddHotel(db.store, "bar hotel", "arizona", 4, nil)
	room := fixtures.AddRoom(db.store, "small", true, 4.4, hotel.ID)
	from := time.Now()
	till := from.AddDate(0, 0, 4)
	booking := fixtures.AddBooking(*db.store, user.ID, room.ID, from, till)
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	admin := app.Group("/", JWTAuthentication(db.store.User), AdminAuth)
	bookingHandler := NewBookingHandler(db.store)
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	adminUser := fixtures.AddUser(*db.store, "admin", "admin", true)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("not a 200 response %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking but got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
	}
	fmt.Println(bookings)

	//test a non-admin user (cannot access booking)
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected a  status unauthorize but got %d", resp.StatusCode)
	}

}

//func TestUserGetBooking(t *testing.T) {
//	db := setup(t)
//	defer db.teardown(t)
//	user := fixtures.AddUser(*db.store, "james", "foo", false)
//	hotel := fixtures.AddHotel(db.store, "bar hotel", "arizona", 4, nil)
//	room := fixtures.AddRoom(db.store, "small", true, 4.4, hotel.ID)
//	from := time.Now()
//	till := from.AddDate(0, 0, 4)
//	booking := fixtures.AddBooking(*db.store, user.ID, room.ID, from, till)
//	app := fiber.New()
//	jwtapproute := app.Group("/", middleware.JWTAuthentication(db.store.User))
//	bookingHandler := NewBookingHandler(db.store)
//	jwtapproute.Get("/:id", bookingHandler.HandleGetBooking)
//	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
//	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
//	resp, err := app.Test(req)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if resp.StatusCode != http.StatusOK {
//		t.Fatalf("non 200 code got %d", resp.StatusCode)
//	}
//	var bookingRes *types.Booking
//	if err := json.NewDecoder(resp.Body).Decode(&bookingRes); err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(bookingRes)
//}
