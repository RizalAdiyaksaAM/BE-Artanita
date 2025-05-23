package usecases

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	dto_base "tugas-akhir/dto/base"
	dto "tugas-akhir/dto/orphanage"
	"tugas-akhir/entities"
	"tugas-akhir/repositories"
	err_util "tugas-akhir/utils/error"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrphanageActivityUsecase interface {
	CreateActivity(c echo.Context, req *dto.ActivityRequest) error
	GetActivityAll(c echo.Context, req *dto_base.PaginationRequest, searchTitle string) (*[]dto.ActivityResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetActivityByID(c echo.Context, id uuid.UUID) (*dto.ActivityResponse, error)
	UpdateActivity(c echo.Context, id uuid.UUID, req *dto.ActivityRequest) error
	DeleteActivity(c echo.Context, id uuid.UUID) error
}

type orphanageActivityUsecase struct {
	activityRepo repositories.OrphanageActivityRepository
}

func NewOrphanageActivityUsecase(activityRepo repositories.OrphanageActivityRepository) OrphanageActivityUsecase {
	return &orphanageActivityUsecase{
		activityRepo: activityRepo,
	}
}

func (uu *orphanageActivityUsecase) CreateActivity(c echo.Context, req *dto.ActivityRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// ActivityID generated here
	activityID := uuid.New()

	// Initialize arrays for images and videos
	activityImages := make([]entities.ActivityImage, len(req.ActivityImages))
	for i, image := range req.ActivityImages {
		activityImages[i] = entities.ActivityImage{
			ID:         uuid.New(),
			ActivityID: activityID,
			ImageUrl:   image.ImageUrl,
		}
	}

	activityVideos := make([]entities.ActivityVideo, len(req.ActivityVideos))
	for i, video := range req.ActivityVideos {
		activityVideos[i] = entities.ActivityVideo{
			ID:         uuid.New(),
			ActivityID: activityID,
			VideoUrl:   video.VideoUrl,
		}
	}

	// Create activity entity
	activity := &entities.OrphanageActivity{
		ID:             activityID,
		Title:          req.Title,
		Description:    req.Description,
		Location:       req.Location,
		Time:           req.Time,
		ActivityImages: activityImages,
		ActivityVideos: activityVideos,
	}

	// Call the repository to save the activity
	err := uu.activityRepo.CreateActivity(ctx, activity)
	if err != nil {
		log.Println("Error creating activity:", err)
		return err
	}

	log.Printf("ActivityImages URLs: %+v", activityImages)
	log.Printf("ActivityVideos URLs: %+v", activityVideos)

	return nil
}

func (uu *orphanageActivityUsecase) GetActivityAll(
    c echo.Context, req *dto_base.PaginationRequest, searchTitle string) (*[]dto.ActivityResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {

    ctx := c.Request().Context()

    baseURL := fmt.Sprintf(
        "%s?limit=%d&page=",
        c.Request().URL.Path,
        req.Limit,
    )

    var (
        next string
        prev string
    )

    // Pastikan page tidak lebih kecil dari 1
    if req.Page < 1 {
        req.Page = 1
    }

    if req.Page > 1 {
        prev = baseURL + strconv.Itoa(req.Page-1)
    }

    // Mendapatkan data aktivitas dan total data dengan searchTitle
    activities, totalData, err := uu.activityRepo.GetActivityAll(ctx, req, searchTitle)
    if err != nil {
        return nil, nil, nil, err
    }

    var responses []dto.ActivityResponse
    for _, a := range activities {

        var activityImages []dto.ActivityImageResponse
        for _, image := range a.ActivityImages {
            activityImages = append(activityImages, dto.ActivityImageResponse{
                ImageUrl: image.ImageUrl,
            })
        }

        var activityVideos []dto.ActivityVideoResponse
        for _, video := range a.ActivityVideos {
            activityVideos = append(activityVideos, dto.ActivityVideoResponse{
                VideoUrl: video.VideoUrl,
            })
        }

        responses = append(responses, dto.ActivityResponse{
            ID:             a.ID.String(),
            Title:          a.Title,
            Description:    a.Description,
			Location:       a.Location,
			Time:           a.Time,
            ActivityImages: activityImages,
            ActivityVideos: activityVideos,
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
    } else {
        next = baseURL + strconv.Itoa(req.Page+1)
    }

    link := &dto_base.Link{
        Next: next,
        Prev: prev,
    }

    return &responses, paginationMetadata, link, nil
}


func (uu *orphanageActivityUsecase) GetActivityByID(c echo.Context, id uuid.UUID) (*dto.ActivityResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	activity, err := uu.activityRepo.GetActivityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var activityImages []dto.ActivityImageResponse
	for _, image := range activity.ActivityImages {
		activityImages = append(activityImages, dto.ActivityImageResponse{
			ImageUrl: image.ImageUrl,
		})
	}

	var activityVideos []dto.ActivityVideoResponse
	for _, video := range activity.ActivityVideos {
		activityVideos = append(activityVideos, dto.ActivityVideoResponse{
			VideoUrl: video.VideoUrl,
		})
	}

	responses := dto.ActivityResponse{
		ID:             activity.ID.String(),
		Title:          activity.Title,
		Description:    activity.Description,
		Location:       activity.Location,
		Time:           activity.Time,
		ActivityImages: activityImages,
		ActivityVideos: activityVideos,
	}

	return &responses, nil
}

func (uu *orphanageActivityUsecase) UpdateActivity(c echo.Context, id uuid.UUID, req *dto.ActivityRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	activity, err := uu.activityRepo.GetActivityByID(ctx, id)
	if err != nil {
		return err
	}

	activity.Title = req.Title
	activity.Description = req.Description
	activity.Location = req.Location
	activity.Time = req.Time

	err = uu.activityRepo.UpdateActivity(ctx, id, activity)
	if err != nil {
		return err
	}

	return nil
}

func (uu *orphanageActivityUsecase) DeleteActivity(c echo.Context, id uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	err := uu.activityRepo.DeleteActivity(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
