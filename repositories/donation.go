package repositories

import (
	"context"
	"errors"
	dto_base "tugas-akhir/dto/base"
	"tugas-akhir/entities"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DonationRepository interface {
	CreateDonation(ctx context.Context, donation *entities.Donation) error
	Update(ctx context.Context, donation *entities.Donation) error
	FindById(ctx context.Context, id uuid.UUID) (*entities.Donation, error)
	GetDonations(ctx context.Context, programDonationID uuid.UUID, searchName string, req *dto_base.PaginationRequest) ([]entities.Donation, int64, error)
	GetDonation(ctx context.Context) (*[]entities.Donation, error)
	GetDonationsLanding(ctx context.Context) (*[]entities.Donation, error)
	GetDonationByID(ctx context.Context, donationID uuid.UUID) (*entities.Donation, error)
	GetNotifikasi(ctx context.Context) (*[]entities.TransactionNotification, error)
	GetNotifikasiByDonationID(ctx context.Context, donationID uuid.UUID) (*entities.TransactionNotification, error)
	GetDonationByProgramID(ctx context.Context, programID uuid.UUID) (*[]entities.Donation, error)
}

type donationRepo struct {
	DB *gorm.DB
}

// NewDonationRepository creates a new instance of DonationRepository
func NewDonationRepository(db *gorm.DB) DonationRepository {
	return &donationRepo{
		DB: db,
	}
}

// FindById retrieves a donation by its ID

// CreateDonation creates a new donation record
func (dr *donationRepo) CreateDonation(ctx context.Context, donation *entities.Donation) error {
	log := logrus.New()

	log.Infof("Inserting donation into DB: %+v", donation)
	// Use GORM's WithContext to set the context for the transaction
	if err := dr.DB.WithContext(ctx).Create(donation).Error; err != nil {
		log.WithError(err).Error("Failed to insert donation into DB")
		return err
	}
	return nil
}

// FindById retrieves a donation by its ID
func (dr *donationRepo) FindById(ctx context.Context, id uuid.UUID) (*entities.Donation, error) {
	var donation entities.Donation
	if err := dr.DB.WithContext(ctx).First(&donation, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("donation not found")
		}
		return nil, err
	}
	return &donation, nil
}

// Update updates the donation record in the database (e.g., updating the donation status)
func (dr *donationRepo) Update(ctx context.Context, donation *entities.Donation) error {
	// Update the donation record based on the donation ID
	if err := dr.DB.WithContext(ctx).Save(donation).Error; err != nil {
		return err
	}
	return nil
}

func (dr *donationRepo) GetDonations(ctx context.Context, programDonationID uuid.UUID, searchName string, req *dto_base.PaginationRequest) ([]entities.Donation, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var donations []entities.Donation
	var totalData int64

	offset := (req.Page - 1) * req.Limit

	query := dr.DB.WithContext(ctx).Model(&entities.Donation{})

	// Jika ada filter ProgramDonationID, tambahkan kondisi filter
	if programDonationID != uuid.Nil {
		query = query.Where("program_donation_id = ?", programDonationID)
	}

	// Jika ada pencarian berdasarkan nama, tambahkan kondisi LIKE
	if searchName != "" {
		query = query.Where("name LIKE ?", "%"+searchName+"%")
	}

	// Menghitung total data
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// Query untuk mengambil data donasi sesuai dengan pagination dan filter
	query = query.Order(req.SortBy).Limit(req.Limit).Offset(offset)

	if err := query.Find(&donations).Error; err != nil {
		return nil, 0, err
	}

	return donations, totalData, nil
}

func (dr *donationRepo) GetDonationByID(ctx context.Context, donationID uuid.UUID) (*entities.Donation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var donation entities.Donation

	// Query untuk mendapatkan donasi berdasarkan ID
	if err := dr.DB.WithContext(ctx).Where("id = ?", donationID).First(&donation).Error; err != nil {
		return nil, err
	}

	return &donation, nil
}

func (dr *donationRepo) GetDonationsLanding(ctx context.Context) (*[]entities.Donation, error) {
	// Memeriksa konteks
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var donations []entities.Donation

	// Query untuk mendapatkan donasi dengan status 1
	if err := dr.DB.WithContext(ctx).Where("status = ?", 1).Find(&donations).Error; err != nil {
		return nil, err
	}

	return &donations, nil
}

func (dr *donationRepo) GetNotifikasi(ctx context.Context) (*[]entities.TransactionNotification, error) {
	// Memeriksa konteks
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var notifications []entities.TransactionNotification

	// Query untuk mendapatkan donasi dengan status 1
	if err := dr.DB.WithContext(ctx).Find(&notifications).Error; err != nil {
		return nil, err
	}

	return &notifications, nil
}

// GetDonationByProgramID mengambil donasi berdasarkan ProgramDonationID
func (dr *donationRepo) GetDonationByProgramID(ctx context.Context, programID uuid.UUID) (*[]entities.Donation, error) {
	var donations []entities.Donation
	if err := dr.DB.WithContext(ctx).Where("program_donation_id = ?", programID).Find(&donations).Error; err != nil {
		return nil, err
	}
	return &donations, nil
}

func (dr *donationRepo) GetDonation(ctx context.Context) (*[]entities.Donation, error) {
	var donations []entities.Donation
	if err := dr.DB.WithContext(ctx).Find(&donations).Error; err != nil {
		return nil, err
	}
	return &donations, nil
}

func (dr *donationRepo) GetNotifikasiByDonationID(ctx context.Context, donationID uuid.UUID) (*entities.TransactionNotification, error) {
	var notification entities.TransactionNotification
	if err := dr.DB.WithContext(ctx).
		Where("order_id = ? AND transaction_status = ?", donationID, "settlement").
		First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // atau custom handling
		}
		return nil, err
	}

	return &notification, nil
}
