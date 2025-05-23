package user

import (
	"tugas-akhir/controllers"
	"tugas-akhir/repositories"
	"tugas-akhir/usecases"
	"tugas-akhir/utils/validation"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitUserRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	userRepo := repositories.NewUserRepository(db)

	userUseCase := usecases.NewUserUseCase(userRepo)
	userController := controllers.NewUserController(userUseCase, v)

	g.POST("/donation", userController.CreateUser)
	g.GET("/donation", userController.GetUserAll)
}
