package database

import (
	"fmt"
	"log"
	"os"

	msg "tugas-akhir/constant/messages"
	"tugas-akhir/entities"
	log_util "tugas-akhir/utils/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DB_HOST      string
	DB_USERNAME  string
	DB_PASSWORD  string
	DB_NAME      string
	DB_PORT      string
	DB_SSL       string
	DB_TZ        string
	DB_LOG_LEVEL string
}

func ConnectDB(config Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.DB_HOST,
		config.DB_USERNAME,
		config.DB_PASSWORD,
		config.DB_NAME,
		config.DB_PORT,
		config.DB_SSL,
		config.DB_TZ,
	)

	logLevel := log_util.GetDBLogLevel(config.DB_LOG_LEVEL)
	logger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:        0,
		LogLevel:             logLevel,
		Colorful:             true,
		ParameterizedQueries: true,
	})
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		log.Fatal(msg.FAILED_CONNECT_DB, err)
	}
	migrate(db)
	return db
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entities.User{},
		&entities.OrphanageActivity{},
		&entities.ActivityImage{},
		&entities.ActivityVideo{},
		&entities.OrphanageUser{},
		&entities.Admin{},
		&entities.ProgramDonation{},
		&entities.ProgramDonationImage{},
		&entities.Donation{},
		&entities.TransactionNotification{},
	)
	if err != nil {
		log.Fatal(msg.FAILED_MIGRATE_DB, err)
	}
}
