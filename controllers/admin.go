package controllers

import (
	"fmt"
	"net/http"
	msg "tugas-akhir/constant/messages"
	dto "tugas-akhir/dto/admin"
	"tugas-akhir/usecases"
	http_util "tugas-akhir/utils/http"
	"tugas-akhir/utils/token"
	"tugas-akhir/utils/validation"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type adminController struct {
	adminUsecase usecases.AdminUseCase
	validator    *validation.Validator
	tokenUtil    token.TokenUtil
}

func NewAdminController(adminUsecase usecases.AdminUseCase, validator *validation.Validator, tokenUtil token.TokenUtil) *adminController {
	return &adminController{
		adminUsecase: adminUsecase,
		validator:    validator,
		tokenUtil:    tokenUtil,
	}
}

func (ac *adminController) Register(c echo.Context) error {
	log := logrus.New()
	request := new(dto.RegisterRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := ac.validator.Validate(request); err != nil {
		log.Error(err)
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := ac.adminUsecase.Register(c, request)
	if err != nil {
		log.Error(err)
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_CREATE_ADMIN)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.ADMIN_CREATED_SUCCESS, nil)
}

func (ac *adminController) Login(c echo.Context) error {
	request := new(dto.LoginRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := ac.validator.Validate(request); err != nil {
		fmt.Println(err)
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Menangkap kedua nilai kembalian dari Login
	response, err := ac.adminUsecase.Login(c, request)
	if err != nil {
		fmt.Println("Error: ", err)
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_LOGIN)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.LOGIN_SUCCESS, response)
}
