package usecase

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"
	"github.com/Scotfarel/db-tp-api/internal/pkg/user"
)

type UserUsecase struct {
	userRepo user.Repository
}

func CreateUseCase(ur user.Repository) user.Usecase {
	return &UserUsecase{
		ur,
	}
}

func (useCase *UserUsecase) InsertUserInto(nick string, user *models.User) ([]*models.User, error) {
	nickUser, err := useCase.userRepo.GetUserByNickname(nick)
	if err != nil && err != utils.ErrDoesntExists {
		return nil, err
	}
	emailUser, err := useCase.userRepo.GetUserByEmail(user.Email)
	if err != nil && err != utils.ErrDoesntExists {
		return nil, err
	}

	if nickUser != nil || emailUser != nil {
		matchedUsers := []*models.User{}
		if nickUser != nil {
			matchedUsers = append(matchedUsers, nickUser)
			if emailUser != nil && nickUser.Nickname != emailUser.Nickname {
				matchedUsers = append(matchedUsers, emailUser)
			}
		} else if emailUser != nil {
			matchedUsers = append(matchedUsers, emailUser)
		}
		return matchedUsers, utils.ErrUserExistWith
	}

	user.SetNickname(nick)
	if err = useCase.userRepo.InsertUserInto(user); err != nil {
		return nil, err
	}

	return []*models.User{user}, nil
}

func (useCase *UserUsecase) GetUserByNickname(nick string) (*models.User, error) {
	userByNickname, err := useCase.userRepo.GetUserByNickname(nick)
	if err != nil {
		return nil, err
	}

	return userByNickname, nil
}

func (useCase *UserUsecase) UpdateUser(nick string, user *models.User) error {
	userByNickname, err := useCase.userRepo.GetUserByNickname(nick)
	if err != nil {
		if err == utils.ErrDoesntExists {
			return utils.ErrUserDoesntExists
		}
		return err
	}

	newEmailCheckUser, err := useCase.userRepo.GetUserByEmail(user.Email)
	if err != nil && err != utils.ErrDoesntExists {
		return err
	}
	if err != utils.ErrDoesntExists && newEmailCheckUser.Nickname != userByNickname.Nickname {
		return utils.ErrUserExistWith
	}

	user.SetNickname(nick)
	if user.Email == "" {
		user.Email = userByNickname.Email
	}
	if user.Fullname == "" {
		user.Fullname = userByNickname.Fullname
	}
	if user.About == "" {
		user.About = userByNickname.About
	}

	if err = useCase.userRepo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}
