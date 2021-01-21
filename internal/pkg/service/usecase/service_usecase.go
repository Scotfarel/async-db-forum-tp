package usecase

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/service"
)

type UseCase struct {
	serviceRepo service.Repository
}

func CreateUseCase(sr service.Repository) service.Usecase {
	return &UseCase{
		serviceRepo: sr,
	}
}

func (useCase *UseCase) GetInfo() (*models.Status, error) {
	countForum, err := useCase.serviceRepo.ForumInfo()
	countPost, err := useCase.serviceRepo.PostInfo()
	countThread, err := useCase.serviceRepo.ThreadInfo()
	countUser, err := useCase.serviceRepo.UserInfo()

	if err != nil {
		return nil, err
	}

	info := &models.Status{
		ForumsCount:  countForum,
		PostsCount:   countPost,
		ThreadsCount: countThread,
		UsersCount:   countUser,
	}
	return info, nil
}

func (useCase *UseCase) Drop() error {
	err := useCase.serviceRepo.DropVotes()
	err = useCase.serviceRepo.DropPost()
	err = useCase.serviceRepo.DropThread()
	err = useCase.serviceRepo.DropForum()
	err = useCase.serviceRepo.DropUser()

	if err != nil {
		return err
	}

	return nil
}
