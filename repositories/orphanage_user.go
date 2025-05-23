package repositories

import (
	"context"
	"errors"
	dto_base "tugas-akhir/dto/base"
	"tugas-akhir/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrphanageUserRepository interface {
	GetUserAll(
		ctx context.Context, req *dto_base.PaginationRequest,
		searchName, filterAddress, filterEducation, filterPosition, filterAge string,
	) ([]entities.OrphanageUser, int64, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.OrphanageUser, error)
	GetUserByPosition(ctx context.Context, position string) (*entities.OrphanageUser, error)
	CreateUser(ctx context.Context, user *entities.OrphanageUser) error
	UpdateUser(ctx context.Context, id uuid.UUID, user *entities.OrphanageUser) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type orphanageUserRepo struct {
	DB *gorm.DB
}

func NewOrphanageUserRepository(db *gorm.DB) OrphanageUserRepository {
	return &orphanageUserRepo{
		DB: db,
	}
}

func (our *orphanageUserRepo) CreateUser(ctx context.Context, user *entities.OrphanageUser) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := our.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (our *orphanageUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.OrphanageUser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var user entities.OrphanageUser

	if our.DB == nil {
		return nil, errors.New("db is nil")
	}

	if err := our.DB.WithContext(ctx).Model(&entities.OrphanageUser{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (our *orphanageUserRepo) GetUserByPosition(ctx context.Context, position string) (*entities.OrphanageUser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var user entities.OrphanageUser

	if our.DB == nil {
		return nil, errors.New("db is nil")
	}

	if err := our.DB.WithContext(ctx).Model(&entities.OrphanageUser{}).Where("position = ?", position).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (our *orphanageUserRepo) GetUserAll(
	ctx context.Context, req *dto_base.PaginationRequest,
	searchName, filterAddress, filterEducation, filterPosition, filterAge string,
) ([]entities.OrphanageUser, int64, error) {
	if our == nil {
		return nil, 0, errors.New("user repo is nil")
	}

	if req == nil {
		return nil, 0, errors.New("pagination request is nil")
	}

	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var users []entities.OrphanageUser
	var totalData int64

	offset := (req.Page - 1) * req.Limit

	if our.DB == nil {
		return nil, 0, errors.New("db is nil")
	}

	// Membuat query untuk menghitung total data dengan kondisi pencarian dan filter
	query := our.DB.WithContext(ctx).Model(&entities.OrphanageUser{})

	if searchName != "" {
		query = query.Where("name LIKE ?", "%"+searchName+"%")
	}

	if filterAddress != "" {
		query = query.Where("address LIKE ?", "%"+filterAddress+"%")
	}

	if filterEducation != "" {
		query = query.Where("education LIKE ?", "%"+filterEducation+"%")
	}

	if filterPosition != "" {
		query = query.Where("position LIKE ?", "%"+filterPosition+"%")
	}

	if filterAge != "" {
		query = query.Where("age = ?", filterAge)
	}

	// Menghitung total data
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// Query untuk mengambil data pengguna sesuai dengan filter dan pagination
	query = query.Order(req.SortBy).Limit(req.Limit).Offset(offset)

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalData, nil
}

func (our *orphanageUserRepo) UpdateUser(ctx context.Context, id uuid.UUID, user *entities.OrphanageUser) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := our.DB.WithContext(ctx).Model(&entities.OrphanageUser{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return err
	}

	return nil
}

func (our *orphanageUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := our.DB.WithContext(ctx).Model(&entities.OrphanageUser{}).Where("id = ?", id).Delete(&entities.OrphanageUser{}).Error; err != nil {
		return err
	}

	return nil
}
