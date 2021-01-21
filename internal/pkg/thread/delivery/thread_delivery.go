package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/thread"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	useCase thread.Usecase
}

func CreateHandler(routs *echo.Echo, useCase thread.Usecase) *Handler {
	handler := &Handler{
		useCase,
	}

	routs.POST("/api/thread/:slug_or_id/create", handler.CreatePost())
	routs.POST("/api/thread/:slug_or_id/details", handler.UpdateThread())
	routs.POST("/api/thread/:slug_or_id/vote", handler.ThreadVote())
	routs.GET("/api/thread/:slug_or_id/details", handler.GetThreadDetails())
	routs.GET("/api/thread/:slug_or_id/posts", handler.GetThreadPosts())

	return handler
}

func (handler *Handler) CreatePost() echo.HandlerFunc {
	type request struct {
		Author  string `json:"author" binding:"required"`
		Message string `json:"message" binding:"required"`
		Parent  uint64 `json:"parent" binding:"required"`
	}

	return func(context echo.Context) error {
		requestPost := []*request{}
		if err := json.NewDecoder(context.Request().Body).Decode(&requestPost); err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}
		slugOrID := context.Param("slug_or_id")

		posts := make([]*models.Post, 0, len(requestPost))
		for _, oldPost := range requestPost {
			post := &models.Post{
				Author:   oldPost.Author,
				Message:  oldPost.Message,
				ParentID: oldPost.Parent,
				IsEdited: false,
			}

			posts = append(posts, post)
		}

		returnPosts, err := handler.useCase.CreateThreadPosts(slugOrID, posts)
		if err != nil {
			if err == utils.ErrThreadDoesntExists || err == utils.ErrUserDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}

			return context.JSON(http.StatusConflict, utils.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.JSON(http.StatusCreated, returnPosts)
	}
}

func (handler *Handler) GetThreadDetails() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrID := context.Param("slug_or_id")

		returnThread, err := handler.useCase.GetThreadBySlugOrID(slugOrID)
		if err != nil {
			if err == utils.ErrThreadDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}

			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		return context.JSON(http.StatusOK, returnThread)
	}
}

func (handler *Handler) UpdateThread() echo.HandlerFunc {
	type updateThreadRequest struct {
		Message string `json:"message" binding:"require"`
		Title   string `json:"title" binding:"require"`
	}
	return func(context echo.Context) error {
		updateThreadRequest := &updateThreadRequest{}
		if err := context.Bind(updateThreadRequest); err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		slugOrID := context.Param("slug_or_id")
		threadModel := &models.Thread{
			About: updateThreadRequest.Message,
			Title: updateThreadRequest.Title,
		}

		updateThread, err := handler.useCase.UpdateThread(slugOrID, threadModel)
		if err != nil {
			if err == utils.ErrThreadDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.JSON(http.StatusOK, updateThread)
	}
}

func (handler *Handler) ThreadVote() echo.HandlerFunc {
	type voteReq struct {
		Nickname string `json:"nickname" binding:"require"`
		Voice    int64  `json:"voice" binding:"require"`
	}
	return func(context echo.Context) error {
		voteRequest := &voteReq{}
		if err := context.Bind(voteRequest); err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		slugOrID := context.Param("slug_or_id")

		votesR := &models.Vote{
			Nickname: voteRequest.Nickname,
			Voice:    voteRequest.Voice,
		}

		rThread, err := handler.useCase.InsertVote(slugOrID, votesR)
		if err != nil {
			if err == utils.ErrThreadDoesntExists || err == utils.ErrUserDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.JSON(http.StatusOK, rThread)
	}
}

func (handler *Handler) GetThreadPosts() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrID := context.Param("slug_or_id")
		offset := uint64(0)
		time := uint64(0)
		var err error
		if l := context.QueryParam("limit"); l != "" {
			offset, err = strconv.ParseUint(l, 10, 64)
			if err != nil {
				return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
		}
		if s := context.QueryParam("since"); s != "" {
			time, err = strconv.ParseUint(s, 10, 64)
			if err != nil {
				return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
		}

		if err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		sort := context.QueryParam("sort")
		desc := false
		if descVal := context.QueryParam("desc"); descVal == "true" {
			desc = true
		}

		returnPosts, err := handler.useCase.GetThreadPosts(slugOrID, offset, time, sort, desc)
		if err != nil {
			if err == utils.ErrThreadDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}

			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		return context.JSON(http.StatusOK, returnPosts)
	}
}
