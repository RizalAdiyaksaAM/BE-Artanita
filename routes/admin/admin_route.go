package admin

import (
	"tugas-akhir/controllers"
	"tugas-akhir/repositories"
	"tugas-akhir/usecases"
	"tugas-akhir/utils/password"
	"tugas-akhir/utils/token"
	"tugas-akhir/utils/validation"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitAdminRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {

	passwordUtil := password.NewPasswordUtil()
	tokenUtil := token.NewTokenUtil()

	adminRepo := repositories.NewAdminRepository(db)
	adminUseCase := usecases.NewAdminUsecase(adminRepo, passwordUtil,  tokenUtil)
	adminController := controllers.NewAdminController(adminUseCase, v, tokenUtil)

	// Public routes
	g.POST("/admin/login", adminController.Login)
	g.POST("/admin/register", adminController.Register)

}
