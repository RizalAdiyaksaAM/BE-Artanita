package controllers

import (
	"net/http"
	"strconv"
	"strings"
	dto_base "tugas-akhir/dto/base"
	dto "tugas-akhir/dto/user"
	"tugas-akhir/usecases"
	"tugas-akhir/utils/validation"

	http_util "tugas-akhir/utils/http"
	msg "tugas-akhir/constant/messages"

	"github.com/labstack/echo/v4"
)

type userController struct {
	userUsecase usecases.UserUsecase
	validator   *validation.Validator
}

func NewUserController(userUsecase usecases.UserUsecase, validator *validation.Validator) *userController {
	return &userController{
		userUsecase: userUsecase,
		validator:   validator,
	}
}

func (uc *userController) CreateUser(ctx echo.Context) error {
	request := new(dto.UserRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := uc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}


	if err := uc.userUsecase.CreateUser(ctx, request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_CREATE_USER)
	}


	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_CREATE_USER, nil)
}

func (uc *userController) GetUserAll(c echo.Context) error {
	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := uc.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy, 
	}

	if err := uc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, metadata, link, err := uc.userUsecase.GetUserAll(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_USER_ALL)
	}

	return http_util.HandlePaginationResponse(c, msg.SUCCESS_GET_USER_ALL, result, metadata, link)

	
}

func (pc *userController) convertQueryParams(page, limit string) (int, int, error) {
	if page == "" {
		page = "1"
	}

	if limit == "" {
		limit = "10"
	}

	var (
		intPage, intLimit int
		err               error
	)

	intPage, err = strconv.Atoi(page)
	if err != nil {
		return 0, 0, err
	}

	intLimit, err = strconv.Atoi(limit)
	if err != nil {
		return 0, 0, err
	}

	return intPage, intLimit, nil
}
