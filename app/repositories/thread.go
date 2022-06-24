package repositories

import (
	"context"
	"db_project/app/models"
	"db_project/utils/constants"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
	"time"
)

type IThreadRepository interface {
	GetBySlug(slug string) (thread *models.Thread, err error)
	UpdateByID(thread *models.Thread) (updatedThread *models.Thread, err error)
	GetByID(id int) (thread *models.Thread, err error)
	CreatePosts(threadId int, forumSlug string, post []*models.Post) (createdPosts []*models.Post, err error)
	GetPosts(threadId int, params *models.PostsQueryParams) (posts []*models.Post, err error)
	VoteBySlug(slug string, vote *models.Vote) (err error)
	UpdateBySlug(thread *models.Thread) (updatedThread *models.Thread, err error)
	VoteByID(threadId int, vote *models.Vote) (err error)
}

type ThreadRepository struct {
	db *pgxpool.Pool
}

func CreateThreadRepository(db *pgxpool.Pool) IThreadRepository {
	return &ThreadRepository{db: db}
}

func (repo *ThreadRepository) GetBySlug(slug string) (thread *models.Thread, err error) {
	row := repo.db.QueryRow(context.Background(), constants.ThreadQuery["GetBySlug"], slug)

	thread = &models.Thread{}
	err = row.Scan(&thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Title, &thread.Msg, &thread.Created, &thread.Votes)
	return
}
func (repo *ThreadRepository) GetByID(id int) (thread *models.Thread, err error) {
	row := repo.db.QueryRow(context.Background(), constants.ThreadQuery["GetByID"], id)

	thread = &models.Thread{}
	err = row.Scan(&thread.ID, &thread.Slug, &thread.Author, &thread.Forum,
		&thread.Title, &thread.Msg, &thread.Created, &thread.Votes)
	return
}
func (repo *ThreadRepository) UpdateBySlug(thread *models.Thread) (updatedThread *models.Thread, err error) {
	row := repo.db.QueryRow(context.Background(), constants.ThreadQuery["UpdateBySlug"], thread.Title, thread.Msg, thread.Slug)
	updatedThread = &models.Thread{}
	err = row.Scan(&updatedThread.ID, &updatedThread.Slug, &updatedThread.Author, &updatedThread.Forum,
		&updatedThread.Title, &updatedThread.Msg, &updatedThread.Created, &updatedThread.Votes)
	return
}
func (repo *ThreadRepository) UpdateByID(thread *models.Thread) (updatedThread *models.Thread, err error) {
	row := repo.db.QueryRow(context.Background(), constants.ThreadQuery["UpdateByID"], thread.Title, thread.Msg, thread.ID)
	updatedThread = &models.Thread{}
	err = row.Scan(&updatedThread.ID, &updatedThread.Slug, &updatedThread.Author, &updatedThread.Forum,
		&updatedThread.Title, &updatedThread.Msg, &updatedThread.Created, &updatedThread.Votes)
	return
}

func (repo *ThreadRepository) VoteBySlug(slug string, vote *models.Vote) (err error) {
	_, err = repo.db.Exec(context.Background(), constants.ThreadQuery["VoteBySlug"], vote.Username, slug, vote.Voice)
	return
}

func (repo *ThreadRepository) VoteByID(id int, vote *models.Vote) (err error) {
	_, err = repo.db.Exec(context.Background(), constants.ThreadQuery["VoteByID"], vote.Username, id, vote.Voice)
	return
}

func (repo *ThreadRepository) CreatePostsBatch(threadId int, forumSlug string, posts []*models.Post) (createdPosts []*models.Post, err error) {
	ctx := context.Background()
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			trErr := tx.Rollback(ctx)
			if trErr != nil {
				err = trErr
			}
		} else {
			trErr := tx.Commit(ctx)
			if trErr != nil {
				err = trErr
			}
		}
	}()

	batch := new(pgx.Batch)
	createdTime := time.Now()

	for _, post := range posts {
		batch.Queue(constants.ThreadQuery["CreatePostsBatch"], post.Parent, post.Author, forumSlug, threadId, post.Message, createdTime)
	}

	batchRes := tx.SendBatch(ctx, batch)
	defer func() {
		batchErr := batchRes.Close()
		if batchErr != nil {
			err = batchErr
		}
	}()

	createdPosts = make([]*models.Post, 0)

	for i := 0; i < batch.Len(); i++ {
		createdPost := &models.Post{}

		row := batchRes.QueryRow()
		err = row.Scan(
			&createdPost.ID,
			&createdPost.Parent,
			&createdPost.Author,
			&createdPost.Forum,
			&createdPost.Thread,
			&createdPost.Created,
			&createdPost.IsEdited,
			&createdPost.Message)

		if err != nil {
			createdPosts = nil
			return
		}

		createdPosts = append(createdPosts, createdPost)
	}

	return
}

func (repo *ThreadRepository) CreatePosts(threadId int, forumSlug string, posts []*models.Post) (createdPosts []*models.Post, err error) {
	query := constants.ThreadQuery["PostsCreate"]

	createdTime := time.Now()
	var values []interface{}
	for i, _ := range posts {
		indexShift := 6 * i
		query += fmt.Sprintf("(NULLIF($%v, 0), $%v, $%v, $%v, $%v, $%v),",
			indexShift+1,
			indexShift+2,
			indexShift+3,
			indexShift+4,
			indexShift+5,
			indexShift+6)
		values = append(values, posts[i].Parent, posts[i].Author, forumSlug, threadId, posts[i].Message, createdTime)
	}
	query = strings.TrimSuffix(query, ",")
	query += constants.ThreadQuery["CreatePostsTwo"]

	rows, err := repo.db.Query(context.Background(), query, values...)
	defer rows.Close()

	if err != nil {
		return
	}

	createdPosts = make([]*models.Post, 0)
	for rows.Next() {
		createdPost := &models.Post{}
		err = rows.Scan(
			&createdPost.ID,
			&createdPost.Parent,
			&createdPost.Author,
			&createdPost.Forum,
			&createdPost.Thread,
			&createdPost.Created,
			&createdPost.IsEdited,
			&createdPost.Message)

		if err != nil {
			createdPosts = nil
			return
		}

		createdPosts = append(createdPosts, createdPost)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}

func (repo *ThreadRepository) GetPosts(threadId int, params *models.PostsQueryParams) (posts []*models.Post, err error) {
	var rows pgx.Rows

	if params.Since == 0 {
		if params.Desc {
			rows, err = repo.db.Query(context.Background(), constants.DescNoSincePostQuery[params.SortType],
				threadId, params.Limit)
		} else {
			rows, err = repo.db.Query(context.Background(), constants.AscNoSincePostQuery[params.SortType],
				threadId, params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = repo.db.Query(context.Background(), constants.DescSincePostQuery[params.SortType],
				threadId, params.Since, params.Limit)
		} else {
			rows, err = repo.db.Query(context.Background(), constants.AscSincePostQuery[params.SortType],
				threadId, params.Since, params.Limit)
		}
	}

	defer rows.Close()
	if err != nil {
		return
	}

	posts = make([]*models.Post, 0)
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(
			&post.ID,
			&post.Parent,
			&post.Author,
			&post.Forum,
			&post.Thread,
			&post.Created,
			&post.IsEdited,
			&post.Message)
		if err != nil {
			posts = nil
			return
		}
		posts = append(posts, post)
	}
	return
}
