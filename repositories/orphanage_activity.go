package repositories

import (
	"context"
	"errors"
	"log"
	dto_base "tugas-akhir/dto/base"
	"tugas-akhir/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrphanageActivityRepository interface {
	CreateActivity(ctx context.Context, activity *entities.OrphanageActivity) error
	GetActivityAll(ctx context.Context, req *dto_base.PaginationRequest, searchTitle string) ([]entities.OrphanageActivity, int64, error)
	GetActivityByID(ctx context.Context, id uuid.UUID) (*entities.OrphanageActivity, error)
	UpdateActivity(ctx context.Context, id uuid.UUID, activity *entities.OrphanageActivity) error
	DeleteActivity(ctx context.Context, id uuid.UUID) error
}

type orphanageActivityRepo struct {
	DB *gorm.DB
}

func NewOrphanageActivityRepository(db *gorm.DB) OrphanageActivityRepository {
	return &orphanageActivityRepo{
		DB: db,
	}
}

func (oar *orphanageActivityRepo) CreateActivity(ctx context.Context, activity *entities.OrphanageActivity) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	// Save the orphanage activity along with its images and videos in the database
	if err := oar.DB.Create(activity).Error; err != nil {
		log.Println("Error creating orphanage activity:", err)
		return err
	}

	return nil
}

func (oar *orphanageActivityRepo) GetActivityAll(
	ctx context.Context, req *dto_base.PaginationRequest, searchTitle string) ([]entities.OrphanageActivity, int64, error) {

	if oar == nil {
		return nil, 0, errors.New("user repo is nil")
	}

	if req == nil {
		return nil, 0, errors.New("pagination request is nil")
	}

	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var activities []entities.OrphanageActivity
	var totalData int64

	offset := (req.Page - 1) * req.Limit

	if oar.DB == nil {
		return nil, 0, errors.New("db is nil")
	}

	// Membuat query untuk menghitung total data dengan kondisi pencarian berdasarkan title
	query := oar.DB.WithContext(ctx).Model(&entities.OrphanageActivity{})

	if searchTitle != "" {
		query = query.Where("title LIKE ?", "%"+searchTitle+"%")
	}

	// Menghitung total data
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// Query untuk mengambil data program donasi sesuai dengan filter dan pagination
	query = query.Preload("ActivityImages").Preload("ActivityVideos").Order(req.SortBy).Limit(req.Limit).Offset(offset)

	if err := query.Find(&activities).Error; err != nil {
		return nil, 0, err
	}

	return activities, totalData, nil
}

func (oar *orphanageActivityRepo) GetActivityByID(ctx context.Context, id uuid.UUID) (*entities.OrphanageActivity, error) {
	if oar == nil {
		return nil, errors.New("user repo is nil")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var activity entities.OrphanageActivity

	if oar.DB == nil {
		return nil, errors.New("db is nil")
	}

	if err := oar.DB.WithContext(ctx).Preload("ActivityImages").Preload("ActivityVideos").Model(&entities.OrphanageActivity{}).Where("id = ?", id).First(&activity).Error; err != nil {
		return nil, err
	}

	return &activity, nil
}

func (oar *orphanageActivityRepo) UpdateActivity(ctx context.Context, id uuid.UUID, activity *entities.OrphanageActivity) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	tx := oar.DB.WithContext(ctx).Begin()

	if err := tx.Model(&entities.OrphanageActivity{}).Where("id = ?", id).Updates(activity).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (oar *orphanageActivityRepo) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := oar.DB.WithContext(ctx).Model(&entities.OrphanageActivity{}).Where("id = ?", id).Delete(&entities.OrphanageActivity{}).Error; err != nil {
		return err
	}
	return nil
}
