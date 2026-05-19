package handler

import (
	"ImageGenerationService/config"
	"ImageGenerationService/services/DatabaseService/consts"
	"ImageGenerationService/services/DatabaseService/database"
	"ImageGenerationService/services/DatabaseService/datastore/request"
	"ImageGenerationService/services/DatabaseService/models"
	"ImageGenerationService/services/MailSenderService/gmail"
	"ImageGenerationService/services/ObjectStorageService/objectstoremanager"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func ProcessReadyRequests() error {
	db, err := database.GetConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err.Error())
	}

	reqDataStore := request.New(db)
	requests, err := findReadyRequests(reqDataStore)
	if err != nil {
		return err
	}
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	apiURL := cfg.ImageGeneratorAPI.URL
	apiKey := cfg.ImageGeneratorAPI.KEY

	for _, reqElement := range requests {
		log.Printf("Processing request ID: %v", reqElement.ID)
		caption := reqElement.ImageCaption
		log.Printf("Image caption for request ID: %v is %v", reqElement.ID, caption)
		imageBytes, err := generateImage(caption, apiURL, apiKey)
		if err != nil {
			return fmt.Errorf("failed to generate image via external API: %v", err.Error())
		}
		log.Printf("Image successfully generated for request ID: %v", reqElement.ID)
		fileName, err := storeImageInObjectStorage(imageBytes, reqElement.ID, reqElement.Email)
		if err != nil {
			return err
		}
		log.Printf("Image successfully stored in object storage for request ID: %v", reqElement.ID)
		imageURL, err := objectstoremanager.GetPictureURL("generated-pictures/", fileName)
		if err != nil {
			return err
		}
		log.Printf("URL of the image successfully retrieved from object storage for request ID: %v", reqElement.ID)

		reqElement.NewImageURL = imageURL
		err = reqDataStore.UpdateRequest(&reqElement)
		if err != nil {
			return err
		}

		log.Printf("URL successfully updated in the database for request ID: %v", reqElement.ID)

		err = gmail.SendCustomEmail("Your Requested URL", imageURL, reqElement.Email)
		if err != nil {
			return err
		}

		log.Printf("Email successfully sent for request ID: %v", reqElement.ID)

		reqElement.Status = consts.DONE_STATUS
		err = reqDataStore.UpdateRequest(&reqElement)
		if err != nil {
			return err
		}

		log.Printf("Status successfully updated in the database for request ID: %v", reqElement.ID)
	}
	return nil
}

func findReadyRequests(reqDataStore request.RequestDatastore) ([]models.Request, error) {
	reqs, err := reqDataStore.GetReadyRequests()
	if err != nil {
		return nil, err
	}
	return reqs, nil
}

func generateImage(caption string, apiURL string, apiKey string) ([]byte, error) {
	reqBody := map[string]string{
		"inputs": caption,
	}
	payloadBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Failed to close response body")
		}
	}(resp.Body)

	imageBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned an error: %s", string(imageBytes))
	}

	return imageBytes, nil
}

func storeImageInObjectStorage(imageBytes []byte, id uint, email string) (string, error) {
	imageReader := bytes.NewReader(imageBytes)

	fileName := strconv.Itoa(int(id)) + "_" + email

	err := objectstoremanager.AddPictureToDataStorage("generated-pictures/", fileName, imageReader)
	if err != nil {
		return "", fmt.Errorf("failed to store image: %v", err)
	}

	return fileName, nil
}
