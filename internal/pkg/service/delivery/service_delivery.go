package delivery

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/service"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"

	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	serviceUsecase service.Usecase
}

func CreateHandler(routs *echo.Echo, useCase service.Usecase) *Handler {
	handler := &Handler{
		serviceUsecase: useCase,
	}

	routs.GET("/api/service/status", handler.GetStatus())
	routs.POST("/api/service/clear", handler.Clear())

	return handler
}

func (handler *Handler) GetStatus() echo.HandlerFunc {
	return func(context echo.Context) error {
		stat, err := handler.serviceUsecase.GetInfo()
		if err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.JSON(http.StatusOK, stat)
	}
}

func (handler *Handler) Clear() echo.HandlerFunc {
	return func(context echo.Context) error {
		err := handler.serviceUsecase.Drop()
		if err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.NoContent(http.StatusOK)
	}
}
