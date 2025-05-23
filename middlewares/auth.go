package middlewares

import (
	"net/http"
	msg "tugas-akhir/constant/messages"
	http_util "tugas-akhir/utils/http"
	"tugas-akhir/utils/token"

	"github.com/labstack/echo/v4"
)


func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := token.NewTokenUtil().GetClaims(c)
		
		// Memeriksa apakah role pengguna adalah admin
		if claims.Role != "admin" {
			return http_util.HandleErrorResponse(
				c,
				http.StatusForbidden, 
				msg.FORBIDDEN_ACCESS,
			)
		}

		return next(c)
	}
}



// HasAnyRole memeriksa apakah pengguna memiliki salah satu dari role yang diberikan
func HasAnyRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := token.NewTokenUtil().GetClaims(c)
			
			// Memeriksa apakah role pengguna terdapat dalam daftar roles yang diizinkan
			for _, role := range roles {
				if claims.Role == role {
					return next(c)
				}
			}
			
			return http_util.HandleErrorResponse(
				c,
				http.StatusForbidden, 
				msg.FORBIDDEN_ACCESS,
			)
		}
	}
}