package datastore

import "RequestManagementService/services/DatabaseService/models"

type Request interface {
	CreateRequest(request *models.Request) (models.Request, error)
	GetRequestById(id int) (models.Request, error)
	GetLastRequestId() (int, error)
	CreateRequestMetadata(requestMetadata *models.RequestMetadata) (models.RequestMetadata, error)
	GetRequestMetadataByRequestId(id int) (models.RequestMetadata, error)
	UpdateRequest(request *models.Request) error
}
