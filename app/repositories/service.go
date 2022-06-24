package repositories

import (
	"context"
	"db_project/app/models"
	"db_project/utils/constants"
	"github.com/jackc/pgx/v4/pgxpool"
)

type IServiceRepository interface {
	Clear() (err error)
	Status() (status *models.ForumStatus, err error)
}

type ServiceRepository struct {
	db *pgxpool.Pool
}

func CreateServiceRepository(db *pgxpool.Pool) IServiceRepository {
	return &ServiceRepository{db: db}
}

func (repo *ServiceRepository) Clear() (err error) {
	_, err = repo.db.Exec(context.Background(), constants.ServiceQuery["Clear"])
	return
}

func (repo *ServiceRepository) Status() (status *models.ForumStatus, err error) {
	ctx := context.Background()
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			trErr := tx.Commit(ctx)
			if trErr != nil {
				err = trErr
			}
		} else {
			trErr := tx.Rollback(ctx)
			if trErr != nil {
				err = trErr
			}
		}
	}()

	status = &models.ForumStatus{}
	if err = tx.QueryRow(ctx, constants.ServiceQuery["queryUsers"]).Scan(&status.User); err != nil {
		status = nil
		return
	}
	if err = tx.QueryRow(ctx, constants.ServiceQuery["queryForums"]).Scan(&status.Forum); err != nil {
		status = nil
		return
	}
	if err = tx.QueryRow(ctx, constants.ServiceQuery["queryThreads"]).Scan(&status.Thread); err != nil {
		status = nil
		return
	}
	if err = tx.QueryRow(ctx, constants.ServiceQuery["queryPosts"]).Scan(&status.Post); err != nil {
		status = nil
		return
	}

	return
}
