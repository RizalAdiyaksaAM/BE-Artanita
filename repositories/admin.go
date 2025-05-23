package repositories

import (
	"context"
	"errors"
	"tugas-akhir/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(ctx context.Context, admin *entities.Admin) error
	GetAdminByID(ctx context.Context, id uuid.UUID) (*entities.Admin, error)
	GetAdmin(ctx context.Context, admin *entities.Admin) (*entities.Admin, error)
	UpdateAdmin(ctx context.Context, id uuid.UUID, admin *entities.Admin) error
	DeleteAdmin(ctx context.Context, id uuid.UUID) error
}

type adminRepo struct {
	DB *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepo{
		DB: db,
	}
}

func (ar *adminRepo) CreateAdmin(ctx context.Context, admin *entities.Admin) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return ar.DB.Create(admin).Error
}

func (ar *adminRepo) GetAdmin(ctx context.Context, admin *entities.Admin) (*entities.Admin, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := ar.DB.Where(admin).First(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

// func (ar *adminRepo) GetAdminAll(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Admin, int64, error) {
// 	if ar == nil {
// 		return nil, 0, errors.New("adminRepository is nil")
// 	}

// 	if req == nil {
// 		return nil, 0, errors.New("PaginationRequest is nil")
// 	}

// 	if err := ctx.Err(); err != nil {
// 		return nil, 0, err
// 	}

// 	var admins []entities.Admin
// 	var totalData int64

// 	offset := (req.Page - 1) * req.Limit

// 	// Pastikan ar.DB tidak nil sebelum menggunakan
// 	if ar.DB == nil {
// 		return nil, 0, errors.New("DB connection is nil")
// 	}

// 	// Menghitung total data
// 	if err := ar.DB.WithContext(ctx).Model(&entities.Admin{}).Count(&totalData).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	query := ar.DB.WithContext(ctx).Model(&entities.Admin{}).Order(req.SortBy).Limit(req.Limit).Offset(offset)
// 	if err := query.Find(&admins).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	return admins, totalData, nil
// }

func (ar *adminRepo) GetAdminByID(ctx context.Context, id uuid.UUID) (*entities.Admin, error) {
	if ar == nil {
		return nil, errors.New("adminRepository is nil")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var admin entities.Admin	
	if err := ar.DB.WithContext(ctx).Model(&entities.Admin{}).Where("id = ?", id).First(&admin).Error; err != nil {
		return nil, err
	}

	return &admin, nil
}

func (ar *adminRepo) UpdateAdmin(ctx context.Context, id uuid.UUID, admin *entities.Admin) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return ar.DB.Model(&entities.Admin{}).Where("id = ?", id).Updates(admin).Error
}

func (ar *adminRepo) DeleteAdmin(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return ar.DB.Model(&entities.Admin{}).Where("id = ?", id).Delete(&entities.Admin{}).Error
}