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

func InitOrphanageUserRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	cloudinaryInstance, _ := config.SetupCloudinary()
	cloudinaryService := cloudinary.NewCloudinaryService(cloudinaryInstance)

	orphanageUserRepo := repositories.NewOrphanageUserRepository(db)
	orphanageUserUsecase := usecases.NewOrphanageUserUsecase(orphanageUserRepo, cloudinaryService)
	orphanageUserController := controllers.NewOrphanageUserController(orphanageUserUsecase, v)

	g.GET("/users", orphanageUserController.GetOrphanageUserAll)
	g.GET("/users/:id", orphanageUserController.GetOrphanageUserByID)
	g.GET("/users/position/:position", orphanageUserController.GetOrphanageUserByPosition)

	g.POST("/users", orphanageUserController.CreateUser)
	g.PUT("/users/:id", orphanageUserController.UpdateOrphanageUser)
	g.DELETE("/users/:id", orphanageUserController.DeleteOrphanageUser)
	
	g.Use(echojwt.WithConfig(token.GetJWTConfig()), middlewares.IsAdmin)


}
