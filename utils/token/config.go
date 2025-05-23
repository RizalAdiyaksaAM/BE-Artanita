package token

import (
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	msg "tugas-akhir/constant/messages"
	http_util "tugas-akhir/utils/http"
)

type JWTClaim struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"` // Menambahkan field role
	jwt.RegisteredClaims
}

func GetJWTConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTClaim)
		},
		ErrorHandler: jwtErrorHandler,
		SigningKey:   []byte(os.Getenv("JWT_KEY")),
	}
}

func jwtErrorHandler(c echo.Context, err error) error {
	code := http.StatusUnauthorized

	if errors.Is(err, echojwt.ErrJWTInvalid) {
		return http_util.HandleErrorResponse(
			c,
			code,
			msg.INVALID_TOKEN,
		)
	}

	return http_util.HandleErrorResponse(
		c,
		code,
		msg.UNAUTHORIZED,
	)
}