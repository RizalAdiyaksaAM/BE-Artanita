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

type ProgramDonationRepository interface {
	CreateProgramDonation(ctx context.Context, programDonation *entities.ProgramDonation) error
	GetProgramDonationByID(ctx context.Context, programDonationID uuid.UUID) (*entities.ProgramDonation, error)
	GetProgramDonation(ctx context.Context, programDonation *entities.ProgramDonation) (*entities.ProgramDonation, error)
	GetProgramDonationAll(ctx context.Context, req *dto_base.PaginationRequest, searchTitle string) ([]entities.ProgramDonation, int64, error)
	UpdateProgramDonation(ctx context.Context, programDonation *entities.ProgramDonation) error
	DeleteProgramDonation(ctx context.Context, programDonationID uuid.UUID) error
	GetDashboardData(ctx context.Context) (int64, int, int64, error)
	UpdateCurrentAmount(ctx context.Context, programDonationID uuid.UUID, amount int) error
	GetFirstProgramDonation(ctx context.Context) (*entities.ProgramDonation, error)
}

type programDonationRepo struct {
	DB *gorm.DB
}

func NewProgramDonationRepository(db *gorm.DB) ProgramDonationRepository {
	return &programDonationRepo{
		DB: db,
	}
}

func (pdr *programDonationRepo) CreateProgramDonation(ctx context.Context, programDonation *entities.ProgramDonation) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	// Log program donation before saving
	logger := logrus.New()
	logger.Infof("Saving program donation with ID: %s", programDonation.ID.String())

	if err := pdr.DB.Create(programDonation).Error; err != nil {
		logger.Error("Failed to save program donation: ", err)
		return err
	}

	logger.Info("Program donation saved successfully")
	return nil
}

func (pdr *programDonationRepo) GetProgramDonationByID(ctx context.Context, programDonationID uuid.UUID) (*entities.ProgramDonation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var programDonation entities.ProgramDonation
	if err := pdr.DB.WithContext(ctx).Preload("DonationImage").Model(&entities.ProgramDonation{}).Where("id = ?", programDonationID).First(&programDonation).Error; err != nil {
		return nil, err
	}
	return &programDonation, nil
}

func (pdr *programDonationRepo) GetProgramDonationAll(ctx context.Context, req *dto_base.PaginationRequest, searchTitle string) ([]entities.ProgramDonation, int64, error) {
    if pdr == nil {
        return nil, 0, errors.New("program donation repo is nil")
    }

    if req == nil {
        return nil, 0, errors.New("pagination request is nil")
    }

    if err := ctx.Err(); err != nil {
        return nil, 0, err
    }

    var programDonations []entities.ProgramDonation
    var totalData int64

    offset := (req.Page - 1) * req.Limit

    if pdr.DB == nil {
        return nil, 0, errors.New("db is nil")
    }

    // Membuat query untuk menghitung total data dengan kondisi pencarian berdasarkan title
    query := pdr.DB.WithContext(ctx).Model(&entities.ProgramDonation{})
    
    if searchTitle != "" {
        query = query.Where("title LIKE ?", "%"+searchTitle+"%")
    }

    // Menghitung total data
    if err := query.Count(&totalData).Error; err != nil {
        return nil, 0, err
    }

    // Query untuk mengambil data program donasi sesuai dengan filter dan pagination
    query = query.Preload("DonationImage").Order(req.SortBy).Limit(req.Limit).Offset(offset)

    if err := query.Find(&programDonations).Error; err != nil {
        return nil, 0, err
    }

    return programDonations, totalData, nil
}


func (pdr *programDonationRepo) UpdateProgramDonation(ctx context.Context, programDonation *entities.ProgramDonation) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := pdr.DB.Save(programDonation).Error; err != nil {
		return err
	}
	return nil
}

func (pdr *programDonationRepo) DeleteProgramDonation(ctx context.Context, programDonationID uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := pdr.DB.WithContext(ctx).Model(&entities.ProgramDonation{}).Where("id = ?", programDonationID).Delete(&entities.ProgramDonation{}).Error; err != nil {
		return err
	}
	return nil
}

func (pdr *programDonationRepo) GetProgramDonation(ctx context.Context, programDonation *entities.ProgramDonation) (*entities.ProgramDonation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := pdr.DB.Where(programDonation).First(programDonation).Error; err != nil {
		return nil, err
	}

	return programDonation, nil
}

// Implementasi Repository
func (pdr *programDonationRepo) GetDashboardData(ctx context.Context) (int64, int, int64, error) {
	var programCount int64
	var totalDonation int
	var uniqueDonatorsCount int64

	// Menghitung jumlah program donasi
	if err := pdr.DB.Model(&entities.ProgramDonation{}).Count(&programCount).Error; err != nil {
		return 0, 0, 0, err
	}

	// Menghitung total donasi terkumpul dengan status = 1
	if err := pdr.DB.Model(&entities.Donation{}).Where("status = ?", 1).Select("SUM(amount)").Scan(&totalDonation).Error; err != nil {
		return 0, 0, 0, err
	}

	// Menghitung jumlah donatur unik (berdasarkan email atau nomor WA)
	if err := pdr.DB.Model(&entities.Donation{}).Distinct("email").Count(&uniqueDonatorsCount).Error; err != nil {
		return 0, 0, 0, err
	}

	return programCount, totalDonation, uniqueDonatorsCount, nil
}

func (pdr *programDonationRepo) UpdateCurrentAmount(ctx context.Context, programDonationID uuid.UUID, amount int) error {
    if pdr.DB == nil {
        return errors.New("db is nil")
    }

    // Memulai transaction untuk memastikan konsistensi data
    tx := pdr.DB.Begin()

    if tx.Error != nil {
        return tx.Error
    }

    // Ambil data ProgramDonation berdasarkan ID
    var programDonation entities.ProgramDonation
    if err := tx.Where("id = ?", programDonationID).First(&programDonation).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Update currentAmount
    programDonation.CurrentAmount += amount

    // Simpan perubahan ke database
    if err := tx.Save(&programDonation).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Commit transaction jika berhasil
    return tx.Commit().Error
}


func (pdr *programDonationRepo) GetFirstProgramDonation(ctx context.Context) (*entities.ProgramDonation, error) {
    if err := ctx.Err(); err != nil {
        return nil, err
    }

    var programDonation entities.ProgramDonation
    err := pdr.DB.WithContext(ctx).
        Preload("DonationImage").
        Model(&entities.ProgramDonation{}).
        Order("created_at ASC").
        Limit(1).
        First(&programDonation).Error
    if err != nil {
        return nil, err
    }

    return &programDonation, nil
}

