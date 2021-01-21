package delivery

import (
	"fmt"

	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"
	"github.com/Scotfarel/db-tp-api/internal/pkg/user"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	useCase user.Usecase
}

func CreateHandler(routs *echo.Echo, useCase user.Usecase) *Handler {
	handler := &Handler{
		useCase: useCase,
	}

	routs.GET("/api/user/:nickname/profile", handler.GetProfile())
	routs.POST("/api/user/:nickname/create", handler.CreateUser())
	routs.POST("/api/user/:nickname/profile", handler.UpdateProfile())

	return handler
}

func (handler *Handler) CreateUser() echo.HandlerFunc {
	type userRequest struct {
		Email    string `json:"email" binding:"required" validate:"email"`
		Fullname string `json:"fullname" binding:"required"`
		About    string `json:"about"`
	}

	return func(c echo.Context) error {
		createRequest := &userRequest{}
		if err := c.Bind(createRequest); err != nil {
			logrus.Error(fmt.Errorf("Binding error %s", err))
			return c.JSON(http.StatusBadRequest, utils.ErrorResponce{err.Error()})
		}

		nick := c.Param("nickname")

		creatingUser := &models.User{
			Email:    createRequest.Email,
			Fullname: createRequest.Fullname,
			About:    createRequest.About,
		}

		inserted, err := handler.useCase.InsertUserInto(nick, creatingUser)
		if err != nil {
			if err == utils.ErrUserExistWith {
				return c.JSON(http.StatusConflict, inserted)
			}

			logrus.Error(fmt.Errorf("Request error %s", err))
			return c.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, inserted[0])
	}
}

func (handler *Handler) GetProfile() echo.HandlerFunc {
	return func(context echo.Context) error {
		nick := context.Param("nickname")

		returnUser, err := handler.useCase.GetUserByNickname(nick)
		if err != nil && err != utils.ErrDoesntExists {
			logrus.Error(fmt.Errorf("Request error %s", err))
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{err.Error()})
		}

		if err == utils.ErrDoesntExists {
			return context.JSON(http.StatusNotFound, utils.ErrorResponce{err.Error()})
		}

		return context.JSON(http.StatusOK, returnUser)
	}
}

func (handler *Handler) UpdateProfile() echo.HandlerFunc {
	type userRequset struct {
		Email    string `json:"email" binding:"required"`
		Fullname string `json:"fullname" binding:"required"`
		About    string `json:"about"`
	}

	return func(context echo.Context) error {
		updateRequest := &userRequset{}
		if err := context.Bind(updateRequest); err != nil {
			logrus.Error(fmt.Errorf("Binding error %s", err))
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{err.Error()})
		}

		nickname := context.Param("nickname")

		u := &models.User{
			Email:    updateRequest.Email,
			Fullname: updateRequest.Fullname,
			About:    updateRequest.About,
		}

		err := handler.useCase.UpdateUser(nickname, u)
		if err != nil {
			if err == utils.ErrUserExistWith {
				return context.JSON(http.StatusConflict, utils.ErrorResponce{err.Error()})
			}
			if err == utils.ErrUserDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{err.Error()})
			}
			logrus.Error(fmt.Errorf("Request error %s", err))
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{err.Error()})
		}

		return context.JSON(http.StatusOK, u)
	}
}
