package repository

import (
	"fmt"

	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/thread"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"

	"github.com/jackc/pgx"
)

type Repository struct {
	db *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) thread.Repository {
	return &Repository{db}
}

func (repository *Repository) GetThreadCountByForumID(id uint64) (uint64, error) {
	var num uint64
	if err := repository.db.QueryRow(
		"SELECT count(*) FROM threads WHERE forum = $1 ",
		id,
	).
		Scan(&num); err != nil {
		return 0, err
	}

	return num, nil
}

func (repository *Repository) InsertThreadInto(insertedThread *models.Thread) error {
	if err := repository.db.QueryRow(
		"INSERT INTO threads (slug, author, title, message, forum, created) VALUES (NULLIF ($1, ''), $2, $3, $4, $5, $6) RETURNING id",
		insertedThread.Slug,
		insertedThread.AuthorID,
		insertedThread.Title,
		insertedThread.About,
		insertedThread.ForumID,
		insertedThread.CreationDate,
	).
		Scan(&insertedThread.ID); err != nil {
		return err
	}

	return nil
}

func (repository *Repository) GetThreadBySlug(slug string) (*models.Thread, error) {
	threadBySlug := &models.Thread{}
	if err := repository.db.QueryRow(
		"SELECT t.id, u.nickname, t.created, t.forum, f.slug, t.message, coalesce (t.slug, ''), t.title, t.votes FROM threads AS t JOIN users AS u ON (t.author = u.id) JOIN forums AS f ON (f.id = t.forum) WHERE lower(t.slug) = lower($1)",
		slug,
	).
		Scan(
		&threadBySlug.ID,
		&threadBySlug.Author,
		&threadBySlug.CreationDate,
		&threadBySlug.ForumID,
		&threadBySlug.Forum,
		&threadBySlug.About,
		&threadBySlug.Slug,
		&threadBySlug.Title,
		&threadBySlug.Votes,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.ErrDoesntExists
		}
		return nil, err
	}

	return threadBySlug, nil
}

func (repository *Repository) GetThreadByID(id uint64) (*models.Thread, error) {
	threadById := &models.Thread{}
	if err := repository.db.QueryRow(
		"SELECT t.id, u.nickname, t.created, t.forum, f.slug, t.message, coalesce (t.slug, ''), t.title, t.votes FROM threads AS t JOIN users AS u ON (t.author = u.id) JOIN forums AS f ON (f.id = t.forum) WHERE t.id = $1",
		id,
	).
		Scan(
		&threadById.ID,
		&threadById.Author,
		&threadById.CreationDate,
		&threadById.ForumID,
		&threadById.Forum,
		&threadById.About,
		&threadById.Slug,
		&threadById.Title,
		&threadById.Votes,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.ErrDoesntExists
		}
		return nil, err
	}
	return threadById, nil
}

func (repository *Repository) GetThreadByForumSlug(slug string, offset uint64, time string, ranged bool) ([]*models.Thread, error) {
	threadByForumSlug := []*models.Thread{}

	query := "SELECT t.id, u.nickname, t.created, f.slug, t.message, coalesce (t.slug, ''), t.title, t.votes FROM threads AS t JOIN users AS u ON (t.author = u.id) JOIN forums AS f ON (f.id = t.forum) WHERE lower(f.slug) = lower($1)"

	order := " ORDER BY t.created"

	if ranged {
		order += " DESC"
	}

	if offset != 0 {
		order += fmt.Sprintf(" LIMIT %d", offset)
	}

	var rows *pgx.Rows
	var err error

	if time != "" {
		if ranged {
			query += " AND t.created <= $2"
		} else {
			query += " AND t.created >= $2"
		}
		rows, err = repository.db.Query(query + order, slug, time)
	} else {
		rows, err = repository.db.Query(query + order, slug)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		threadSlug := &models.Thread{}
		err = rows.Scan(
			&threadSlug.ID,
			&threadSlug.Author,
			&threadSlug.CreationDate,
			&threadSlug.Forum,
			&threadSlug.About,
			&threadSlug.Slug,
			&threadSlug.Title,
			&threadSlug.Votes,
		)
		if err != nil {
			return nil, err
		}

		threadByForumSlug = append(threadByForumSlug, threadSlug)
	}

	return threadByForumSlug, nil
}

func (repository *Repository) UpdateThread(updatingThread *models.Thread) error {
	if _, err := repository.db.Exec(
		"UPDATE threads SET message = $2, title = $3 WHERE id = $1",
		updatingThread.ID,
		updatingThread.About,
		updatingThread.Title,
	); err != nil {
		return err
	}
	return nil
}
