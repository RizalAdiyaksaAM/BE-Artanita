package repositories

import (
	"context"
	"errors"
	"tugas-akhir/entities"

	dto_base "tugas-akhir/dto/base"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserAll(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.User, int64, error)
}

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		DB: db,
	}
}

func (ur *userRepo) CreateUser(ctx context.Context, user *entities.User) error {
    if err := ctx.Err(); err != nil {
        return err
    }
    if err := ur.DB.Create(user).Error; err != nil {
        return err
    }
    return nil
}

func (ur *userRepo) GetUserAll(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.User, int64, error) {
	if ur == nil {
		return nil, 0, errors.New("user repo is nil")
	}

	if req == nil {
		return nil, 0, errors.New("pagination request is nil")
	}

	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var users []entities.User
	var totalData int64

	offset := (req.Page - 1) * req.Limit

	if ur.DB == nil {
		return nil, 0, errors.New("db is nil")
	}

	// Menghitung total data
	if err := ur.DB.WithContext(ctx).Model(&entities.User{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	query := ur.DB.WithContext(ctx).Model(&entities.User{}).Order(req.SortBy).Limit(req.Limit).Offset(offset)
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalData, nil
}
