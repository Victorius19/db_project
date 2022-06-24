package repositories

import (
	"context"
	"db_project/app/models"
	"db_project/utils/constants"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type ForumRepository struct {
	db *pgxpool.Pool
}

func CreateForumRepository(db *pgxpool.Pool) IForumRepository {
	return &ForumRepository{db: db}
}

type IForumRepository interface {
	CreateThread(thread *models.Thread) (createdThread *models.Thread, err error)
	GetThreads(slug string, params *models.ForumQueryParams) (threads []*models.Thread, err error)
	GetUsers(slug string, params *models.ForumUserQueryParams) (users []*models.User, err error)
	Create(forum *models.Forum) (createdForum *models.Forum, err error)
	Get(slug string) (forum *models.Forum, err error)
}

func (repo *ForumRepository) Create(forum *models.Forum) (createdForum *models.Forum, err error) {
	row := repo.db.QueryRow(context.Background(), constants.ForumQuery["Create"], forum.Slug, forum.Title, forum.User)

	createdForum = &models.Forum{}
	err = row.Scan(
		&createdForum.Slug,
		&createdForum.Title,
		&createdForum.User,
		&createdForum.Posts,
		&createdForum.Threads)
	return
}

func (repo *ForumRepository) Get(slug string) (forum *models.Forum, err error) {
	row := repo.db.QueryRow(context.Background(), constants.ForumQuery["Get"], slug)

	forum = &models.Forum{}
	err = row.Scan(
		&forum.ID,
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads)
	return
}

func (repo *ForumRepository) CreateThread(thread *models.Thread) (createdThread *models.Thread, err error) {

	row := repo.db.QueryRow(context.Background(), constants.ForumQuery["CreateThread"], thread.Slug, thread.Author, thread.Forum, thread.Title, thread.Msg, thread.Created)

	createdThread = &models.Thread{}
	err = row.Scan(
		&createdThread.ID,
		&createdThread.Slug,
		&createdThread.Author,
		&createdThread.Forum,
		&createdThread.Title,
		&createdThread.Msg,
		&createdThread.Created,
		&createdThread.Votes)
	return
}

func (repo *ForumRepository) GetThreads(slug string, params *models.ForumQueryParams) (threads []*models.Thread, err error) {
	query := constants.ForumQuery["GetThreads"]
	var rows pgx.Rows
	if !params.Since.Equal(time.Time{}) {
		if params.Desc {
			query += constants.ForumQuery["GetThreadsDesc"]
		} else {
			query += constants.ForumQuery["GetThreadsNoDesc"]
		}
		rows, err = repo.db.Query(context.Background(), query, slug, params.Since, params.Limit)
	} else {
		if params.Desc {
			query += constants.ForumQuery["GetThreadsSinceDesc"]
		} else {
			query += constants.ForumQuery["GetThreadsSinceNoDesc"]
		}
		rows, err = repo.db.Query(context.Background(), query, slug, params.Limit)
	}

	defer rows.Close()
	if err != nil {
		return
	}

	threads = make([]*models.Thread, 0)
	for rows.Next() {
		thread := &models.Thread{}
		err = rows.Scan(
			&thread.ID,
			&thread.Slug,
			&thread.Author,
			&thread.Forum,
			&thread.Title,
			&thread.Msg,
			&thread.Created,
			&thread.Votes)
		if err != nil {
			threads = nil
			return
		}
		threads = append(threads, thread)
	}

	return
}

func (repo *ForumRepository) GetUsers(slug string, params *models.ForumUserQueryParams) (users []*models.User, err error) {
	query := constants.ForumQuery["GetUsers"]

	var rows pgx.Rows
	if params.Since != "" {
		if params.Desc {
			query += constants.ForumQuery["GetUsersSinceDesc"]
		} else {
			query += constants.ForumQuery["GetUsersSinceNoDesc"]
		}
		rows, err = repo.db.Query(context.Background(), query, slug, params.Since, params.Limit)
	} else {
		if params.Desc {
			query += constants.ForumQuery["GetUsersDesc"]
		} else {
			query += constants.ForumQuery["GetUsersNoDesc"]
		}
		rows, err = repo.db.Query(context.Background(), query, slug, params.Limit)
	}

	defer rows.Close()
	if err != nil {
		return
	}

	users = make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(
			&user.Username,
			&user.FullName,
			&user.About,
			&user.Email)
		if err != nil {
			users = nil
			return
		}
		users = append(users, user)
	}

	return
}
