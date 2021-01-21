package usecase

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/forum"
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/post"
	"github.com/Scotfarel/db-tp-api/internal/pkg/thread"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"
	"github.com/Scotfarel/db-tp-api/internal/pkg/user"
	"github.com/Scotfarel/db-tp-api/internal/pkg/vote"
)

type UseCase struct {
	forumRepo  forum.Repository
	postRepo   post.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
	voteRepo   vote.Repository
}

func CreateUseCase(postRepo post.Repository, forumRepo forum.Repository, voteRepo vote.Repository, threadRepo thread.Repository, userRepo user.Repository) post.Usecase {
	return &UseCase{
		forumRepo:  forumRepo,
		postRepo:   postRepo,
		threadRepo: threadRepo,
		userRepo:   userRepo,
		voteRepo:   voteRepo,
	}
}

func (useCase *UseCase) GetPost(id uint64, related []string) (*models.PostFull, error) {
	postById, err := useCase.postRepo.GetPostByID(id)
	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrPostDoesntExists
		}
		return nil, err
	}

	newPost := &models.PostFull{PostData: postById}

	for _, types := range related {
		switch types {
		case "user":
			userByNickname, err := useCase.userRepo.GetUserByNickname(postById.Author)
			if err != nil {
				if err == utils.ErrDoesntExists {
					return nil, utils.ErrUserDoesntExists
				}
				return nil, err
			}
			newPost.Author = userByNickname
		case "thread":
			threadByID, err := useCase.threadRepo.GetThreadByID(postById.ThreadID)
			if err != nil {
				if err == utils.ErrDoesntExists {
					return nil, utils.ErrThreadDoesntExists
				}
				return nil, err
			}
			newPost.Thread = threadByID
		case "forum":
			forumBySlug, err := useCase.forumRepo.GetForumBySlug(postById.Forum)
			if err != nil {
				if err == utils.ErrDoesntExists {
					return nil, utils.ErrForumDoesntExists
				}
				return nil, err
			}
			newPost.Forum = forumBySlug
		}
	}

	return newPost, nil
}

func (useCase *UseCase) UpdatePost(newPost *models.Post) (*models.Post, error) {
	updatedPost, err := useCase.postRepo.GetPostByID(newPost.ID)
	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrPostDoesntExists
		}

		return nil, err
	}
	if newPost.Message == "" || newPost.Message == updatedPost.Message {
		return updatedPost, nil
	}
	updatedPost.Message = newPost.Message
	updatedPost.IsEdited = true

	if err = useCase.postRepo.UpdatePost(updatedPost); err != nil {
		return nil, err
	}

	return updatedPost, nil
}
