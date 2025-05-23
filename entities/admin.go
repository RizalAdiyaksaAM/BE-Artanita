package entities

import "github.com/google/uuid"

type Admin struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name     string    `gorm:"type:varchar(50); not null"`
	Role     string    `gorm:"type:varchar(50); not null"`
	Email    string    `gorm:"type:varchar(50); not null"`
	Password string    `gorm:"type:varchar(255); not null"`
}
