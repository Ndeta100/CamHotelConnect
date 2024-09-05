package api

import (
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/db"
	"github.com/Ndeta100/CamHotelConnect/types"
	"github.com/Ndeta100/CamHotelConnect/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type HotelHandler struct {
	//it depends on a room store, and it also depends on a hotel store
	store         *db.Store
	ImageUploader utils.ImageUploader
}

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

func NewHotelHandler(store *db.Store, imageUploader utils.ImageUploader) *HotelHandler {
	return &HotelHandler{
		store:         store,
		ImageUploader: imageUploader,
	}
}

// HandleGetHotels TODO
//
//	@Summary		Get all hotels
//	@Description	Get a list of all hotels with optional filtering and pagination.
//	@Tags			hotels
//	@Accept			json
//	@Produce		json
//	@Param			rating	query		int				false	"Filter hotels by rating"
//	@Param			page	query		int				false	"Page number for pagination"	default(1)
//	@Param			limit	query		int				false	"Number of hotels per page"		default(10)
//	@Success		200		{object}	ResourceResp	"List of hotels retrieved successfully"
//	@Failure		400		{object}	GenericResp		"Bad Request"
//	@Failure		404		{object}	GenericResp		"Hotels not found"
//	@Failure		500		{object}	GenericResp		"Internal Server Error"
//	@Router			/hotels [get]
func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var queryFilter db.HotelFilter
	if err := c.QueryParser(&queryFilter); err != nil {
		return ErrBadRequest()
	}
	filter := bson.M{
		//"rating": queryFilter.Rating,
	}
	if queryFilter.Page < 1 {
		queryFilter.Page = 1
	}
	if queryFilter.Limit < 1 {
		queryFilter.Limit = 10
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &queryFilter)
	if err != nil {
		return ErrResourceNotFound("hotels")
	}
	resp := ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(queryFilter.Page),
	}
	return c.JSON(resp)
}

// HandleGetRooms TODO
//
//	@Summary		Get all rooms in a hotel
//	@Description	Get a list of all rooms for a specific hotel.
//	@Tags			hotels
//	@Produce		json
//	@Param			id	path		string		true	"Hotel ID"
//	@Success		200	{array}		types.Room	"List of rooms retrieved successfully"
//	@Failure		400	{object}	GenericResp	"Invalid Hotel ID"
//	@Failure		404	{object}	GenericResp	"Hotel not found"
//	@Failure		500	{object}	GenericResp	"Internal Server Error"
//	@Router			/hotels/{id}/rooms [get]
func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidId()
	}
	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(rooms)
}

// HandleGetHotel TODO:
//
//	@Summary		Get a specific hotel
//	@Description	Get details of a specific hotel by its ID.
//	@Tags			hotels
//	@Produce		json
//	@Param			id	path		string		true	"Hotel ID"
//	@Success		200	{object}	types.Hotel	"Hotel details retrieved successfully"
//	@Failure		400	{object}	GenericResp	"Invalid Hotel ID"
//	@Failure		404	{object}	GenericResp	"Hotel not found"
//	@Failure		500	{object}	GenericResp	"Internal Server Error"
//	@Router			/hotels/{id} [get]
func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotelID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	userInterface := c.Locals("user")
	fmt.Println("ID", userInterface)
	user, ok := userInterface.(*types.User)
	if !ok || user == nil {
		fmt.Printf("Invalid user type or nil user in context: %v\n", userInterface)
		return ErrUnauthorized()
	}
	if !userIsAdmin(user) {
		return fmt.Errorf("user is not admin")
	}
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), user.ID, hotelID)
	if err != nil {
		//Check if hotel not found
		if err.Error() == "hotel not found" {
			//	TODO: This error occurs when user is trying to get a hotel they did not add
			return c.Status(fiber.StatusConflict).JSON("hotel does not exists")
		}
		return c.Status(fiber.StatusInternalServerError).JSON(GenericResp{Msg: "Internal Server Error"})
	}
	// Check if the user is authorized (either the owner or an admin)
	if hotel.UserId != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(GenericResp{Msg: "You do not have permission to access this resource"})
	}
	return c.JSON(hotel)
}

// HandleAddHotel CreateHotel handles the creation of a new hotel
//
//	@Summary		Add a new hotel
//	@Description	Add a new hotel. Only users with the "ADMIN" role can create hotels.
//	@Tags			hotels
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			user	body		types.Hotel	true	"Hotel data"
//	@Param			images	formData	file		true	"Hotel images"
//	@Success		201		{object}	types.Hotel	"Hotel created successfully"
//	@Failure		400		{object}	GenericResp	"Bad Request"
//	@Failure		401		{object}	GenericResp	"Unauthorized - Missing or invalid credentials"
//	@Failure		403		{object}	GenericResp	"Forbidden - User does not have the required role"
//	@Failure		409		{object}	GenericResp	"Conflict - Hotel already exists"
//	@Failure		500		{object}	GenericResp	"Internal Server Error"
//	@Router			/hotels [post]
func (h *HotelHandler) HandleAddHotel(c *fiber.Ctx) error {
	log.Println("Sending request to add hotel")
	//get the current authenticated user
	user := c.Locals("user")
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	currentUser, ok := user.(*types.User)
	if !ok {
		log.Println("Unauthorized access attempt: user is not of type *types.User")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	// Check if the user is an admin
	if !userIsAdmin(currentUser) {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: Only users with role ADMIN can create hotels",
		})
	}
	var hotel types.Hotel
	if err := c.BodyParser(&hotel); err != nil {
		log.Printf("Error parsing hotel json: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": "Invalid input"})
	}
	//Get form data with images
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse form",
		})
	}

	//get image files from the form data
	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "no files provided to upload",
		})
	}
	exists, err := h.store.Hotel.HotelExists(c.Context(), &hotel)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to check if hotel exists",
		})
	}
	if exists {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "hotel already exists",
		})
	}
	//upload image and get urls
	imageUrls, err := h.ImageUploader.UploadImages(c.Context(), files)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to upload images",
		})
	}
	hotel.Images = imageUrls
	log.Println("Imageurl:", imageUrls)
	if err := c.BodyParser(&hotel); err != nil {
		log.Printf("Failed to parse hotel data: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	insertedHotel, err := h.store.Hotel.InsertHotel(c.Context(), currentUser.ID, &hotel)
	if err != nil {
		if err.Error() == "hotel already exists" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"error": "hotel already exist",
			})
		}
		log.Printf("Failed to insert hotel: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create hotel",
		})
	}

	return c.Status(http.StatusCreated).JSON(insertedHotel)
}

func userIsAdmin(user *types.User) bool {
	return user.Role == "ADMIN"
}
