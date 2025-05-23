package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionNotification represents a record of notification received from Midtrans
type TransactionNotification struct {
	ID                uuid.UUID `gorm:"primaryKey;type:uuid"`
	OrderID           string    `gorm:"type:varchar(50);not null"`
	TransactionStatus string    `gorm:"type:varchar(50);not null"`
	GrossAmount       string    `gorm:"type:varchar(50);not null"`
	TransactionTime   string    `gorm:"type:varchar(50);not null"`
	SignatureKey      string    `gorm:"type:varchar(255);not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
