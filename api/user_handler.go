package api

import (
	"errors"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/db"
	"github.com/Ndeta100/CamHotelConnect/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

// HandlePostUser TODO: send email to newly created usr
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided details.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		types.CreateUserParams	true	"User data"
//	@Success		200		{object}	types.User				"User created successfully"
//	@Failure		400		{object}	GenericResp				"Bad Request"
//	@Failure		409		{string}	string					"User already exists"
//	@Failure		500		{object}	GenericResp				"Internal Server Error"
//	@Router			/users [post]
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errorS := params.Validate(); len(errorS) > 0 {
		return c.JSON(errorS)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(fiber.StatusConflict).JSON("User already exists")
		}
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	return c.JSON(insertedUser)
}

// HandlePutUser TODO:check if user exist before updating
//
//	@Summary		Update a user
//	@Description	Update user details by user ID.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"User ID"
//	@Param			user	body		types.UpdateUserParams	true	"Updated user data"
//	@Success		200		{object}	map[string]string		"User updated successfully"
//	@Failure		400		{object}	GenericResp				"Invalid User ID"
//	@Failure		404		{object}	GenericResp				"User not found"
//	@Failure		500		{object}	GenericResp				"Internal Server Error"
//	@Router			/users/{id} [put]
func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	//values := bson.M{}
	params := types.UpdateUserParams{}
	userID := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidId()
	}
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": userID})
}

// HandleDeleteUser TODO: a lot to do like still to decide how delete will work in the app
//
//	@Summary		Delete a user
//	@Description	Delete a user by their ID.
//	@Tags			users
//	@Produce		json
//	@Param			id	path		string				true	"User ID"
//	@Success		200	{object}	map[string]string	"User deleted successfully"
//	@Failure		404	{object}	GenericResp			"User not found"
//	@Failure		500	{object}	GenericResp			"Internal Server Error"
//	@Router			/users/{id} [delete]
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": userID})
}

// HandleGetUser TODO:checks to be made no idea yet
//
//	@Summary		Get a user
//	@Description	Get a user by their ID.
//	@Tags			users
//	@Produce		json
//	@Param			id	path		string				true	"User ID"
//	@Success		200	{object}	types.User			"User details retrieved successfully"
//	@Failure		404	{object}	map[string]string	"User not found"
//	@Failure		500	{object}	GenericResp			"Internal Server Error"
//	@Router			/users/{id} [get]
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// Convert userID to ObjectID if necessary
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Invalid user ID format in token:", id)
		return ErrUnauthorized()
	}
	user, err := h.userStore.GetUserByID(c.Context(), objectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"err": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

// HandleGetUsers TODO: checks to be made still
//
//	@Summary		Get all users
//	@Description	Get a list of all users.
//	@Tags			users
//	@Produce		json
//	@Success		200	{array}		types.User	"List of users retrieved successfully"
//	@Failure		404	{object}	GenericResp	"Users not found"
//	@Failure		500	{object}	GenericResp	"Internal Server Error"
//	@Router			/users [get]
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return ErrResourceNotFound("user")
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUserByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return ErrBadRequest()
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"err": "not found"})
		}
		return err
	}
	return c.JSON(user)
}
