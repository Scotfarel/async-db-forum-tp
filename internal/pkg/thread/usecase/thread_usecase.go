package usecase

import (
	"strconv"

	"github.com/Scotfarel/db-tp-api/internal/pkg/forum"
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/post"
	"github.com/Scotfarel/db-tp-api/internal/pkg/thread"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"
	"github.com/Scotfarel/db-tp-api/internal/pkg/user"
	"github.com/Scotfarel/db-tp-api/internal/pkg/vote"

)

type ThreadUseCase struct {
	forumRepo  forum.Repository
	postRepo   post.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
	voteRepo   vote.Repository
}

func CreateUseCase(
	threadRepo thread.Repository,
	userRepo user.Repository,
	forumRepo forum.Repository,
	postRepo post.Repository,
	voteRepo vote.Repository,
	) thread.Usecase {
	return &ThreadUseCase{
		forumRepo:  forumRepo,
		postRepo:   postRepo,
		threadRepo: threadRepo,
		userRepo:   userRepo,
		voteRepo:   voteRepo,
	}
}

func (useCase *ThreadUseCase) InsertThreadInto(insertedThread *models.Thread) (*models.Thread, error) {
	forumIn, err := useCase.forumRepo.GetForumBySlug(insertedThread.Forum)
	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrForumDoesntExists
		}

		return nil, err
	}

	userByNickname, err := useCase.userRepo.GetUserByNickname(insertedThread.Author)
	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrUserDoesntExists
		}
	}

	if insertedThread.Slug != "" {
		threadBySlug, err := useCase.threadRepo.GetThreadBySlug(insertedThread.Slug)
		if err != nil && err != utils.ErrDoesntExists {
			return nil, err
		}
		if threadBySlug != nil {
			return threadBySlug, utils.ErrExistWithSlug
		}
	}

	insertedThread.Forum = forumIn.Slug
	insertedThread.Author = userByNickname.Nickname
	insertedThread.AuthorID = userByNickname.ID
	insertedThread.ForumID = forumIn.ID

	if err := useCase.threadRepo.InsertThreadInto(insertedThread); err != nil {
		return nil, err
	}

	return insertedThread, nil
}

func (useCase *ThreadUseCase) CreateThreadPosts(idOrSlug string, posts []*models.Post) ([]*models.Post, error) {
	creatingThread := &models.Thread{}

	id, err := strconv.ParseUint(idOrSlug, 10, 64)
	if err != nil {
		creatingThread, err = useCase.threadRepo.GetThreadBySlug(idOrSlug)
	} else {
		creatingThread, err = useCase.threadRepo.GetThreadByID(id)
	}

	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrThreadDoesntExists
		}

		return nil, err
	}

	_, err = useCase.userRepo.CheckNicknames(posts)
	if err != nil {
		return nil, err
	}

	_, err = useCase.postRepo.CheckPostParentPosts(posts, creatingThread.ID)
	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return []*models.Post{}, nil
	}

	for _, oldPost := range posts {
		oldPost.ThreadID = creatingThread.ID
		oldPost.Forum = creatingThread.Forum
		oldPost.ForumID = creatingThread.ForumID
	}
	if err = useCase.postRepo.InsertIntoPost(posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (useCase *ThreadUseCase) GetThreadBySlugOrID(idOrSlug string) (*models.Thread, error) {
	threadBySlugOrId := &models.Thread{}

	id, err := strconv.ParseUint(idOrSlug, 10, 64)
	if err != nil {
		threadBySlugOrId, err = useCase.threadRepo.GetThreadBySlug(idOrSlug)
	} else {
		threadBySlugOrId, err = useCase.threadRepo.GetThreadByID(id)
	}

	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrThreadDoesntExists
		}

		return nil, err
	}

	return threadBySlugOrId, nil
}

func (useCase *ThreadUseCase) UpdateThread(idOrSlug string, thread *models.Thread) (*models.Thread, error) {
	updatingThread := &models.Thread{}

	id, err := strconv.ParseUint(idOrSlug, 10, 64)
	if err != nil {
		updatingThread, err = useCase.threadRepo.GetThreadBySlug(idOrSlug)
	} else {
		updatingThread, err = useCase.threadRepo.GetThreadByID(id)
	}

	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrThreadDoesntExists
		}

		return nil, err
	}

	if thread.About != "" {
		updatingThread.About = thread.About
	}
	if thread.Title != "" {
		updatingThread.Title = thread.Title
	}

	if err = useCase.threadRepo.UpdateThread(updatingThread); err != nil {
		return nil, err
	}

	updatingThread.Votes, err = useCase.voteRepo.GetVotes(updatingThread.ID)
	if err != nil {
		return nil, err
	}

	return updatingThread, nil
}

func (useCase *ThreadUseCase) InsertVote(idOrSlug string, insertedVote *models.Vote) (*models.Thread, error) {
	insertedThread := &models.Thread{}

	id, err := strconv.ParseUint(idOrSlug, 10, 64)
	if err != nil {
		insertedThread, err = useCase.threadRepo.GetThreadBySlug(idOrSlug)
	} else {
		insertedThread, err = useCase.threadRepo.GetThreadByID(id)
	}

	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrThreadDoesntExists
		}

		return nil, err
	}

	userByNickname, err := useCase.userRepo.GetUserByNickname(insertedVote.Nickname)
	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrUserDoesntExists
		}
	}

	insertedVote.ThreadID = insertedThread.ID
	insertedVote.UserID = userByNickname.ID
	if err = useCase.voteRepo.AddVote(insertedVote); err != nil {
		return nil, err
	}

	insertedThread.Votes, err = useCase.voteRepo.GetVotes(insertedThread.ID)
	if err != nil {
		return nil, err
	}

	return insertedThread, nil
}

func (useCase *ThreadUseCase) GetThreadPosts(idOrSlug string, offset uint64, start uint64, sort string, ranging bool) ([]*models.Post, error) {
	postThread := &models.Thread{}

	id, err := strconv.ParseUint(idOrSlug, 10, 64)
	if err != nil {
		postThread, err = useCase.threadRepo.GetThreadBySlug(idOrSlug)
	} else {
		postThread, err = useCase.threadRepo.GetThreadByID(id)
	}

	if err != nil {
		if err == utils.ErrDoesntExists {
			return nil, utils.ErrThreadDoesntExists
		}

		return nil, err
	}

	returnPosts, err := useCase.postRepo.GetPostByThread(postThread.ID, offset, start, sort, ranging)
	if err != nil {
		return nil, err
	}

	return returnPosts, err
}
