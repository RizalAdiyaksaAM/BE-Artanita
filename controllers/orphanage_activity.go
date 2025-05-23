package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	msg "tugas-akhir/constant/messages"
	"tugas-akhir/drivers/cloudinary"
	dto_base "tugas-akhir/dto/base"
	dto "tugas-akhir/dto/orphanage"
	"tugas-akhir/usecases"
	http_util "tugas-akhir/utils/http"
	"tugas-akhir/utils/validation"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrphanageActivityController struct {
	orphanageActivityUsecase usecases.OrphanageActivityUsecase
	validator                *validation.Validator
	cloudinaryService        cloudinary.CloudinaryService
}

func NewOrphanageActivityController(orphanageActivityUsecase usecases.OrphanageActivityUsecase, validator *validation.Validator, cloudinaryService cloudinary.CloudinaryService) *OrphanageActivityController {
	return &OrphanageActivityController{
		orphanageActivityUsecase: orphanageActivityUsecase,
		validator:                validator,
		cloudinaryService:        cloudinaryService,
	}
}

func (oac *OrphanageActivityController) CreateActivity(ctx echo.Context) error {
	log.Println("Received request to create activity")

	// Reading multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Println("Error reading multipart form:", err)
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}
	log.Println("Multipart form read successfully")

	// Binding request to DTO
	var request dto.ActivityRequest
	if err := ctx.Bind(&request); err != nil {
		log.Println("Error binding request:", err)
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}
	log.Printf("Request bound successfully: %+v", request)

	// Process Image files
	files := form.File["image"]
	if len(files) == 0 {
		log.Println("No image files provided in the request")
	}
	for _, file := range files {
		log.Printf("Processing image file: %s\n", file.Filename)
		src, err := file.Open()
		if err != nil {
			log.Println("Error opening image file:", err)
			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "1")
		}
		defer src.Close()

		secureURL, err := oac.cloudinaryService.UploadImage(ctx.Request().Context(), src, "artanita/activity/images")
		if err != nil {
			log.Println("Error uploading image to Cloudinary:", err)
			return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPLOAD_IMAGE)
		}
		log.Printf("Image uploaded successfully, URL: %s", secureURL)

		images := dto.ActivityImageRequest{
			ImageUrl: &secureURL,
		}
		request.ActivityImages = append(request.ActivityImages, images)
	}

	// Process Video files
	files = form.File["video"]
	if len(files) == 0 {
		log.Println("No video files provided in the request")
	}
	for _, file := range files {
		log.Printf("Processing video file: %s\n", file.Filename)
		src, err := file.Open()
		if err != nil {
			log.Println("Error opening video file:", err)
			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "2")
		}
		defer src.Close()

		secureURL, err := oac.cloudinaryService.UploadImage(ctx.Request().Context(), src, "artanita/activity/videos")
		if err != nil {
			log.Println("Error uploading video to Cloudinary:", err)
			return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPLOAD_IMAGE)
		}
		log.Printf("Video uploaded successfully, URL: %s", secureURL)

		videos := dto.ActivityVideoRequest{
			VideoUrl: &secureURL,
		}
		request.ActivityVideos = append(request.ActivityVideos, videos)
	}

	// Validate request
	log.Println("Validating request data")
	if err := oac.validator.Validate(request); err != nil {
		log.Println("Validation failed:", err)
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "3")
	}
	log.Println("Validation successful")

	// Call CreateActivity usecase
	log.Println("Calling CreateActivity usecase")
	if err := oac.orphanageActivityUsecase.CreateActivity(ctx, &request); err != nil {
		log.Println("Error creating activity:", err)
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_CREATE_ACTIVITY)
	}

	log.Println("Activity created successfully")
	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_CREATE_ACTIVITY, nil)
}

func (oac *OrphanageActivityController) GetActivityAll(ctx echo.Context) error {
    // Mengambil query parameter
    page := strings.TrimSpace(ctx.QueryParam("page"))
    limit := strings.TrimSpace(ctx.QueryParam("limit"))
    sortBy := ctx.QueryParam("sort_by")

    searchTitle := ctx.QueryParam("search_title") // Mengambil query parameter search_title

    intPage, intLimit, err := oac.convertQueryParams(page, limit)
    if err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
    }

    req := &dto_base.PaginationRequest{
        Page:   intPage,
        Limit:  intLimit,
        SortBy: sortBy,
    }

    if err := oac.validator.Validate(req); err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
    }

    // Memanggil usecase untuk mendapatkan aktivitas dengan search_title
    result, metadata, link, err := oac.orphanageActivityUsecase.GetActivityAll(ctx, req, searchTitle)
    if err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_ACTIVITY_ALL)
    }

    return http_util.HandlePaginationResponse(ctx, msg.SUCCESS_GET_ACTIVITY_ALL, result, metadata, link)
}


func (oac *OrphanageActivityController) GetActivityById(ctx echo.Context) error {
	activityID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	result, err := oac.orphanageActivityUsecase.GetActivityByID(ctx, activityID)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_ACTIVITY_BY_ID)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_ACTIVITY_BY_ID, result)
}

func (oac *OrphanageActivityController) UpdateActivity(ctx echo.Context) error {
	activityID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	request := new(dto.ActivityRequest)
	if err := ctx.Bind(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := oac.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := oac.orphanageActivityUsecase.UpdateActivity(ctx, activityID, request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPDATE_ACTIVITY)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_UPDATE_ACTIVITY, nil)
}

func (oac *OrphanageActivityController) DeleteActivity(ctx echo.Context) error {
	activityID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := oac.orphanageActivityUsecase.DeleteActivity(ctx, activityID); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_DELETE_ACTIVITY)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_DELETE_ACTIVITY, nil)

}

func (oac *OrphanageActivityController) convertQueryParams(page, limit string) (int, int, error) {
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
