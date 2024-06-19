package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/db/fixtures"
	"github.com/Ndeta100/CamHotelConnect/types"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	//insertedUser := insertTestUser(t, tdb.store.User)
	insertedUser := fixtures.AddUser(*tdb.store, "james", "foo", types.UserRoleAdmin)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuth)
	params := AuthParams{
		Email:    insertedUser.Email,
		Password: "james_foo",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected http status of 200 but got %d", resp.StatusCode)
	}
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Error(err)
	}
	if authResp.Token == "" {
		t.Fatalf("Expeected JWT token and auth response")
	}
	//Set the encrypted password to an empty string, because we do not return that in any
	//json response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		fmt.Println("Inserted------------", insertedUser)
		fmt.Println("AUTH-------------", authResp.User)
		t.Fatalf("expected the user to be the inserted user")
	}
	fmt.Println(resp)
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	wrongPassword := fixtures.AddUser(*tdb.store, "james", "foo", types.UserRoleUser)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuth)
	params := AuthParams{
		Email:    wrongPassword.Email,
		Password: "notme",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected http status of 400 but got %d", resp.StatusCode)
	}
	var genResp GenericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatalf("Expect generic response of type error, but got %s", genResp.Type)
	}
	if genResp.Msg != "Invalid credentials" {
		t.Fatalf("expected generic response of type <invalid credentials> but got %s", genResp.Msg)
	}
}
