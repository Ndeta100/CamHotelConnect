package api

import (
	"errors"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/db"
	"github.com/Ndeta100/CamHotelConnect/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"time"
)

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type AuthHandler struct {
	userStore db.UserStore
}

type GenericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(GenericResp{
		Type: "error",
		Msg:  "Invalid credentials",
	})
}

// HandleAuth ----------------------GENERAL OVERVIEW
// HandleAuth Handler should only do:
//   - Serialization of the incoming request(JSON)
//    - do some data fetching from db
//    - call some business logic
//    - return data back to user
func (h *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}

	}
	if !types.IsValidPassword(user.EncryptedPassword, authParams.Password) {
		return invalidCredentials(c)

	}
	resp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}
	return c.JSON(resp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr
}
