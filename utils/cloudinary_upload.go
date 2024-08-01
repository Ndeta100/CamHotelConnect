package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	uuid2 "github.com/google/uuid"
	"log"
	"mime/multipart"
)

type ImageUploader interface {
	UploadImages(ctx context.Context, filePaths []*multipart.FileHeader) ([]string, error)
}

type CloudinaryUploader struct {
	cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryUploader(cloudName, apiKey, apiSecret string) (*CloudinaryUploader, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &CloudinaryUploader{cloudinary: cld}, nil
}

func (upload *CloudinaryUploader) UploadImages(ctx context.Context, filePaths []*multipart.FileHeader) ([]string, error) {
	var urls []string
	if len(filePaths) == 0 {
		return nil, errors.New("no file in path")
	}
	//iterate through filepath
	for _, filePath := range filePaths {
		err := func(filePath *multipart.FileHeader) error {
			// open file
			file, err := filePath.Open()
			if err != nil {
				return fmt.Errorf("error opening file %s: %w", filePath.Filename, err)
			}
			//Ensure each file is immediately closed when upload is finish
			defer func() {
				if err := file.Close(); err != nil {
					log.Printf("error closing file %s: %w", filePath.Filename, err)
				}
			}()

			uniqueField := true
			publicID, err := generateUniquePublicID()
			if err != nil {
				fmt.Println("Error generating public ID")
			}
			overwrite := false
			uploadParams := uploader.UploadParams{
				Folder:         "",
				UniqueFilename: &uniqueField,
				PublicID:       publicID,
				Overwrite:      &overwrite,
			}
			uploadResults, err := upload.cloudinary.Upload.Upload(ctx, filePath, uploadParams)
			if err != nil {
				log.Println("Error uploading file", err)
			}
			urls = append(urls, uploadResults.URL)
			return nil
		}(filePath)
		if err != nil {
			return nil, fmt.Errorf("error uploading file %s: %w", filePath.Filename, err)
		}
	}
	return urls, nil
}

func generateUniquePublicID() (string, error) {
	uuid, err := uuid2.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}
	return uuid.String(), err
}
