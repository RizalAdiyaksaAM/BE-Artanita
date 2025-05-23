package orphanage

import (
	"tugas-akhir/config"
	"tugas-akhir/controllers"
	"tugas-akhir/drivers/cloudinary"
	"tugas-akhir/middlewares"
	"tugas-akhir/repositories"
	"tugas-akhir/usecases"
	"tugas-akhir/utils/token"
	"tugas-akhir/utils/validation"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitActivityRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	cloudinaryInstance, _ := config.SetupCloudinary()
	cloudinaryService := cloudinary.NewCloudinaryService(cloudinaryInstance)

	activityRepo := repositories.NewOrphanageActivityRepository(db)
	activityUseCase := usecases.NewOrphanageActivityUsecase(activityRepo)
	activityController := controllers.NewOrphanageActivityController(activityUseCase, v, cloudinaryService)

	g.GET("/activities/:id", activityController.GetActivityById)
	g.GET("/activities", activityController.GetActivityAll)
	g.POST("/activities", activityController.CreateActivity)
	g.PUT("/activities/:id", activityController.UpdateActivity)
	g.DELETE("/activities/:id", activityController.DeleteActivity)
	
	g.Use(echojwt.WithConfig(token.GetJWTConfig()), middlewares.IsAdmin)
	
}
