package api

import (
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

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(hotel)
}

// HandleAddHotel CreateHotel handles the creation of a new hotel
func (h *HotelHandler) HandleAddHotel(c *fiber.Ctx) error {
	log.Println("Sending request to add hotel")
	// Assuming you have a way to get the current authenticated user
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
	//upload image and get urls
	imageUrls, err := h.ImageUploader.UploadImages(c.Context(), files)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to upload images",
		})
	}
	hotel.Images = imageUrls

	if err := c.BodyParser(&hotel); err != nil {
		log.Printf("Failed to parse hotel data: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	insertedHotel, err := h.store.Hotel.InsertHotel(c.Context(), &hotel)
	if err != nil {
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
