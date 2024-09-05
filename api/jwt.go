package api

import (
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"os"
	"time"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Println("X-Api-Token not present in header")
			return ErrUnauthorized()
		}
		fmt.Println("Token received:", token[0])
		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		//check token expiration
		if time.Now().Unix() > expires {
			// TODO createRefreshToken
			log.Println("Token expired")
			return NewError(http.StatusUnauthorized, "Expired token")
		}
		userID := claims["id"].(string)
		// Convert userID to ObjectID if necessary
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			fmt.Println("Invalid user ID format in token:", objectID)
			return ErrUnauthorized()
		}
		user, err := userStore.GetUserByID(c.Context(), objectID)
		if err != nil {
			return ErrUnauthorized()
		}
		// Ensure user is not nil
		if user == nil {
			fmt.Println("User not found in database")
			return ErrUnauthorized()
		}
		//set authenticated user to the  context
		//c.Context().SetUserValue("user", user)
		c.Locals("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method", token.Header["alg"])
			return nil, ErrUnauthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Failed to parse jwt token:", err)
		return nil, ErrUnauthorized()
	}
	if !token.Valid {
		fmt.Println("invalid")
		return nil, ErrUnauthorized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnauthorized()
	}
	return claims, nil
}
