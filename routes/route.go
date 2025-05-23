package routes

import (
	"tugas-akhir/utils/validation"

	"tugas-akhir/routes/admin"
	"tugas-akhir/routes/donation"
	"tugas-akhir/routes/orphanage"
	"tugas-akhir/routes/user"
	"tugas-akhir/routes/program"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoute(e *echo.Echo, db *gorm.DB, v *validation.Validator) {
	baseRoute := e.Group("/api/v1")

	userRoute := baseRoute.Group("")
	activityRoute := baseRoute.Group("")
	orphanageUserRoute := baseRoute.Group("")
	adminRoute := baseRoute.Group("")
	programDonationRoute := baseRoute.Group("")
	donationRoute := baseRoute.Group("")

	user.InitUserRoute(userRoute, db, v)
	orphanage.InitActivityRoute(activityRoute, db, v)
	orphanage.InitOrphanageUserRoute(orphanageUserRoute, db, v)
	admin.InitAdminRoute(adminRoute, db, v)
	program.InitProgramDonationRoute(programDonationRoute, db, v)
	donation.InitDonationRoute(donationRoute, db, v)
}
