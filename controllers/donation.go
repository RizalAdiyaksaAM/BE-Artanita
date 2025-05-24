package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	msg "tugas-akhir/constant/messages"
	dto_base "tugas-akhir/dto/base"
	dto "tugas-akhir/dto/donation"
	"tugas-akhir/usecases"
	http_util "tugas-akhir/utils/http"
	midtran "tugas-akhir/utils/midtrans"
	"tugas-akhir/utils/validation"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type DonationController struct {
	donationUsecase usecases.DonationUsecase
	validator       *validation.Validator
	midtransClient  *midtran.Client
}

func NewDonationController(donationUsecase usecases.DonationUsecase, validator *validation.Validator, midtransClient *midtran.Client) *DonationController {
	return &DonationController{
		donationUsecase: donationUsecase,
		validator:       validator,
		midtransClient:  midtransClient, // Menyuntikkan client Midtrans
	}
}
func (d *DonationController) CreateDonation(c echo.Context) error {
	log := logrus.New()

	// Bind request data to the DTO
	request := dto.DonationRequest{}
	if err := c.Bind(&request); err != nil {
		log.WithError(err).Error("Error binding request data")
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid request data")
	}

	// Validate the request data using validator
	if err := d.validator.Validate(request); err != nil {
		log.WithError(err).Error("Validation failed for donation request")
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid data in the request")
	}

	// Call usecase to create the donation
	response, err := d.donationUsecase.CreateDonation(c, request)
	if err != nil {
		log.WithError(err).Error("Error creating donation in usecase")
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, "Failed to create donation transaction")
	}

	// Return the success response with the donation details
	log.Infof("Donation created successfully: %+v", response)
	return http_util.HandleSuccessResponse(c, http.StatusCreated, "Donation created successfully", response)
}

func (d *DonationController) MidtransWebhook(c echo.Context) error {
    log := logrus.New()
    log.Info("Received webhook from Midtrans")

    // Cetak raw body untuk debugging
    bodyBytes, err := io.ReadAll(c.Request().Body)
    if err != nil {
        log.WithError(err).Error("Failed to read request body")
        return c.JSON(http.StatusOK, map[string]interface{}{"status": "error", "message": "Failed to read request body"})
    }

    c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

    log.Infof("Raw webhook payload: %s", string(bodyBytes))

    var notification midtran.Notification
    if err := c.Bind(&notification); err != nil {
        log.WithError(err).Error("Error binding webhook data")
        return c.JSON(http.StatusOK, map[string]interface{}{"status": "error", "message": "Invalid request format"})
    }

    log.Infof("Notification data: %+v", notification)

    if notification.SignatureKey != "" {
        if !d.midtransClient.VerifyNotificationSignature(notification) {
            log.Error("Invalid signature from Midtrans webhook")
            return c.JSON(http.StatusOK, map[string]interface{}{"status": "error", "message": "Invalid signature"})
        }
    } else {
        log.Warn("Signature key is empty, skipping verification")
    }

    // Memanggil usecase untuk memproses status donasi dan update currentAmount
    err = d.donationUsecase.UpdateDonationStatus(c, notification)
    if err != nil {
        log.WithError(err).Error("Failed to update donation status")
        return c.JSON(http.StatusOK, map[string]interface{}{"status": "error", "message": "Failed to update donation status"})
    }

    log.Info("Webhook processed successfully")
    return c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "message": "Webhook processed successfully"})
}

func (d *DonationController) GetDonations(ctx echo.Context) error {
	programDonationIDStr := ctx.QueryParam("program_donation_id")
	searchName := ctx.QueryParam("search_name")
	page := strings.TrimSpace(ctx.QueryParam("page"))
	limit := strings.TrimSpace(ctx.QueryParam("limit"))
	sortBy := ctx.QueryParam("sort_by")

	// Parse ProgramDonationID jika ada
	var programDonationID uuid.UUID
	if programDonationIDStr != "" {
		parsedID, err := uuid.Parse(programDonationIDStr)
		if err != nil {
			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid Program Donation ID format")
		}
		programDonationID = parsedID
	}

	intPage, intLimit, err := d.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	// Panggil usecase untuk mendapatkan donasi dengan filter atau semua donasi
	result, metadata, link, err := d.donationUsecase.GetDonations(ctx, programDonationID, searchName, req)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_DONATIONS)
	}

	return http_util.HandlePaginationResponse(ctx, msg.SUCCESS_GET_DONATIONS, result, metadata, link)
}

func (d *DonationController) convertQueryParams(page, limit string) (int, int, error) {
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

func (dc *DonationController) GetDonationByID(ctx echo.Context) error {
    // Mendapatkan parameter donationID dari URL
    donationIDStr := ctx.Param("id")
    donationID, err := uuid.Parse(donationIDStr)
    if err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid Donation ID format")
    }

    // Memanggil usecase untuk mendapatkan donasi berdasarkan ID
    donation, err := dc.donationUsecase.GetDonationByID(ctx, donationID)
    if err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_DONATIONS)
    }

    // Mengirimkan response dengan data donasi
    return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_DONATIONS, donation)
}

func (d *DonationController) GetDonationsLanding(ctx echo.Context) error {
	// Memanggil usecase untuk mendapatkan semua donasi yang berstatus 1
	donations, err := d.donationUsecase.GetDonationLanding(ctx)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_DONATIONS)
	}

	// Mengirimkan response dengan data donasi
	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_DONATIONS, donations)
}


func (d *DonationController) GetChartDonation(ctx echo.Context) error {
	donations, err := d.donationUsecase.GetChartDonation(ctx)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_DONATIONS)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_DONATIONS, donations)
}

func (d *DonationController) GetDonaturNotifikasi(ctx echo.Context) error {
	donations, err := d.donationUsecase.GetDonaturNotifikasi(ctx)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_DONATIONS)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_DONATIONS, donations)
}

func (d *DonationController) GetDonaturByProgramDonation(ctx echo.Context) error {
	programDonationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid Program Donation ID format")
	}

	donations, err := d.donationUsecase.GetDonaturByProgramDonation(ctx, programDonationID)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_DONATIONS)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_DONATIONS, donations)
}

func (d *DonationController) GetNotifikasi(ctx echo.Context) error {
	donationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid Donation ID format")
	}

	notifications, err := d.donationUsecase.GetNotifikasiStatus(ctx, donationID)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_GET_DONATIONS)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_GET_DONATIONS, notifications)
}