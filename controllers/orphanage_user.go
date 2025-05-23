package controllers

import (
	"net/http"
	"strconv"
	"strings"
	msg "tugas-akhir/constant/messages"
	dto_base "tugas-akhir/dto/base"
	dto "tugas-akhir/dto/orphanage"
	"tugas-akhir/usecases"
	http_util "tugas-akhir/utils/http"
	"tugas-akhir/utils/validation"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrphanageUserController struct {
	orphanageUserUsecase usecases.OrphanageUserUsecase
	validator            *validation.Validator
}

func NewOrphanageUserController(orphanageUserUsecase usecases.OrphanageUserUsecase, validator *validation.Validator) *OrphanageUserController {
	return &OrphanageUserController{
		orphanageUserUsecase: orphanageUserUsecase,
		validator:            validator,
	}
}

func (ouc *OrphanageUserController) CreateUser(ctx echo.Context) error {
	request := new(dto.OrphanageUserRequest)
	if err := ctx.Bind(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := ouc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := ouc.orphanageUserUsecase.CreateOrphanageUser(ctx, request)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_CREATE_USER)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusCreated, msg.SUCCESS_CREATE_USER, nil)

}

func (ouc *OrphanageUserController) GetOrphanageUserAll(ctx echo.Context) error {
    // Mengambil query parameter
    page := strings.TrimSpace(ctx.QueryParam("page"))
    limit := strings.TrimSpace(ctx.QueryParam("limit"))
    sortBy := ctx.QueryParam("sort_by")

    // Search dan filter parameter
    searchName := ctx.QueryParam("search_name")
    filterAddress := ctx.QueryParam("filter_address")
    filterEducation := ctx.QueryParam("filter_education")
    filterPosition := ctx.QueryParam("filter_position")
    filterAge := ctx.QueryParam("filter_age")

    // Mengonversi page dan limit ke integer
    intPage, intLimit, err := ouc.convertQueryParams(page, limit)
    if err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
    }

    // Membuat request pagination
    req := &dto_base.PaginationRequest{
        Page:   intPage,
        Limit:  intLimit,
        SortBy: sortBy,
    }

    // Validasi pagination request
    if err := ouc.validator.Validate(req); err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
    }

    // Memanggil usecase untuk mendapatkan data orphanage users dengan filter dan search
    result, metadata, link, err := ouc.orphanageUserUsecase.GetOrphanageUserAll(
        ctx, req, searchName, filterAddress, filterEducation, filterPosition, filterAge,
    )
    if err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_ORPHANAGE_USER_ALL)
    }

    // Mengembalikan response dengan data hasil
    return http_util.HandlePaginationResponse(ctx, msg.SUCCESS_GET_ORPHANAGE_USER_ALL, result, metadata, link)
}


func (ouc *OrphanageUserController) GetOrphanageUserByID(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	result, err := ouc.orphanageUserUsecase.GetOrphanageUserByID(ctx, userID)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_USER_BY_ID)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_USER_BY_ID, result)

}

func (ouc *OrphanageUserController) DeleteOrphanageUser(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := ouc.orphanageUserUsecase.DeleteOrphanageUser(ctx, userID); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_DELETE_USER)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_DELETE_USER, nil)

}

func (ouc *OrphanageUserController) UpdateOrphanageUser(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	request := new(dto.OrphanageUserRequest)
	if err := ctx.Bind(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := ouc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := ouc.orphanageUserUsecase.UpdateOrphanageUser(ctx, userID, request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPDATE_USER)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_UPDATE_USER, nil)

}

func (ouc *OrphanageUserController) GetOrphanageUserByPosition(ctx echo.Context) error {
	result, err := ouc.orphanageUserUsecase.GetOrphanageUserByPosition(ctx, ctx.Param("position"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_USER_BY_POSITION)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_USER_BY_POSITION, result)

}

func (ouc *OrphanageUserController) convertQueryParams(page, limit string) (int, int, error) {
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
