package usecases

import (
	"context"
	"fmt"
	"math"
	"strconv"
	dto_base "tugas-akhir/dto/base"
	dto "tugas-akhir/dto/donation"
	"tugas-akhir/entities"
	"tugas-akhir/repositories"
	err_util "tugas-akhir/utils/error"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ProgramDonationUsecase interface {
	CreateProgramDonation(c echo.Context, req *dto.ProgramDonationRequest) error
	GetProgramDonationByID(c echo.Context, programDonationID uuid.UUID) (*dto.ProgramDonationResponse, error)
	GetProgramDonationAll(c echo.Context, req *dto_base.PaginationRequest, searchTitle string) (*[]dto.ProgramDonationResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	UpdateProgramDonation(c echo.Context, programDonationID uuid.UUID, req *dto.ProgramDonationRequest) error
	DeleteProgramDonation(c echo.Context, programDonationID uuid.UUID) error
	GetDashboardData(c echo.Context) (*dto.DashboardDonationResponse, error)
}

type programDonationUsecase struct {
	programDonationRepo repositories.ProgramDonationRepository
}

func NewProgramDonationUsecase(programDonationRepo repositories.ProgramDonationRepository) ProgramDonationUsecase {
	return &programDonationUsecase{
		programDonationRepo: programDonationRepo,
	}
}

func (pdu *programDonationUsecase) CreateProgramDonation(c echo.Context, req *dto.ProgramDonationRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	programID := uuid.New()

	// Log the program donation ID
	logger := logrus.New()
	logger.Infof("Creating program donation with ID: %s", programID.String())

	programImages := make([]entities.ProgramDonationImage, len(req.ProgramDonationImages))
	for i, image := range req.ProgramDonationImages {
		programImages[i] = entities.ProgramDonationImage{
			ID:        uuid.New(),
			ProgramID: programID,
			ImageUrl:  image.ImageUrl,
		}
	}

	programDonation := &entities.ProgramDonation{
		ID:            programID,
		Title:         req.Title,
		Deskripsi:     req.Deskripsi,
		GoalAmount:    req.GoalAmount,
		CurrentAmount: 0,
		DonationImage: programImages,
	}

	// Log program donation details before saving
	logger.Infof("Program donation details: %+v", programDonation)

	// Call repository to save program donation
	err := pdu.programDonationRepo.CreateProgramDonation(ctx, programDonation)
	if err != nil {
		logger.Error("Failed to create program donation in repository: ", err)
		return err
	}

	logger.Info("Program donation created successfully")
	return nil
}

func (pdu *programDonationUsecase) GetProgramDonationByID(c echo.Context, programDonationID uuid.UUID) (*dto.ProgramDonationResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	programDonation, err := pdu.programDonationRepo.GetProgramDonationByID(ctx, programDonationID)
	if err != nil {
		return nil, err
	}

	var programDonationImages []dto.ProgramDonationImageResponse
	for _, image := range programDonation.DonationImage {
		programDonationImages = append(programDonationImages, dto.ProgramDonationImageResponse{
			ImageUrl: image.ImageUrl,
		})
	}

	responses := dto.ProgramDonationResponse{
		ID:                    programDonation.ID.String(),
		Title:                 programDonation.Title,
		Deskripsi:             programDonation.Deskripsi,
		GoalAmount:            programDonation.GoalAmount,
		CurrentAmount:         programDonation.CurrentAmount,
		ProgramDonationImages: programDonationImages,
	}

	return &responses, nil
}

func (pdu *programDonationUsecase) GetProgramDonationAll(c echo.Context, req *dto_base.PaginationRequest, searchTitle string) (*[]dto.ProgramDonationResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx := c.Request().Context() // Menggunakan context dari request Echo

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

	// Mendapatkan data program donasi dan total data dengan searchTitle
	programDonations, totalData, err := pdu.programDonationRepo.GetProgramDonationAll(ctx, req, searchTitle)
	if err != nil {
		return nil, nil, nil, err
	}

	// Membentuk response dari data yang diperoleh
	var responses []dto.ProgramDonationResponse
	for i, p := range programDonations {
		var programDonationImages []dto.ProgramDonationImageResponse
		for _, image := range p.DonationImage {
			programDonationImages = append(programDonationImages, dto.ProgramDonationImageResponse{
				ImageUrl: image.ImageUrl,
			})
		}

		responses = append(responses, dto.ProgramDonationResponse{
			Number:                req.Limit*(req.Page-1) + i + 1,
			ID:                    p.ID.String(),
			Title:                 p.Title,
			Deskripsi:             p.Deskripsi,
			GoalAmount:            p.GoalAmount,
			CurrentAmount:         p.CurrentAmount,
			ProgramDonationImages: programDonationImages,
		})
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(req.Limit)))
	paginationMetadata := &dto_base.PaginationMetadata{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: req.Page,
	}

	// Menangani kasus ketika halaman yang diminta lebih besar dari total halaman
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

func (pdu *programDonationUsecase) UpdateProgramDonation(c echo.Context, programDonationID uuid.UUID, req *dto.ProgramDonationRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	programDonation, err := pdu.programDonationRepo.GetProgramDonationByID(ctx, programDonationID)
	if err != nil {
		return err
	}

	programDonation.Title = req.Title
	programDonation.Deskripsi = req.Deskripsi
	programDonation.GoalAmount = req.GoalAmount

	err = pdu.programDonationRepo.UpdateProgramDonation(ctx, programDonation)
	if err != nil {
		return err
	}

	return nil
}

func (pdu *programDonationUsecase) DeleteProgramDonation(c echo.Context, programDonationID uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	err := pdu.programDonationRepo.DeleteProgramDonation(ctx, programDonationID)
	if err != nil {
		return err
	}

	return nil
}

func (pdu *programDonationUsecase) GetDashboardData(c echo.Context) (*dto.DashboardDonationResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	programCount, totalDonation, uniqueDonatorsCount, err := pdu.programDonationRepo.GetDashboardData(ctx)
	if err != nil {
		return nil, err
	}

	// Membuat response untuk dashboard
	response := &dto.DashboardDonationResponse{
		ProgramCount:        programCount,
		TotalDonation:       totalDonation,
		UniqueDonatorsCount: uniqueDonatorsCount,
	}

	return response, nil
}
