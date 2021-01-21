package delivery

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/forum"
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/thread"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	forumCase  forum.Usecase
	threadCase thread.Usecase
}

func CreateHandler(routes *echo.Echo, forumUseCase forum.Usecase, threadUseCase thread.Usecase) *Handler {
	handler := &Handler{
		forumCase:  forumUseCase,
		threadCase: threadUseCase,
	}

	routes.POST("/api/forum/create", handler.CreateForum())
	routes.POST("/api/forum/:fslug/create", handler.CreateThread())
	routes.GET("/api/forum/:slug/details", handler.GetForumDetails())
	routes.GET("/api/forum/:slug/threads", handler.GetForumThreads())
	routes.GET("/api/forum/:slug/users", handler.GetForumUsers())

	return handler
}

func (handler *Handler) CreateForum() echo.HandlerFunc {
	type forumRequest struct {
		Slug  string `json:"slug" binding:"required"`
		Title string `json:"title" binding:"required"`
		User  string `json:"user" binding:"required"`
	}

	return func(context echo.Context) error {
		forumRequest := &forumRequest{}
		if err := context.Bind(forumRequest); err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		request := &models.Forum{
			Slug:          forumRequest.Slug,
			Title:         forumRequest.Title,
			AdminNickname: forumRequest.User,
		}

		returnForum, err := handler.forumCase.InsertIntoForum(request)
		if err != nil {
			if err == utils.ErrExistWithSlug {
				return context.JSON(http.StatusConflict, returnForum)
			}
			if err == utils.ErrUserDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.JSON(http.StatusCreated, returnForum)
	}
}

func (handler *Handler) CreateThread() echo.HandlerFunc {
	type CreateThreadRequest struct {
		Author  string    `json:"author" binding:"require"`
		Created time.Time `json:"created" binding:"omitempty"`
		Message string    `json:"message" binding:"require"`
		Title   string    `json:"title" binding:"require"`
		Slug    string    `json:"slug" binding:"omitempty"`
	}
	return func(context echo.Context) error {
		threadRequest := &CreateThreadRequest{}
		if err := context.Bind(threadRequest); err != nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		if _, err := strconv.ParseInt(threadRequest.Slug, 10, 64); err == nil {
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: utils.ErrIncorrectSlug.Error(),
			})
		}

		slug := context.Param("fslug")

		if threadRequest.Created.IsZero() {
			threadRequest.Created = time.Now()
		}

		requested := &models.Thread{
			Author:       threadRequest.Author,
			CreationDate: threadRequest.Created,
			About:        threadRequest.Message,
			Title:        threadRequest.Title,
			Slug:         threadRequest.Slug,
			Forum:        slug,
		}

		returnThread, err := handler.threadCase.InsertThreadInto(requested)
		if err != nil {
			if err == utils.ErrForumDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
			if err == utils.ErrUserDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
			if err == utils.ErrExistWithSlug {
				return context.JSON(http.StatusConflict, returnThread)
			}
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		return context.JSON(http.StatusCreated, returnThread)
	}
}

func (handler *Handler) GetForumDetails() echo.HandlerFunc {
	return func(context echo.Context) error {
		slug := context.Param("slug")

		forumDetails, err := handler.forumCase.GetForumBySlug(slug)
		if err != nil {
			if err == utils.ErrForumDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}

			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}

		return context.JSON(http.StatusOK, forumDetails)
	}
}

func (handler *Handler) GetForumThreads() echo.HandlerFunc {
	return func(context echo.Context) error {
		slug := context.Param("slug")
		offset := uint64(0)
		var err error

		if l := context.QueryParam("limit"); l != "" {
			offset, err = strconv.ParseUint(l, 10, 64)
			if err != nil {
				return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
		}
		timeParam := context.QueryParam("since")

		desc := false
		if descVal := context.QueryParam("desc"); descVal == "true" {
			desc = true
		}

		forumThreads, err := handler.forumCase.GetForumThreads(slug, offset, timeParam, desc)
		if err != nil {
			if err == utils.ErrForumDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.JSON(http.StatusOK, forumThreads)
	}
}

func (handler *Handler) GetForumUsers() echo.HandlerFunc {
	return func(context echo.Context) error {
		slug := context.Param("slug")

		offset := uint64(0)
		var err error
		if l := context.QueryParam("limit"); l != "" {
			offset, err = strconv.ParseUint(l, 10, 64)
			if err != nil {
				return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
		}
		since := context.QueryParam("since")

		desc := false
		if descVal := context.QueryParam("desc"); descVal == "true" {
			desc = true
		}

		forumsUsers, err := handler.forumCase.GetForumUsers(slug, offset, since, desc)
		if err != nil {
			if err == utils.ErrForumDoesntExists {
				return context.JSON(http.StatusNotFound, utils.ErrorResponce{
					Message: err.Error(),
				})
			}
			return context.JSON(http.StatusBadRequest, utils.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.JSON(http.StatusOK, forumsUsers)
	}
}
