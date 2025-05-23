package donation

import (
	"github.com/google/uuid"
)

type DonationRequest struct {
    Name      string    `json:"name" form:"name" validate:"required"`
    Address   string    `json:"address" form:"address" validate:"required"`
    NoWA      int       `json:"no_wa" form:"no_wa" validate:"required"`
    Email     string    `json:"email" form:"email" validate:"required"`
    Amount    int       `json:"amount" form:"amount" validate:"required"`
    Message   string    `json:"message" form:"message" validate:"required"`
    ProgramID uuid.UUID `json:"program_id" form:"program_id" validate:"omitempty,uuid4"` 
}


type DonationResponse struct {
	ID           string `json:"id"`
	Number       int    `json:"number"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	NoWA         int    `json:"no_wa"`
	Email        string `json:"email"`
	Amount       int    `json:"amount"`
	Message      string `json:"message"`
	Status       int    `json:"status"`
	SnapURL      string `json:"snap_url"`
	ProgramID    string `json:"program_id"`
	ProgramTitle string `json:"program_title"`
}

type TopUpReq struct {
	Amount int `json:"amount"`
}

type TopUpResp struct {
	SnapURL string `json:"snap_url"`
}

type ProgramDonationRequest struct {
	Title                 string                        `json:"title" form:"title"`
	Deskripsi             string                        `json:"deskripsi" form:"deskripsi" `
	GoalAmount            int                           `json:"goal_amount" form:"goal_amount" `
	ProgramDonationImages []ProgramDonationImageRequest `json:"program_donation_images" form:"program_donation_images"`
}

type ProgramDonationImageRequest struct {
	ImageUrl *string `json:"image_url" form:"image_url"`
}

type ProgramDonationResponse struct {
	ID                    string                         `json:"id"`
	Number                int                            `json:"number"`
	Title                 string                         `json:"title"`
	Deskripsi             string                         `json:"deskripsi"`
	GoalAmount            int                            `json:"goal_amount"`
	CurrentAmount         int                            `json:"current_amount"`
	ProgramDonationImages []ProgramDonationImageResponse `json:"program_donation_images"`
}

type ProgramDonationImageResponse struct {
	ImageUrl *string `json:"image_url"`
}

type DashboardDonationResponse struct {
	ProgramCount        int64 `json:"program_count"`
	TotalDonation       int   `json:"total_donation"`
	UniqueDonatorsCount int64 `json:"unique_donators_count"`
}

type DonationLandingResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Amount       int    `json:"amount"`
	Message      string `json:"message"`
}

type DonationChartResponse struct {
	ID string `json:"id"`
	Amount int `json:"amount"`
	Date string `json:"date"`
	ProgramDonation string `json:"program_donation"`
}

type DonaturNotifikasiResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Amount string `json:"amount"`
	ProgramDonation string `json:"program_donation"`
	Message string `json:"message"`
	Date string `json:"date"`
	Status string `json:"status"`
}