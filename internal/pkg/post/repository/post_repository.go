package repository

import (
	"fmt"
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/post"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"
	"github.com/jackc/pgx"
	"strconv"
)

type PostRepository struct {
	db *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) post.Repository {
	return &PostRepository{
		db: db,
	}
}

func (repository *PostRepository) InsertIntoPost(posts []*models.Post) error {
	query := "INSERT INTO posts (author, forum, message, parent, thread) VALUES "

	for _, p := range posts {
		query += fmt.Sprintf(
			"(%d, %d, '%s', %d, %d),",
			p.AuthorID,
			p.ForumID,
			p.Message,
			p.ParentID,
			p.ThreadID,
		)
		_, err := repository.db.Exec(
			"INSERT INTO forums_users (user_id, forum_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			p.AuthorID,
			p.ForumID,
		)
		if err != nil {
			return err
		}
	}
	query = query[0 : len(query)-1]
	query += " RETURNING id, created"

	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	rows, err := tx.Query(query)
	if err != nil {
		return err
	}

	defer func() {
		rows.Close()
	}()

	index := 0
	for rows.Next() {
		if err := rows.Scan(&posts[index].ID, &posts[index].CreationDate); err != nil {
			return err
		}
		index++
	}
	return tx.Commit()
}

func (repository *PostRepository) GetPostCountByForumID(id uint64) (uint64, error) {
	var num uint64
	if err := repository.db.QueryRow("SELECT count(*) from posts WHERE forum = $1", id).
		Scan(&num); err != nil {
		return 0, err
	}

	return num, nil
}

func (repository *PostRepository) GetPostByID(id uint64) (*models.Post, error) {
	postById := &models.Post{}
	if err := repository.db.QueryRow(
		"SELECT p.id, u.nickname, f.slug, p.thread, p.message, p.created, p.isEdited, coalesce(path[array_length(path, 1) - 1], 0) FROM posts AS p JOIN users AS u ON (u.id = p.author) JOIN forums AS f ON (f.id = p.forum) WHERE p.id = $1",
		id,
	).
		Scan(
		&postById.ID,
		&postById.Author,
		&postById.Forum,
		&postById.ThreadID,
		&postById.Message,
			&postById.CreationDate,
		&postById.IsEdited,
		&postById.ParentID,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.ErrDoesntExists
		}
		return nil, err
	}

	return postById, nil
}

func (repository *PostRepository) UpdatePost(post *models.Post) error {
	if _, err := repository.db.Exec("UPDATE posts SET message = $2, isEdited = TRUE WHERE id = $1",
		post.ID,
		post.Message,
	); err != nil {
		return err
	}

	return nil
}

func (repository *PostRepository) GetPostByThread(id uint64, offset uint64, time uint64, sort string, desc bool) ([]*models.Post, error) {
	queryStringFmt := "SELECT p.id, u.nickname, f.slug, p.thread, p.created, p.message, p.isEdited, coalesce(p.path[array_length(p.path, 1) - 1], 0) FROM posts AS p JOIN users AS u ON (u.id = p.author) JOIN forums AS f ON (f.id = p.forum) WHERE %s %s"

	var where string
	var order string

	switch sort {
	case "flat", "":
		where = "p.thread = $1"
		if time != 0 {
			if desc {
				where += " AND p.id < $2"
			} else {
				where += " AND p.id > $2"
			}
		}
		order = "ORDER BY "
		if sort == "flat" {
			order += "p.created"
			if desc {
				order += " DESC"
			}
			order += ", p.id"
			if desc {
				order += " DESC"
			}
		} else {
			order += "p.id"
			if desc {
				order += " DESC"
			}
		}
		if offset != 0 {
			if time != 0 {
				order += " LIMIT $3"
			} else {
				order += " LIMIT $2"
			}
		}
	case "tree":
		where = "p.thread = $1"
		if time != 0 {
			if desc {
				where += " AND coalesce(path < (select path FROM posts where id = $2), true)"
			} else {
				where += " AND coalesce(path > (select path FROM posts where id = $2), true)"
			}
		}
		order = "ORDER BY p.path[1]"
		if desc {
			order += " DESC"
		}
		order += ", p.path[2:]"
		if desc {
			order += " DESC"
		}
		order += " NULLS FIRST"
		if offset != 0 {
			if time != 0 {
				order += " LIMIT $3"
			} else {
				order += " LIMIT $2"
			}
		}
	case "parent_tree":
		where = "p.path[1] IN (SELECT path[1] FROM posts WHERE thread = $1 AND array_length(path, 1) = 1"
		if time != 0 {
			if desc {
				where += " AND id < (SELECT path[1] FROM posts WHERE id = $2)"
			} else {
				where += " AND id > (SELECT path[1] FROM posts WHERE id = $2)"
			}
		}

		where += " ORDER BY id"
		if desc {
			where += " DESC"
		}

		if offset != 0 {
			if time != 0 {
				where += " LIMIT $3"
			} else {
				where += " LIMIT $2"
			}
		}
		where += ")"
		order = "ORDER BY p.path[1]"
		if desc {
			order += " DESC"
		}
		order += ", p.path[2:] NULLS FIRST"
	}

	rows := &pgx.Rows{}
	var err error
	if time != 0 {
		if offset != 0 {
			rows, err = repository.db.Query(fmt.Sprintf(queryStringFmt, where, order), id, time, offset)
		} else {
			rows, err = repository.db.Query(fmt.Sprintf(queryStringFmt, where, order), id, time)
		}
	} else {
		if offset != 0 {
			rows, err = repository.db.Query(fmt.Sprintf(queryStringFmt, where, order), id, offset)
		} else {
			rows, err = repository.db.Query(fmt.Sprintf(queryStringFmt, where, order), id)
		}
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	returnPosts := []*models.Post{}

	for rows.Next() {
		pet := &models.Post{}

		err = rows.Scan(
			&pet.ID,
			&pet.Author,
			&pet.Forum,
			&pet.ThreadID,
			&pet.CreationDate,
			&pet.Message,
			&pet.IsEdited,
			&pet.ParentID,
		)
		if err != nil {
			return nil, err
		}

		returnPosts = append(returnPosts, pet)
	}

	return returnPosts, nil
}

func (repository *PostRepository) CheckPostParentPosts(posts []*models.Post, threadID uint64) (bool, error) {
	links := map[uint64]uint64{}
	values := []interface{}{threadID}

	query := "SELECT count(*) FROM posts WHERE thread = $1 AND id in ("

	i := 2
	for _, oldPost := range posts {
		if oldPost.ParentID > 0 {
			query += "$" + strconv.Itoa(i) + ","
			links[oldPost.ParentID] += 1
			values = append(values, oldPost.ParentID)
			i++
		}
	}

	if len(links) == 0 {
		return true, nil
	}

	query = query[0:len(query)-1] + ")"
	var num int
	if err := repository.db.QueryRow(query, values...).Scan(&num); err != nil {
		return false, err
	}
	if num != len(links) {
		return false, utils.ErrParentPostDoesntExists
	}
	return true, nil
}
