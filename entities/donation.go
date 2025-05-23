package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Donation struct {
	ID                uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name              string    `gorm:"type:varchar(50); not null"`
	Address           string    `gorm:"type:varchar(255); not null"`
	NoWA              int       `gorm:"type:int"`
	Email             string    `gorm:"type:varchar(50); not null"`
	Amount            int       `gorm:"type:int"`
	Message           string    `gorm:"type:text; not null"`
	Status            int       `gorm:"type:int"`
	SnapURL           string    `gorm:"type:varchar(255); not null"`
	ProgramDonationID uuid.UUID `gorm:"type:uuid;not null"` // Foreign Key
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

type ProgramDonation struct {
	ID            uuid.UUID              `gorm:"primaryKey;type:uuid"`
	Title         string                 `gorm:"type:varchar(50); not null"`
	Deskripsi     string                 `gorm:"type:text; not null"`
	GoalAmount    int                    `gorm:"type:int"`
	CurrentAmount int                    `gorm:"type:int"`
	DonationImage []ProgramDonationImage `gorm:"foreignKey:ProgramID;references:ID"`
	Donations     []Donation             `gorm:"foreignKey:ProgramDonationID;references:ID"` // Relasi One-to-Many
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type ProgramDonationImage struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	ImageUrl  *string   `gorm:"type:varchar(255)"`
	ProgramID uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
