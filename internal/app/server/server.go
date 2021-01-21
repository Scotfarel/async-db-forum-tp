package server


import (
	forumDelivery "github.com/Scotfarel/db-tp-api/internal/pkg/forum/delivery"
	postDelivery "github.com/Scotfarel/db-tp-api/internal/pkg/post/delivery"
	serviceDelivery "github.com/Scotfarel/db-tp-api/internal/pkg/service/delivery"
	threadDelivery "github.com/Scotfarel/db-tp-api/internal/pkg/thread/delivery"
	userDelivery "github.com/Scotfarel/db-tp-api/internal/pkg/user/delivery"

	"github.com/Scotfarel/db-tp-api/configs"

	forumRepo "github.com/Scotfarel/db-tp-api/internal/pkg/forum/repository"
	forumUseCase "github.com/Scotfarel/db-tp-api/internal/pkg/forum/usecase"
	postRepo "github.com/Scotfarel/db-tp-api/internal/pkg/post/repository"
	postUseCase "github.com/Scotfarel/db-tp-api/internal/pkg/post/usecase"
	serviceRepo "github.com/Scotfarel/db-tp-api/internal/pkg/service/repostitory"
	serviceUseCase "github.com/Scotfarel/db-tp-api/internal/pkg/service/usecase"
	threadRepo "github.com/Scotfarel/db-tp-api/internal/pkg/thread/repository"
	threadUseCase "github.com/Scotfarel/db-tp-api/internal/pkg/thread/usecase"
	userRepo "github.com/Scotfarel/db-tp-api/internal/pkg/user/repository"
	userUseCase "github.com/Scotfarel/db-tp-api/internal/pkg/user/usecase"
	voteRepo "github.com/Scotfarel/db-tp-api/internal/pkg/vote/repository"


	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var config = pgx.ConnConfig{
	Host:     configs.Host,
	Port:     configs.DBPort,
	Database: configs.Database,
	User:     configs.User,
	Password: configs.Password,
}

type Server struct {
	Url   string
}

func (server *Server) StartApiServer() {
	psqlConnector, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: config,
			MaxConnections: 10000,
		})
	if err != nil {
		logrus.Fatal(err)
	}

	router := echo.New()

	userRepository := userRepo.CreateRepository(psqlConnector)
	forumRepository := forumRepo.CreateRepository(psqlConnector)
	threadRepository := threadRepo.CreateRepository(psqlConnector)
	postRepository := postRepo.CreateRepository(psqlConnector)
	voteRepository := voteRepo.CreateRepository(psqlConnector)

	serviceRepository := serviceRepo.CreateRepository(psqlConnector)

	user := userUseCase.CreateUseCase(userRepository)
	forum := forumUseCase.CreateUseCase(forumRepository, userRepository, postRepository, threadRepository)
	thread := threadUseCase.CreateUseCase(threadRepository, userRepository, forumRepository, postRepository, voteRepository)
	service := serviceUseCase.CreateUseCase(serviceRepository)
	post := postUseCase.CreateUseCase(postRepository, forumRepository, voteRepository, threadRepository, userRepository)

	userDelivery.CreateHandler(router, user)
	forumDelivery.CreateHandler(router, forum, thread)
	threadDelivery.CreateHandler(router, thread)
	postDelivery.CreateHandler(router, post)

	serviceDelivery.CreateHandler(router, service)

	router.Start(server.Url)
}