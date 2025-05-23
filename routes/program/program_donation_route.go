package program

import (
	"tugas-akhir/config"
	"tugas-akhir/controllers"
	"tugas-akhir/drivers/cloudinary"
	"tugas-akhir/middlewares"
	"tugas-akhir/repositories"
	"tugas-akhir/usecases"
	"tugas-akhir/utils/validation"
	"tugas-akhir/utils/token"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitProgramDonationRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	cloudinaryInstance, _ := config.SetupCloudinary()
	cloudinaryService := cloudinary.NewCloudinaryService(cloudinaryInstance)

	programDonationRepo := repositories.NewProgramDonationRepository(db)
	programDonationUsecase := usecases.NewProgramDonationUsecase(programDonationRepo)
	programDonationController := controllers.NewProgramDonationController(programDonationUsecase, v, cloudinaryService)

	g.GET("/program-donations", programDonationController.GetProgramDonationAll)
	g.GET("/program-donations/:id", programDonationController.GetProgramDonationById)
	g.GET("/dashboard-donations", programDonationController.GetDashboard)
	
	g.POST("/program-donations", programDonationController.CreateProgramDonation)
	g.PUT("/program-donations/:id", programDonationController.UpdateProgramDonation)
	g.DELETE("/program-donations/:id", programDonationController.DeleteProgramDonation)
	
	g.Use(echojwt.WithConfig(token.GetJWTConfig()), middlewares.IsAdmin)
}
