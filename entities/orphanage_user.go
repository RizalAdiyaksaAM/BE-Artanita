package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrphanageUser struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name      string    `gorm:"type:varchar(50); not null"`
	Address   *string    `gorm:"type:varchar(255); not null"`
	Age       *int       `gorm:"type:int"`
	Education *string    `gorm:"type:varchar(50); not null"`
	Position  *string    `gorm:"type:varchar(50); not null"`
	Image     *string    `gorm:"type:varchar(255); not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

