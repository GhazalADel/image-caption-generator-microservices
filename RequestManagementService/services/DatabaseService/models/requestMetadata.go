package models

type RequestMetadata struct {
	ID         uint   `gorm:"primaryKey"`
	RequestID  uint   `gorm:"not null"`
	Extension  string `gorm:"not null"`
	UploadedAt int64  `gorm:"not null"`
}

func (RequestMetadata) TableName() string {
	return "requests_metadata"
}
