package repositories

import (
	"context"
	"db_project/app/models"
	"db_project/utils/constants"
	"github.com/jackc/pgx/v4/pgxpool"
)

type IPostRepository interface {
	Get(id int) (post *models.Post, err error)
	Update(post *models.Post) (updatedPost *models.Post, err error)
}

type PostRepository struct {
	db *pgxpool.Pool
}

func CreatePostRepository(db *pgxpool.Pool) IPostRepository {
	return &PostRepository{db: db}
}

func (repo *PostRepository) Get(id int) (post *models.Post, err error) {
	row := repo.db.QueryRow(context.Background(), constants.PostQuery["Get"], id)
	post = &models.Post{}
	err = row.Scan(
		&post.ID,
		&post.Parent,
		&post.Author,
		&post.Forum,
		&post.Thread,
		&post.Created,
		&post.IsEdited,
		&post.Message)
	return
}

func (repo *PostRepository) Update(post *models.Post) (updatedPost *models.Post, err error) {
	row := repo.db.QueryRow(context.Background(), constants.PostQuery["Update"], post.Message, post.ID)

	updatedPost = &models.Post{}
	err = row.Scan(
		&updatedPost.ID,
		&updatedPost.Parent,
		&updatedPost.Author,
		&updatedPost.Forum,
		&updatedPost.Thread,
		&updatedPost.Created,
		&updatedPost.IsEdited,
		&updatedPost.Message)
	return
}
