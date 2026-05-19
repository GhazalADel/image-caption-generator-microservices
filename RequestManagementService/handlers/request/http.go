package request

import (
	"RequestManagementService/services/DatabaseService/consts"
	"RequestManagementService/services/DatabaseService/datastore"
	"RequestManagementService/services/DatabaseService/models"
	"RequestManagementService/services/ObjectStorageService/objectstoremanager"
	"RequestManagementService/services/QueueService/rabbitmq"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/mail"
	"path/filepath"
	"strconv"
	"time"
)

type RequestHandler struct {
	datastore datastore.Request
}

func New(requests datastore.Request) *RequestHandler {
	return &RequestHandler{datastore: requests}
}

func (r RequestHandler) AddRequest(c echo.Context) error {
	err := c.Request().ParseMultipartForm(15 << 20)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "could not parse form")
	}

	email := c.FormValue("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, "Email is required")
	}
	_, emailErr := mail.ParseAddress(email)
	if emailErr != nil {
		return c.JSON(http.StatusBadRequest, "Email is invalid")
	}

	picture, handler, err := c.Request().FormFile("picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Picture is required")
	}
	defer func() {
		if pictureErr := picture.Close(); pictureErr != nil {
			c.Logger().Warnf("Failed to close picture file: %v", pictureErr)
		}
	}()

	var request models.Request

	request.Email = email
	request.Status = consts.Status(string(consts.PENDING_STATUS))

	createdReq, err := r.datastore.CreateRequest(&request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Could not create request: %v", err.Error()))
	}

	fileExtension := filepath.Ext(handler.Filename)
	fileName := email + "_" + strconv.Itoa(int(createdReq.ID)) + fileExtension

	var requestMetaData models.RequestMetadata
	requestMetaData.RequestID = createdReq.ID
	requestMetaData.Extension = fileExtension
	requestMetaData.UploadedAt = time.Now().Unix()

	_, err = r.datastore.CreateRequestMetadata(&requestMetaData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Could not create request meta data: %v", err.Error()))
	}

	uploadErr := objectstoremanager.AddPictureToDataStorage("pictures/", fileName, picture)
	if uploadErr != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Failed to store picture in data storage: %v", uploadErr.Error()))
	}

	rmq, err := rabbitmq.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Could not connect to rabbitmq: %v", err.Error()))
	}
	defer rmq.Close()

	err = rmq.AddToQueue("requests", strconv.Itoa(int(createdReq.ID)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to add message to queue: %v", err.Error()))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Request Submitted Successfully...Request ID : %v", createdReq.ID))
}

func (r RequestHandler) GetRequest(c echo.Context) error {
	id := c.Param("id")
	requestId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid parameter id")
	}

	req, err := r.datastore.GetRequestById(requestId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "could not retrieve request")
	}

	if req.Status == consts.Status(string(consts.DONE_STATUS)) {
		return c.JSON(http.StatusOK, req.NewImageURL)

	}

	return c.JSON(http.StatusOK, "Your request hasn't been done yet")
}
