package usecases

import (
	"db_project/app/models"
	"db_project/app/repositories"
	"db_project/utils/constants"
	"db_project/utils/errors"
	"github.com/jackc/pgx/v4"
	"strconv"
)

type PostUseCase struct {
	postRepository repositories.IPostRepository
	forumUseCase   IForumUseCase
	userUseCase    IUserUseCase
	threadUseCase  IThreadUseCase
}

type IPostUseCase interface {
	Get(id int, details []string) (postDetailed *models.ParamsPost, err error)
	Update(post *models.Post) (updatedPost *models.Post, err error)
}

func CreatePostUseCase(postRepository repositories.IPostRepository,
	forumUseCase IForumUseCase,
	userUseCase IUserUseCase,
	threadUseCase IThreadUseCase) IPostUseCase {
	return &PostUseCase{
		postRepository: postRepository,
		forumUseCase:   forumUseCase,
		userUseCase:    userUseCase,
		threadUseCase:  threadUseCase,
	}
}

func (usecase *PostUseCase) Get(id int, details []string) (postDetailed *models.ParamsPost, err error) {
	postDetailed = &models.ParamsPost{}
	postDetailed.Post, err = usecase.postRepository.Get(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.PostNotFound
			return
		}
		err = errors.ServerInternal
		return
	}

	for _, detailType := range details {
		switch detailType {
		case constants.PostUser:
			postDetailed.Author, err = usecase.userUseCase.Get(&postDetailed.Post.Author)
			if err != nil {
				postDetailed = nil
				return
			}
		case constants.PostThread:
			postDetailed.Thread, err = usecase.threadUseCase.Get(strconv.Itoa(postDetailed.Post.Thread))
			if err != nil {
				postDetailed = nil
				return
			}
		case constants.PostForum:
			postDetailed.Forum, err = usecase.forumUseCase.Get(postDetailed.Post.Forum)
			if err != nil {
				postDetailed = nil
				return
			}
		default:
			postDetailed = nil
			err = errors.BadRequest.SetTextDetails("неверные query параметры")
			return
		}
	}

	return
}

func (usecase *PostUseCase) Update(post *models.Post) (updatedPost *models.Post, err error) {
	updatedPost, err = usecase.postRepository.Update(post)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.PostNotFound
			return
		}
		err = errors.ServerInternal
		return
	}

	return
}
