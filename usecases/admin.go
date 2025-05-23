package usecases

import (
	"context"
	dto "tugas-akhir/dto/admin"
	"tugas-akhir/entities"
	"tugas-akhir/repositories"
	"tugas-akhir/utils/password"
	"tugas-akhir/utils/token"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AdminUseCase interface {
	Register(c echo.Context, req *dto.RegisterRequest) error
	Login(c echo.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type adminUsecase struct {
	adminRepo    repositories.AdminRepository
	passwordUtil password.PasswordUtil
	tokenUtil    token.TokenUtil
}

func NewAdminUsecase(adminRepo repositories.AdminRepository, passwordUtil password.PasswordUtil, tokenUtil token.TokenUtil) AdminUseCase {
	return &adminUsecase{
		adminRepo:    adminRepo,
		passwordUtil: passwordUtil,
		tokenUtil:    tokenUtil,
	}
}

func (au *adminUsecase) Register(c echo.Context, req *dto.RegisterRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	log := logrus.New()

	hashedPassword, err := au.passwordUtil.HashPassword(req.Password)
	if err != nil {
		return err
	}


	admin := &entities.Admin{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Role:     "admin",
		Password: hashedPassword,
	}

	err = au.adminRepo.CreateAdmin(ctx, admin)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (au *adminUsecase) Login(c echo.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	admin, err := au.adminRepo.GetAdmin(ctx, &entities.Admin{Email: req.Email})
	if err != nil {
		return nil, err
	}
	if err := au.passwordUtil.VerifyPassword(req.Password, admin.Password); err != nil {
		return nil, err
	}

	var token string

	token, err = au.tokenUtil.GenerateToken(admin.ID, admin.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Name:  admin.Name,
		Email: admin.Email,
		Token: token,
		Role:  admin.Role,
	}, nil
}
