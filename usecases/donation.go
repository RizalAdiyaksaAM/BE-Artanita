package usecases

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
	"tugas-akhir/config"
	dto_base "tugas-akhir/dto/base"
	dto "tugas-akhir/dto/donation"
	"tugas-akhir/entities"
	"tugas-akhir/repositories"
	err_util "tugas-akhir/utils/error"
	midtran "tugas-akhir/utils/midtrans"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type DonationUsecase interface {
	CreateDonation(c echo.Context, request dto.DonationRequest) (dto.DonationResponse, error)
	UpdateDonationStatus(c echo.Context, notification midtran.Notification) error
	GetDonations(c echo.Context, programDonationID uuid.UUID, searchName string, req *dto_base.PaginationRequest) (*[]dto.DonationResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetDonationByID(c echo.Context, donationID uuid.UUID) (*dto.DonationResponse, error)
    GetDonationLanding(c echo.Context) (*[]dto.DonationLandingResponse, error)
    GetChartDonation(c echo.Context) (*[]dto.DonationChartResponse, error)
    GetDonaturNotifikasi(c echo.Context) (*[]dto.DonaturNotifikasiResponse, error)
    GetDonaturByProgramDonation(c echo.Context, programDonationID uuid.UUID) (*[]dto.DonationLandingResponse, error)
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{}) // Atau &logrus.TextFormatter{}
	logrus.SetOutput(os.Stdout) // Default output ke stdout
	logrus.SetLevel(logrus.InfoLevel) // Set level logging
}

type donationUsecase struct {
	donationRepository                repositories.DonationRepository
	config                            config.MidtransConfig
	programDonation                   repositories.ProgramDonationRepository
	transactionNotificationRepository repositories.TransactionNotificationRepository
}

func NewDonationUsecase(donationRepository repositories.DonationRepository, config config.MidtransConfig, programDonation repositories.ProgramDonationRepository, transactionNotificationRepository repositories.TransactionNotificationRepository) DonationUsecase {
	return &donationUsecase{
		donationRepository:                donationRepository,
		config:                            config,
		programDonation:                   programDonation,
		transactionNotificationRepository: transactionNotificationRepository,
	}
}

func (d *donationUsecase) CreateDonation(c echo.Context, request dto.DonationRequest) (dto.DonationResponse, error) {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Initialize logger
	log := logrus.New()

       if request.ProgramID == uuid.Nil {
        firstProgram, err := d.programDonation.GetFirstProgramDonation(ctx)
        if err != nil {
            log.WithError(err).Error("Failed to get first program donation")
            return dto.DonationResponse{}, errors.New("failed to get default program donation")
        }
        request.ProgramID = firstProgram.ID
    }


	// Initialize the donation transaction from request
	transactionDonation := entities.Donation{
		ID:                uuid.New(),
		Name:              request.Name,
		Address:           request.Address,
		NoWA:              request.NoWA,
		Email:             request.Email,
		Message:           request.Message,
		Status:            0, // default to 0 (pending)
		SnapURL:           "",
		ProgramDonationID: request.ProgramID, // Setting ProgramDonationID directly from the request
		Amount:            request.Amount,    // Set the Amount to the value from the request
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Verifikasi ProgramDonation yang dipilih (pastikan program ID valid)
	program, err := d.programDonation.GetProgramDonationByID(ctx, request.ProgramID)
	if err != nil {
		log.WithError(err).Error("Program not found")
		return dto.DonationResponse{}, errors.New("program donation not found")
	}

	// Create Midtrans request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionDonation.ID.String(),
			GrossAmt: int64(request.Amount), // Ensure that Amount is greater than 0
		},
	}

	var client snap.Client
	client.New(d.config.ServerKey, midtrans.Sandbox)

	snapResp, err := client.CreateTransaction(req)
	if snapResp == nil {
		return dto.DonationResponse{}, err
	}

	transactionDonation.SnapURL = snapResp.RedirectURL

	// Insert the transaction donation into the repository
	err = d.donationRepository.CreateDonation(ctx, &transactionDonation)
	if err != nil {
		log.WithError(err).Error("Failed to insert donation transaction into the database")
		return dto.DonationResponse{}, err
	}

	// Return the response struct with donation details, including the SnapURL
	return dto.DonationResponse{
		ID:        transactionDonation.ID.String(),
		Name:      transactionDonation.Name,
		Address:   transactionDonation.Address,
		NoWA:      transactionDonation.NoWA,
		Email:     transactionDonation.Email,
		Amount:    transactionDonation.Amount,
		Message:   transactionDonation.Message,
		Status:    transactionDonation.Status,
		SnapURL:   transactionDonation.SnapURL, // Ensure SnapURL is returned
		ProgramID: program.ID.String(),
	}, nil
}

func (d *donationUsecase) UpdateDonationStatus(c echo.Context, notification midtran.Notification) error {
    log := logrus.New()
    log.Infof("Processing notification: %+v", notification)

    // Mengonversi OrderID (string) menjadi uuid.UUID
    donationID, err := uuid.Parse(notification.OrderID)
    if err != nil {
        log.WithError(err).Error("Failed to parse donation ID")
        return errors.New("invalid donation ID")
    }

    // Cari donasi berdasarkan ID (uuid.UUID)
    donation, err := d.donationRepository.FindById(c.Request().Context(), donationID)
    if err != nil {
        log.WithError(err).Error("Donation not found")
        return errors.New("donation not found")
    }

    log.Infof("Found donation: %+v", donation)

    // Tentukan status berdasarkan transaction_status
    previousStatus := donation.Status
    switch notification.TransactionStatus {
    case "capture", "settlement":
        donation.Status = 1 // Pembayaran berhasil
        // Update currentAmount di ProgramDonation
        if err := d.programDonation.UpdateCurrentAmount(c.Request().Context(), donation.ProgramDonationID, donation.Amount); err != nil {
            log.WithError(err).Error("Failed to update currentAmount in ProgramDonation")
            return errors.New("failed to update currentAmount")
        }
    case "pending":
        donation.Status = 0 // Pembayaran pending
    case "deny", "expire", "cancel", "failure":
        donation.Status = 2 // Pembayaran gagal
    default:
        log.Warnf("Unknown transaction status: %s", notification.TransactionStatus)
    }

    log.Infof("Updating donation status from %d to %d", previousStatus, donation.Status)

    // Simpan status donasi yang baru
    err = d.donationRepository.Update(c.Request().Context(), donation)
    if err != nil {
        log.WithError(err).Error("Failed to update donation status")
        return err
    }

    // Simpan notifikasi ke dalam database untuk audit
    transactionNotification := entities.TransactionNotification{
        ID:                uuid.New(),
        OrderID:           notification.OrderID,
        TransactionStatus: notification.TransactionStatus,
        GrossAmount:       notification.GrossAmount,
        TransactionTime:   notification.TransactionTime,
        SignatureKey:      notification.SignatureKey,
    }

    err = d.transactionNotificationRepository.CreateNotification(c.Request().Context(), &transactionNotification)
    if err != nil {
        log.WithError(err).Error("Failed to save transaction notification")
    }

    log.Infof("Donation status updated successfully: %+v", donation)
    return nil
}




func (du *donationUsecase) GetDonations(c echo.Context, programDonationID uuid.UUID, searchName string, req *dto_base.PaginationRequest) (*[]dto.DonationResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
    ctx, cancel := context.WithCancel(c.Request().Context())
    defer cancel()

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

    // Validasi parameter
    donations, totalData, err := du.donationRepository.GetDonations(ctx, programDonationID, searchName, req)
    if err != nil {
        return nil, nil, nil, err
    }

    // Convert entities to DTO response dan menambahkan `number` serta `programTitle`
    var donationResponses []dto.DonationResponse
    for i, donation := range donations {
        // Mengambil ProgramDonation berdasarkan ProgramDonationID dari setiap donasi
        var programDonation *entities.ProgramDonation
        if donation.ProgramDonationID != uuid.Nil {
            programDonation, err = du.programDonation.GetProgramDonationByID(ctx, donation.ProgramDonationID)
            if err != nil {
                return nil, nil, nil, err // Jika terjadi error, kembalikan error
            }
        }

        // Pastikan programDonation ada sebelum mengakses Title
        programTitle := ""
        if programDonation != nil {
            programTitle = programDonation.Title
        }

        donationResponses = append(donationResponses, dto.DonationResponse{
            Number:      req.Limit*(req.Page-1) + i + 1, // Menghitung nomor urut donasi
            ID:          donation.ID.String(),
            Name:        donation.Name,
            Address:     donation.Address,
            NoWA:        donation.NoWA,
            Email:       donation.Email,
            ProgramTitle: programTitle, // Menambahkan title dari ProgramDonation
            ProgramID:   donation.ProgramDonationID.String(),
            Amount:      donation.Amount,
            Status:      donation.Status,
            Message:     donation.Message,
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

    return &donationResponses, paginationMetadata, link, nil
}


func (du *donationUsecase) GetDonationByID(c echo.Context, donationID uuid.UUID) (*dto.DonationResponse, error) {
    ctx, cancel := context.WithCancel(c.Request().Context())
    defer cancel()

    // Mengambil donasi berdasarkan ID
    donation, err := du.donationRepository.GetDonationByID(ctx, donationID)
    if err != nil {
        return nil, err
    }

    // Mengonversi donasi ke bentuk DTO response
    donationResponse := &dto.DonationResponse{
		ID:      donation.ID.String(),
		Name:    donation.Name,
		Address: donation.Address,
		NoWA:    donation.NoWA,
		Email:   donation.Email,
		ProgramID: donation.ProgramDonationID.String(),
		Amount:  donation.Amount,
		Status:  donation.Status,
		Message: donation.Message,
    }

    return donationResponse, nil
}

func (du *donationUsecase) GetDonationLanding(c echo.Context) (*[]dto.DonationLandingResponse, error) {
	// Mengambil konteks dari Echo
	ctx := c.Request().Context()

	// Memanggil repository untuk mendapatkan donasi dengan status 1
	donations, err := du.donationRepository.GetDonationsLanding(ctx)
	if err != nil {
		return nil, err
	}

	// Anda mungkin ingin mengonversi data entities.Donation ke DTO (misalnya, DonationLandingResponse)
	var response []dto.DonationLandingResponse
	for _, donation := range *donations {
		response = append(response, dto.DonationLandingResponse{
            ID:        donation.ID.String(),
			Name:      donation.Name,
			Amount:    donation.Amount,
			Message:   donation.Message,
		})
	}

	return &response, nil
}

// GetChartDonation mengambil data donasi untuk chart
func (du *donationUsecase) GetChartDonation(c echo.Context) (*[]dto.DonationChartResponse, error) {
	// Mengambil konteks dari Echo
	ctx := c.Request().Context()

    // Memanggil repository untuk mendapatkan donasi dengan status 1
    donations, err := du.donationRepository.GetDonation(ctx)
    if err != nil {
        return nil, err
    }

    var response []dto.DonationChartResponse
    for _, donation := range *donations {

        if donation.Status != 1 {
            continue
        }

        program, err := du.programDonation.GetProgramDonationByID(ctx, donation.ProgramDonationID)
        if err != nil {
            return nil, err
        }

        response = append(response, dto.DonationChartResponse{
            ID:          donation.ID.String(),
            Amount:      donation.Amount,
            Date:        donation.UpdatedAt.Format("2006-01-02"),
            ProgramDonation: program.Title,
        })
    }

    return &response, nil
}


func (du *donationUsecase) GetDonaturNotifikasi(c echo.Context) (*[]dto.DonaturNotifikasiResponse, error) {
    // Mengambil konteks dari Echo
    ctx := c.Request().Context()

    // Memanggil repository untuk mendapatkan donasi
    donations, err := du.donationRepository.GetNotifikasi(ctx)
    if err != nil {
        return nil, err
    }

    // Anda mungkin ingin mengonversi data entities.Donation ke DTO (misalnya, DonaturNotifikasiResponse)
    var response []dto.DonaturNotifikasiResponse
    for _, donation := range *donations {
        if donation.TransactionStatus != "settlement" {
            continue
        }

        programID, err := uuid.Parse(donation.OrderID)
        if err != nil {
            return nil, errors.New("ID program donasi tidak valid: " + donation.OrderID)
        }

        donatiUser, err := du.donationRepository.GetDonationByID(ctx, programID)
        if err != nil {
            return nil, err
        }

        program, err := du.programDonation.GetProgramDonationByID(ctx, donatiUser.ProgramDonationID)
        if err != nil {
            return nil, err
        }

        response = append(response, dto.DonaturNotifikasiResponse{
            ID:        donation.ID.String(),
            Name:      donatiUser.Name,
            Amount:    donation.GrossAmount,
            ProgramDonation: program.Title,
            Message:   donatiUser.Message,
            Date:      donation.CreatedAt.Format("2006-01-02"),
            Status:    donation.TransactionStatus,
        })

    }

    return &response, nil
}

func (du *donationUsecase) GetDonaturByProgramDonation(c echo.Context, programDonationID uuid.UUID) (*[]dto.DonationLandingResponse, error) {
    // Mengambil konteks dari Echo
    ctx := c.Request().Context()

    // Memanggil repository untuk mendapatkan donasi
    donations, err := du.donationRepository.GetDonationByProgramID(ctx, programDonationID)
    if err != nil {
        return nil, err
    }


    // Anda mungkin ingin mengonversi data entities.Donation ke DTO (misalnya, DonationLandingResponse)
    var response []dto.DonationLandingResponse
    for _, donation := range *donations {
        response = append(response, dto.DonationLandingResponse{
            ID:        donation.ID.String(),
            Name:      donation.Name,
            Amount:    donation.Amount,
            Message:   donation.Message,
        })
    }

    return &response, nil
}
