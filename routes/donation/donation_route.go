package donation

import (
	"tugas-akhir/config"
	"tugas-akhir/controllers"
	"tugas-akhir/repositories"
	"tugas-akhir/usecases"
	"tugas-akhir/utils/validation"
	midtrans "tugas-akhir/utils/midtrans"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitDonationRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {

	// Inisialisasi konfigurasi Midtrans
	config := config.InitConfigMidtrans()

	// Membuat client Midtrans berdasarkan konfigurasi
	midtransClient := midtrans.NewClient(config)

	// Inisialisasi repository-repository yang diperlukan
	donationRepo := repositories.NewDonationRepository(db)
	programDonationRepo := repositories.NewProgramDonationRepository(db)
	transactionNotificationRepo := repositories.NewTransactionNotificationRepository(db)
	
	// Inisialisasi Usecase untuk Donation
	donationUsecase := usecases.NewDonationUsecase(donationRepo,config,  programDonationRepo, transactionNotificationRepo)

	// Inisialisasi Controller untuk Donation
	donationController := controllers.NewDonationController(donationUsecase, v, midtransClient)

	// Daftarkan route POST untuk membuat donasi
	g.POST("/donations", donationController.CreateDonation)

	// Daftarkan route POST untuk menerima webhook dari Midtrans
	g.POST("/midtrans-webhook", donationController.MidtransWebhook)
	g.GET("/donations-all", donationController.GetDonations)
	g.GET("/donations/:id", donationController.GetDonationByID)
	g.GET("/donations-user", donationController.GetDonationsLanding)
	g.GET("/donations-chart", donationController.GetChartDonation)
	g.GET("/donations-notifikasi", donationController.GetDonaturNotifikasi)
	g.GET("/donations-program/:id", donationController.GetDonaturByProgramDonation)
}
