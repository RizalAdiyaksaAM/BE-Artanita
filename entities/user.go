package entities

import (
	"github.com/google/uuid"
)

type User struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name    string    `gorm:"type:varchar(50); not null"`
	Email   string    `gorm:"type:varchar(50); not null"`
	Address string    `gorm:"type:varchar(255); not null"`
	NoWA    int       `gorm:"type:int"`
}
