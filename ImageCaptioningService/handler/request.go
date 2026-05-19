package handler

import (
	"ImageCaptioningService/config"
	"ImageCaptioningService/services/DatabaseService/consts"
	"ImageCaptioningService/services/DatabaseService/database"
	"ImageCaptioningService/services/DatabaseService/datastore/request"
	"ImageCaptioningService/services/ObjectStorageService/objectstoremanager"

	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func HandleRequest(d amqp.Delivery) {
	requestIDStr := string(d.Body)

	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		log.Printf("Error converting request ID to an integer: %v", err)
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err.Error())
	}

	reqDataStore := request.New(db)
	reqMetadata, err := reqDataStore.GetRequestMetadataByRequestId(requestID)
	if err != nil {
		log.Printf("Failed to get request metadata by request ID : %v", err.Error())
	}
	req, err := reqDataStore.GetRequestById(requestID)
	if err != nil {
		log.Printf("Failed to get request data for request ID: %v", err.Error())
	}
	log.Printf("Successfully retrieved request related to ID: %v", requestID)

	fileName := req.Email + "_" + strconv.Itoa(int(req.ID)) + reqMetadata.Extension

	imageBytes, err := objectstoremanager.GetPictureFromDataStorage("pictures/", fileName)
	if err != nil {
		log.Fatalf("Failed to retrieve the image from object storage: %v", err)
	}

	log.Printf("Successfully retrieved picture related to request ID: %v from object storage", requestID)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load the configuration: %v", err)
	}

	APIKey := cfg.CaptionGeneratorAPI.KEY
	caption, err := generateCaption(imageBytes, APIKey)
	if err != nil {
		log.Fatalf("Failed to generate the caption: %v", err)
	}

	log.Printf("Successfully generated image caption for request ID: %v", requestID)

	req.Status = consts.READY_STATUS
	req.ImageCaption = caption

	err = reqDataStore.UpdateRequest(&req)
	if err != nil {
		log.Printf("Failed to update the request status and caption: %v", err)
		return
	}

	log.Printf("Successfully updated request %d with status 'ready' and the generated caption.", req.ID)
}

func generateCaption(imageBytes []byte, APIKey string) (string, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading configuration for Caption Generator API: %v", err)
	}
	apiURL := cfg.CaptionGeneratorAPI.URL
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(imageBytes))
	if err != nil {
		log.Printf("Error creating request for Caption Generator API: %v", err)
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+APIKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to Caption Generator API: %v", err)
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	log.Printf("Request sent successfully to Caption Generator API. Awaiting response...")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from Caption Generator API: %v", err)
		return "", fmt.Errorf("failed to read response: %v", err)
	}
	log.Printf("Response successfully received from Caption Generator API.")

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code received from Caption Generator API: %d, Response: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully processed the response from Caption Generator API. Parsing the caption...")
	return parseCaptionFromResponse(body)
}

func parseCaptionFromResponse(body []byte) (string, error) {
	var result []map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(result) > 0 {
		if caption, ok := result[0]["generated_text"].(string); ok {
			return caption, nil
		}
	}

	return "", fmt.Errorf("no caption generated or invalid format")
}
