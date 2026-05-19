package request

import (
	"ImageGenerationService/services/DatabaseService/consts"
	"ImageGenerationService/services/DatabaseService/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type RequestDatastore struct {
	db *gorm.DB
}

func New(db *gorm.DB) RequestDatastore {
	return RequestDatastore{db: db}
}

func (r RequestDatastore) CreateRequest(request *models.Request) (models.Request, error) {
	var tmpRequest *models.Request
	tmpRequest = request
	createdRequest := r.db.Create(&tmpRequest)
	if createdRequest.Error != nil {
		return models.Request{}, fmt.Errorf("request creation failed")
	}
	return *tmpRequest, nil
}

func (r RequestDatastore) CreateRequestMetadata(requestMetadata *models.RequestMetadata) (models.RequestMetadata, error) {
	var tmpRequestMetadata *models.RequestMetadata
	tmpRequestMetadata = requestMetadata
	createdRequestMetadata := r.db.Create(&tmpRequestMetadata)
	if createdRequestMetadata.Error != nil {
		return models.RequestMetadata{}, fmt.Errorf("request meta data creation failed")
	}
	return *tmpRequestMetadata, nil
}

func (r RequestDatastore) GetRequestById(id int) (models.Request, error) {
	var request models.Request

	err := r.db.Where("id = ?", id).First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Request{}, fmt.Errorf("request not found")
		}
		return models.Request{}, fmt.Errorf("failed to get request: %v", err)
	}

	return request, nil
}

func (r RequestDatastore) GetRequestMetadataByRequestId(id int) (models.RequestMetadata, error) {
	var requestMetadata models.RequestMetadata

	err := r.db.Where("request_id = ?", id).First(&requestMetadata).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.RequestMetadata{}, fmt.Errorf("request meta data not found")
		}
		return models.RequestMetadata{}, fmt.Errorf("failed to get request meta data: %v", err)
	}

	return requestMetadata, nil
}

func (r RequestDatastore) GetLastRequestId() (int, error) {
	var request models.Request

	err := r.db.Order("id desc").First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("no requests found")
		}
		return 0, fmt.Errorf("failed to get last request id: %v", err)
	}

	return int(request.ID), nil
}

func (r RequestDatastore) UpdateRequest(request *models.Request) error {
	if err := r.db.Save(request).Error; err != nil {
		return err
	}
	return nil
}

func (r RequestDatastore) GetReadyRequests() ([]models.Request, error) {
	var readyRequests []models.Request
	err := r.db.Where("status = ?", consts.READY_STATUS).Find(&readyRequests).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query ready requests: %v", err)
	}

	return readyRequests, nil
}
