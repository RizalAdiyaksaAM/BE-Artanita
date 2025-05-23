package usecases

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	dto_base "tugas-akhir/dto/base"
	dto "tugas-akhir/dto/user"
	"tugas-akhir/entities"
	"tugas-akhir/repositories"
	err_util "tugas-akhir/utils/error"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserUsecase interface {
	CreateUser(c echo.Context, req *dto.UserRequest) error
	GetUserAll(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.UserResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
}

type userUsecase struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uu *userUsecase) CreateUser(c echo.Context, req *dto.UserRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	user := &entities.User{
		ID:      uuid.New(),
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
		NoWA:    req.NoWA,
	}

	err := uu.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uu *userUsecase) GetUserAll(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.UserResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	baseURL := fmt.Sprintf(
		"%s?limit=%d&page=",
		c.Request().URL.Path,
		req.Limit,
	)

	var (
		next = baseURL + strconv.Itoa(req.Page+1)
		prev = baseURL + strconv.Itoa(req.Page-1)
	)

	if uu.userRepo == nil {
		return nil, nil, nil, errors.New("user repo is nil")
	}

	users, totalData, err := uu.userRepo.GetUserAll(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	user := []dto.UserResponse{}
	for _, u := range users {
		user = append(user, dto.UserResponse{
			ID:      u.ID.String(),
			Name:    u.Name,
			Email:   u.Email,
			Address: u.Address,
			NoWA:    u.NoWA,
		})
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(req.Limit)))
	paginationMetadata := &dto_base.PaginationMetadata{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: req.Page,
	}

	if req.Page > totalPage {
		return nil, nil, nil, err_util.ErrPageNotFound
	}

	if req.Page == 1 {
		prev = ""
	}

	if req.Page == totalPage {
		next = ""
	}

	link := &dto_base.Link{
		Next: next,
		Prev: prev,
	}

	return &user, paginationMetadata, link, nil
}
