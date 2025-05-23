package main

import (

	"log"
	"os"
	"tugas-akhir/config"
	"tugas-akhir/drivers/database"
	"tugas-akhir/routes"
	"tugas-akhir/utils/validation"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// "golang.ngrok.com/ngrok"

	// ngrokconfig "golang.ngrok.com/ngrok/config"
	"gorm.io/gorm"
)

var db *gorm.DB
var v *validation.Validator

func init() {
	config.LoadEnv()
	config.InitConfigDB()
	db = database.ConnectDB(config.InitConfigDB())
	v = validation.NewValidator()
}

func main() {
	// Buat instance Echo
	e := echo.New()

	// Tambahkan middleware untuk logging request
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Konfigurasi CORS yang lebih permisif untuk development
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Izinakan semua origin untuk testing
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXCSRFToken,
			// Tambahkan header khusus Midtrans jika ada
		},
		AllowMethods: []string{
			echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))

	// Inisialisasi rute
	routes.InitRoute(e, db, v)

	// // Buat context untuk ngrok
	// ctx := context.Background()

	// // Siapkan listener ngrok
	// listener, err := ngrok.Listen(ctx,
	// 	ngrokconfig.HTTPEndpoint(),
	// 	ngrok.WithAuthtokenFromEnv(), // Memerlukan NGROK_AUTHTOKEN di environment variables
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Tampilkan URL ngrok
	// log.Println("Ngrok Public URL:", listener.URL())

	// // Jalankan server dengan listener ngrok
	// e.Logger.Fatal(e.Server.Serve(listener))

	//Tentukan port dari environment variable atau gunakan default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port jika tidak ditentukan
	}

	// Jalankan server di localhost dengan port yang ditentukan
	log.Println("Server berjalan di http://localhost:" + port)
	e.Logger.Fatal(e.Start(":" + port))
}
