package repositories

import (
	"context"
	"tugas-akhir/entities"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionNotificationRepository interface {
	CreateNotification(ctx context.Context, notification *entities.TransactionNotification) error
}

type transactionNotificationRepo struct {
	DB *gorm.DB
}

func NewTransactionNotificationRepository(db *gorm.DB) TransactionNotificationRepository {
	return &transactionNotificationRepo{
		DB: db,
	}
}

func (t *transactionNotificationRepo) CreateNotification(ctx context.Context, notification *entities.TransactionNotification) error {
	log := logrus.New()

	// Generate UUID jika belum ada
	if notification.ID == uuid.Nil {
		notification.ID = uuid.New()
	}

	log.Infof("Inserting transaction notification into DB: %+v", notification)
	if err := t.DB.WithContext(ctx).Create(notification).Error; err != nil {
		log.WithError(err).Error("Failed to insert transaction notification into DB")
		return err
	}
	return nil
}
