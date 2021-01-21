package repository

import (
	"fmt"
	"strings"

	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"
	"github.com/Scotfarel/db-tp-api/internal/pkg/user"

	"github.com/jackc/pgx"
)

type UserRepository struct {
	db *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) user.Repository {
	return &UserRepository{
		db,
	}
}

func (repository *UserRepository) InsertUserInto(user *models.User) error {
	if _, err := repository.db.Exec(
		"INSERT INTO users (nickname, email, fullname, about) VALUES ($1, $2, $3, $4)",
		user.Nickname,
		user.Email,
		user.Fullname,
		user.About,
	); err != nil {
		return err
	}

	return nil
}

func (repository *UserRepository) GetUserByNickname(nickname string) (*models.User, error) {
	userByNickname := &models.User{}

	if err := repository.db.QueryRow(
		"SELECT id, nickname, email, fullname, about FROM users WHERE lower(nickname) = lower($1)",
		nickname,
	).Scan(
		&userByNickname.ID,
		&userByNickname.Nickname,
		&userByNickname.Email,
		&userByNickname.Fullname,
		&userByNickname.About,
		); err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.ErrDoesntExists
		}
		return nil, err
	}

	return userByNickname, nil
}

func (repository *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	userByEmail := &models.User{}

	if err := repository.db.QueryRow(
		"SELECT id, nickname, email, fullname, about FROM users WHERE lower(email) = lower($1)",
		email,
	).Scan(
		&userByEmail.ID,
		&userByEmail.Nickname,
		&userByEmail.Email,
		&userByEmail.Fullname,
		&userByEmail.About,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.ErrDoesntExists
		}
		return nil, err
	}

	return userByEmail, nil
}

func (repository *UserRepository) UpdateUser(user *models.User) error {
	if _, err := repository.db.Exec(
		"UPDATE users SET email = $2, fullname = $3, about = $4 WHERE lower(nickname) = lower($1)",
		user.Nickname,
		user.Email,
		user.Fullname,
		user.About,
	); err != nil {
		return err
	}

	return nil
}

func (repository *UserRepository) GetUsersByForum(id uint64, limit uint64, since string, desc bool) ([]*models.User, error) {
	usersByForum := []*models.User{}

	query := "SELECT u.nickname, u.email, u.fullname, u.about FROM forums_users fu JOIN users u ON (fu.user_id = u.id) WHERE fu.forum_id = $1"
	groupBy := " ORDER BY lower(u.nickname)"
	if desc {
		groupBy += " DESC"
	}

	if limit != 0 {
		groupBy += fmt.Sprintf(" LIMIT %d", limit)
	}

	var rows *pgx.Rows
	var err error
	if since != "" {
		if desc {
			query += " AND lower(u.nickname) < lower($2)"
		} else {
			query += " AND lower(u.nickname) > lower($2)"
		}
		rows, err = repository.db.Query(
			query + groupBy,
			id,
			since,
			)
	} else {
		rows, err = repository.db.Query(
			query + groupBy,
			id,
		)
	}
	if err != nil {
		//if err == sql.ErrNoRows {
		//	return nil, nil
		//}

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		userForum := &models.User{}
		err = rows.Scan(
			&userForum.Nickname,
			&userForum.Email,
			&userForum.Fullname,
			&userForum.About,
			)
		if err != nil {
			return nil, err
		}

		usersByForum = append(usersByForum, userForum)
	}

	return usersByForum, nil
}

func (repository *UserRepository) CheckNicknames(posts []*models.Post) (bool, error) {
	rows, err := repository.db.Query("SELECT id, lower(nickname) FROM users")
	if err != nil {
		return false, err
	}

	defer rows.Close()

	nicks := make(map[string]uint64)
	for rows.Next() {
		n := ""
		var id uint64
		if err := rows.Scan(
			&id,
			&n,
		); err != nil {
			return false, err
		}

		nicks[n] = id
	}

	for _, post := range posts {
		id := nicks[strings.ToLower(post.Author)]
		if id == 0 {
			return false, utils.ErrUserDoesntExists
		}
		post.AuthorID = id
	}

	return true, nil
}
