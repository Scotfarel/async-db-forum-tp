package delivery

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/post"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	postUsecase post.Usecase
}

func CreateHandler(routes *echo.Echo, useCase post.Usecase) *Handler {
	handler := &Handler{
		postUsecase: useCase,
	}

	routes.GET("/api/post/:id/details", handler.GetPostDetails())
	routes.POST("/api/post/:id/details", handler.UpdatePost())

	return handler
}

func (handler *Handler) GetPostDetails() echo.HandlerFunc {
	return func(context echo.Context) error {
		id, err := strconv.ParseUint(context.Param("id"), 10, 64)
		if err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		related := strings.Split(context.QueryParam("related"), ",")

		returnPost, err := handler.postUsecase.GetPost(id, related)
		if err != nil {
			if err == utils.ErrPostDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}

			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		return context.JSON(http.StatusOK, returnPost)
	}
}

func (handler *Handler) UpdatePost() echo.HandlerFunc {
	type PostRequest struct {
		Message string `json:"message" binding:"require"`
	}
	return func(context echo.Context) error {
		id, err := strconv.ParseUint(context.Param("id"), 10, 64)
		if err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}
		req := &PostRequest{}
		if err = context.Bind(req); err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		postNew := &models.Post{
			ID:      id,
			Message: req.Message,
		}
		returnPost, err := handler.postUsecase.UpdatePost(postNew)
		if err != nil {
			if err == utils.ErrPostDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		return context.JSON(http.StatusOK, returnPost)
	}
}
