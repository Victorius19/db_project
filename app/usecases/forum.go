package usecases

import (
	"db_project/app/models"
	"db_project/app/repositories"
	"db_project/utils/errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type IForumUseCase interface {
	Create(forum *models.Forum) (createdForum *models.Forum, err error)
	Get(slug string) (forum *models.Forum, err error)
	CreateThread(thread *models.Thread) (createdThread *models.Thread, err error)
	GetThreads(slug string, params *models.ForumQueryParams) (threads []*models.Thread, err error)
	GetUsers(slug string, params *models.ForumUserQueryParams) (users []*models.User, err error)
}

type ForumUseCase struct {
	forumRepository  repositories.IForumRepository
	threadRepository repositories.IThreadRepository
}

func CreateForumUseCase(forumRepository repositories.IForumRepository,
	threadRepository repositories.IThreadRepository) IForumUseCase {
	return &ForumUseCase{forumRepository: forumRepository, threadRepository: threadRepository}
}

func (usecase *ForumUseCase) Create(forum *models.Forum) (createdForum *models.Forum, err error) {
	createdForum, err = usecase.forumRepository.Create(forum)

	if err != nil {
		pgconErr, ok := err.(*pgconn.PgError)
		if ok {
			switch pgconErr.SQLState() {
			case errors.Err23502:
				err = errors.NotFoundForumUser
				createdForum = nil
				return

			case errors.Err23505:
				createdForum, err = usecase.forumRepository.Get(forum.Slug)
				if err != nil {
					err = errors.ServerInternal
					return
				}
				err = errors.ForumAlreadyExists
				return

			default:
				err = errors.ServerInternal
			}
		}
	}

	return
}

func (usecase *ForumUseCase) Get(slug string) (forum *models.Forum, err error) {
	forum, err = usecase.forumRepository.Get(slug)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.NotFoundForum
		} else {
			err = errors.ServerInternal
		}
	}

	return
}

func (usecase *ForumUseCase) CreateThread(thread *models.Thread) (createdThread *models.Thread, err error) {
	createdThread, err = usecase.forumRepository.CreateThread(thread)

	if err != nil {
		pgconErr, ok := err.(*pgconn.PgError)
		if ok {
			switch pgconErr.SQLState() {
			case errors.Err23502:
				err = errors.ThreadUserOrForumNotFound
				createdThread = nil
				return

			case errors.Err23505:
				createdThread, err = usecase.threadRepository.GetBySlug(thread.Slug)
				if err != nil {
					err = errors.ServerInternal
					return
				}
				err = errors.ThreadAlreadyExists
				return

			default:
				err = errors.ServerInternal
			}
		}
	}

	return
}

func (usecase *ForumUseCase) GetThreads(slug string, params *models.ForumQueryParams) (threads []*models.Thread, err error) {
	threads, err = usecase.forumRepository.GetThreads(slug, params)
	if err != nil {
		err = errors.ServerInternal
		return
	}

	if len(threads) == 0 {
		_, err = usecase.Get(slug)
		if err != nil {
			return
		}
	}

	return
}

func (usecase *ForumUseCase) GetUsers(slug string, params *models.ForumUserQueryParams) (users []*models.User, err error) {
	users, err = usecase.forumRepository.GetUsers(slug, params)

	if err != nil {
		err = errors.ServerInternal
		return
	}

	if len(users) == 0 {
		if _, err = usecase.Get(slug); err != nil {
			return
		}
		return
	}

	return
}
