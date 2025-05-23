package http

import (
	"net/http"
	"tugas-akhir/constant/status"
	dto "tugas-akhir/dto/base"

	"github.com/labstack/echo/v4"
)

func HandleErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, &dto.BaseResponse{
		Status:  status.STATUS_FAILED,
		Message: message,
	})
}	

func HandleSuccessResponse(c echo.Context, code int, message string, data any) error {
	return c.JSON(code, &dto.BaseResponse{
		Status:  status.STATUS_SUCCESS,
		Message: message,
		Data:    data,
	})
}

func HandleSearchResponse(c echo.Context, message string, data any, metadata *dto.MetadataResponse) error {
	return c.JSON(http.StatusOK, &dto.SearchResponse{
		BaseResponse: dto.BaseResponse{
			Status:  status.STATUS_SUCCESS,
			Message: message,
			Data:    data,
		},
		Metadata: metadata,
	})
}

func HandleLoadResponse(c echo.Context, message string, data any, metadata *dto.MetadataResponse) error {
	return c.JSON(http.StatusOK, &dto.LoadResponse{
		BaseResponse: dto.BaseResponse{
			Status:  status.STATUS_SUCCESS,
			Message: message,
			Data:    data,
		},
		Metadata: metadata,
	})
}

func HandlePaginationResponse(c echo.Context, message string, data any, pagination *dto.PaginationMetadata, link *dto.Link) error {
	return c.JSON(http.StatusOK, &dto.PaginationResponse{
		BaseResponse: dto.BaseResponse{
			Status:  status.STATUS_SUCCESS,
			Message: message,
			Data:    data,
		},
		Pagination: pagination,
		Link:       link,
	})
}
