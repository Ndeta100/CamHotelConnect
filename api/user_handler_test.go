package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Ndeta100/CamHotelConnect/types"
	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	app := fiber.New()
	userHandler := NewUserHandler(tdb.store.User)
	app.Post("/", userHandler.HandlePostUser)
	params := types.CreateUserParams{
		Email:     "ndeta@gmail.com",
		FirstName: "Ndeta",
		LastName:  "Innocent",
		Password:  "URUHRhai938022",
	}
	//marshall params to bytes
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return
	}
	if user.ID == primitive.NilObjectID {
		t.Errorf("expecting a user id")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Not expecting encrypted password in a json res")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
	bb, _ := io.ReadAll(resp.Body)
	//Logging to test results
	fmt.Println(bb, "byte body")
	fmt.Println(user, "user")
	fmt.Println(resp.Status)
}
