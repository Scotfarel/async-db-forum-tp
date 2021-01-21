package usecase

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/forum"
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/post"
	"github.com/Scotfarel/db-tp-api/internal/pkg/thread"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"
	"github.com/Scotfarel/db-tp-api/internal/pkg/user"
)

type UseCase struct {
	forumRepo  forum.Repository
	postRepo   post.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
}

func CreateUseCase(forumRepo forum.Repository, userRepo user.Repository, postRepo post.Repository, threadRepo thread.Repository) forum.Usecase {
	return &UseCase{
		forumRepo:  forumRepo,
		postRepo:   postRepo,
		threadRepo: threadRepo,
		userRepo:   userRepo,
	}
}

func (useCase *UseCase) InsertIntoForum(forumNew *models.Forum) (*models.Forum, error) {
	insertedForum, err := useCase.forumRepo.GetForumBySlug(forumNew.Slug)
	if err != nil && err != utils.ErrDoesntExists {
		return nil, err
	}
	if insertedForum != nil {
		return insertedForum, utils.ErrExistWithSlug
	}
	userNickname, err := useCase.userRepo.GetUserByNickname(forumNew.AdminNickname)
	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrUserDoesntExists
		}
		return nil, err
	}

	forumNew.AdminNickname = userNickname.Nickname
	forumNew.AdminID = userNickname.ID

	if err = useCase.forumRepo.InsertIntoForum(forumNew); err != nil {
		return nil, err
	}

	return forumNew, nil
}

func (useCase *UseCase) GetForumBySlug(slug string) (*models.Forum, error) {
	forumBySlug, err := useCase.forumRepo.GetForumBySlug(slug)
	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrForumDoesntExists
		}
		return nil, err
	}

	return forumBySlug, nil
}

func (useCase *UseCase) GetForumThreads(slug string, offset uint64, time string, desc bool) ([]*models.Thread, error) {
	if _, err := useCase.forumRepo.GetForumBySlug(slug); err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrForumDoesntExists
		}
		return nil, err
	}

	threadsBySlug, err := useCase.threadRepo.GetThreadByForumSlug(slug, offset, time, desc)
	if err != nil {
		return nil, err
	}

	return threadsBySlug, nil
}

func (useCase *UseCase) GetForumUsers(slug string, offset uint64, time string, desc bool) ([]*models.User, error) {
	oldForum, err := useCase.forumRepo.GetForumBySlug(slug)
	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrForumDoesntExists
		}
	}
	forumUsers, err := useCase.userRepo.GetUsersByForum(oldForum.ID, offset, time, desc)
	if err != nil {
		return nil, err
	}

	return forumUsers, nil
}
