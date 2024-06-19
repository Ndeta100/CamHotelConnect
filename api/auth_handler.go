package api

import (
	"errors"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/db"
	"github.com/Ndeta100/CamHotelConnect/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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
		log.Printf("Failed to parse request body: %v", err)
		return err
	}
	log.Printf("Received login request for email: %s", authParams.Email)
	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Invalid login attempt: email not found %s", authParams.Email)
			return invalidCredentials(c)
		}

	}
	if !types.IsValidPassword(user.EncryptedPassword, authParams.Password) {
		log.Printf("Invalid login attempt:Passwords do not match %s", authParams.Password)
		return invalidCredentials(c)

	}
	token := CreateTokenFromUser(user)
	// Set the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(4 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // Set to true in production
		SameSite: "Strict",
	})
	resp := AuthResponse{
		User:  user,
		Token: token,
	}
	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
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

// Add this function to handle the validate_token endpoint
func (h *AuthHandler) HandleValidateToken(c *fiber.Ctx) error {
	tokenStr := c.Cookies("access_token") // Assuming the token is sent in a cookie
	if tokenStr == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token provided",
		})
	}

	claims, err := validateToken(tokenStr)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}
	// If the token is valid, optionally return user details or just a success message
	userID := claims["id"].(string)
	user, err := h.userStore.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.JSON(fiber.Map{
		"user":    user,
		"message": "Token is valid",
	})
}
