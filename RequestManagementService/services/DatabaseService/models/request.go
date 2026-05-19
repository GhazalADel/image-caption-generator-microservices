package models

import "RequestManagementService/services/DatabaseService/consts"

type Request struct {
	ID           uint          `gorm:"primaryKey"`
	Email        string        `gorm:"not null"`
	Status       consts.Status `gorm:"type:varchar(20)"`
	ImageCaption string        `gorm:"type:text"`
	NewImageURL  string        `gorm:"type:text"`
}

func (Request) TableName() string {
	return "requests"
}
