package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrphanageActivity struct {
	ID             uuid.UUID       `gorm:"primaryKey;type:uuid"`
	Title          string          `gorm:"type:varchar(50); not null"`
	Description    string          `gorm:"type:text; not null"`
	Location       string          `gorm:"type:varchar(255); not null"`
	Time           string          `gorm:"type:varchar(255); not null"`
	ActivityImages []ActivityImage `gorm:"foreignKey:ActivityID;references:ID"`
	ActivityVideos []ActivityVideo `gorm:"foreignKey:ActivityID;references:ID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type ActivityImage struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid"`
	ImageUrl   *string   `gorm:"type:varchar(255)"`
	ActivityID uuid.UUID `gorm:"type:uuid"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type ActivityVideo struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid"`
	VideoUrl   *string   `gorm:"type:varchar(255)"`
	ActivityID uuid.UUID `gorm:"type:uuid"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

