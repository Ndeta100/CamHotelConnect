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
	"os"
	"sync"
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
	if len(filePaths) == 0 {
		return nil, errors.New("no file in path")
	}
	var (
		urls   = make([]string, len(filePaths))
		errorS = make([]error, len(filePaths))
		wg     sync.WaitGroup
		mu     sync.Mutex                //To protect shared data
		sem    = make(chan struct{}, 10) //Semaphore to limit concurrent upload
	)
	//iterate through filepath
	for i, filePath := range filePaths {
		log.Println("uploading file", i, filePath)
		wg.Add(1)
		go func(i int, filePath *multipart.FileHeader) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			// open file
			file, err := filePath.Open()
			if err != nil {
				fmt.Errorf("error opening file %s: %w", filePath.Filename, err)
				return
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
			cloudFolder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
			uploadParams := uploader.UploadParams{
				Folder:         cloudFolder,
				UniqueFilename: &uniqueField,
				PublicID:       publicID,
				Overwrite:      &overwrite,
			}
			uploadResults, err := upload.cloudinary.Upload.Upload(ctx, filePath, uploadParams)
			if err != nil {
				log.Println("Error uploading file", err)
				errorS[i] = fmt.Errorf("error uploading file %s: %w", filePath.Filename, err)
				return
			}
			// Safely append url to list
			mu.Lock()
			urls[i] = uploadResults.URL
			defer mu.Unlock()
			return
		}(i, filePath)
	}
	wg.Wait()
	//check for errors if any
	for _, err := range errorS {
		if err != nil {
			return nil, err
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
