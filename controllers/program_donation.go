package controllers

import (
	"net/http"
	"strconv"
	"strings"
	msg "tugas-akhir/constant/messages"
	"tugas-akhir/drivers/cloudinary"
	dto_base "tugas-akhir/dto/base"
	dto "tugas-akhir/dto/donation"
	"tugas-akhir/usecases"
	http_util "tugas-akhir/utils/http"
	"tugas-akhir/utils/validation"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ProgramDonationController struct {
	programDonationUsecase usecases.ProgramDonationUsecase
	validator              *validation.Validator
	cloudinaryService      cloudinary.CloudinaryService
}

func NewProgramDonationController(programDonationUsecase usecases.ProgramDonationUsecase, validator *validation.Validator, cloudinaryService cloudinary.CloudinaryService) *ProgramDonationController {
	return &ProgramDonationController{
		programDonationUsecase: programDonationUsecase,
		validator:              validator,
		cloudinaryService:      cloudinaryService,
	}
}

func (pdc *ProgramDonationController) CreateProgramDonation(ctx echo.Context) error {
	var logger = logrus.New()
	logger.Info("Received request to create program donation")

	// Parse multipart form data
	form, err := ctx.MultipartForm()
	if err != nil {
		logger.Error("Failed to parse multipart form: ", err)
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	var request dto.ProgramDonationRequest
	// Bind request body to DTO
	if err := ctx.Bind(&request); err != nil {
		logger.Error("Failed to bind request data: ", err)
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Log bound request for debugging
	logger.Infof("Bound request: %+v", request)

	// Check if any images are uploaded
	files := form.File["image"]
	if len(files) == 0 {
		logger.Warn("No images uploaded in request")
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Handle image upload
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			logger.Error("Failed to open image file: ", err)
			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
		}

		secureURL, err := pdc.cloudinaryService.UploadImage(ctx.Request().Context(), src, "artanita/programDonation/image")
		if err != nil {
			logger.Error("Failed to upload image to Cloudinary: ", err)
			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
		}

		// Log the image URL
		logger.Infof("Uploaded image URL: %s", secureURL)

		image := dto.ProgramDonationImageRequest{
			ImageUrl: &secureURL,
		}

		// Append uploaded image URL to the request
		request.ProgramDonationImages = append(request.ProgramDonationImages, image)
	}

	// Validate the request data
	if err := pdc.validator.Validate(request); err != nil {
		logger.Error("Request validation failed: ", err)
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Log the validated request for debugging
	logger.Infof("Validated request: %+v", request)

	// Call the usecase to create the program donation
	if err := pdc.programDonationUsecase.CreateProgramDonation(ctx, &request); err != nil {
		logger.Error("Failed to create program donation: ", err)
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_CREATE_PROGRAM_DONATION)
	}

	// Log success and return response
	logger.Info("Successfully created program donation")
	return http_util.HandleSuccessResponse(ctx, http.StatusCreated, msg.SUCCESS_CREATE_PROGRAM_DONATION, nil)
}

func (pdc *ProgramDonationController) GetProgramDonationAll(ctx echo.Context) error {
    page := strings.TrimSpace(ctx.QueryParam("page"))
    limit := strings.TrimSpace(ctx.QueryParam("limit"))
    sortBy := ctx.QueryParam("sort_by")
    searchTitle := ctx.QueryParam("search_title") // Mengambil query parameter search_title

    intPage, intLimit, err := pdc.convertQueryParams(page, limit)
    if err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
    }

    req := &dto_base.PaginationRequest{
        Page:   intPage,
        Limit:  intLimit,
        SortBy: sortBy,
    }

    if err := pdc.validator.Validate(req); err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
    }

    // Memanggil usecase untuk mendapatkan program donasi dengan search atau filter
    result, metadata, link, err := pdc.programDonationUsecase.GetProgramDonationAll(ctx, req, searchTitle)
    if err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_PROGRAM_DONATION_ALL)
    }

    return http_util.HandlePaginationResponse(ctx, msg.SUCCESS_GET_PROGRAM_DONATION_ALL, result, metadata, link)
}


func (pdc *ProgramDonationController) GetProgramDonationById(ctx echo.Context) error {
	programDonationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	result, err := pdc.programDonationUsecase.GetProgramDonationByID(ctx, programDonationID)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_PROGRAM_DONATION_BY_ID)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_PROGRAM_DONATION_BY_ID, result)
}

func (pdc *ProgramDonationController) UpdateProgramDonation(ctx echo.Context) error {
	programDonationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	request := new(dto.ProgramDonationRequest)
	if err := ctx.Bind(&request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := pdc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := pdc.programDonationUsecase.UpdateProgramDonation(ctx, programDonationID, request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPDATE_PROGRAM_DONATION)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_UPDATE_PROGRAM_DONATION, nil)
}

func (pdc *ProgramDonationController) DeleteProgramDonation(ctx echo.Context) error {
	programDonationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := pdc.programDonationUsecase.DeleteProgramDonation(ctx, programDonationID); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_DELETE_PROGRAM_DONATION)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_DELETE_PROGRAM_DONATION, nil)
}

func (pdc *ProgramDonationController) convertQueryParams(page, limit string) (int, int, error) {
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

func (pdc *ProgramDonationController) GetDashboard(ctx echo.Context) error {
	var logger = logrus.New()
	logger.Info("Received request to get dashboard data")

	// Memanggil usecase untuk mendapatkan data dashboard
	data, err := pdc.programDonationUsecase.GetDashboardData(ctx)
	if err != nil {
		logger.Error("Failed to get dashboard data: ", err)
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_DASHBOARD)
	}

	// Mengirimkan response yang berisi data dashboard
	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_DASHBOARD, data)
}
