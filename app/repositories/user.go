package repositories

import (
	"context"
	"db_project/app/models"
	"db_project/utils/constants"
	"github.com/jackc/pgx/v4/pgxpool"
)

type IUserRepository interface {
	Get(nickname *string) (user *models.User, err error)
	Update(user *models.User) (updatedUser *models.User, err error)
	GetUsersByUserNicknameOrEmail(user *models.User) (users []*models.User, err error)
	All() (users *[]models.User, err error)
	Create(user *models.User) (err error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func CreateUserRepository(db *pgxpool.Pool) IUserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Get(nickname *string) (user *models.User, err error) {
	user = &models.User{}
	row := repo.db.QueryRow(context.Background(), constants.UserQuery["Get"], *nickname)
	err = row.Scan(
		&user.Username,
		&user.FullName,
		&user.About,
		&user.Email)
	return
}

func (repo *UserRepository) All() (users *[]models.User, err error) {
	return
}

func (repo *UserRepository) Create(user *models.User) (err error) {
	_, err = repo.db.Exec(context.Background(), constants.UserQuery["Create"], user.Username, user.FullName, user.About, user.Email)
	return
}

func (repo *UserRepository) Update(user *models.User) (updatedUser *models.User, err error) {
	row := repo.db.QueryRow(context.Background(), constants.UserQuery["Update"], user.FullName, user.About, user.Email, user.Username)
	updatedUser = &models.User{}
	err = row.Scan(
		&updatedUser.Username,
		&updatedUser.FullName,
		&updatedUser.About,
		&updatedUser.Email)
	return
}

func (repo *UserRepository) GetUsersByUserNicknameOrEmail(user *models.User) (users []*models.User, err error) {
	rows, err := repo.db.Query(context.Background(), constants.UserQuery["GetUsersByUserNOE"], user.Username, user.Email)
	defer rows.Close()

	if err != nil {
		return
	}

	users = make([]*models.User, 0)
	for rows.Next() {
		conflictUser := &models.User{}
		err = rows.Scan(
			&conflictUser.Username,
			&conflictUser.FullName,
			&conflictUser.About,
			&conflictUser.Email)
		if err != nil {
			users = nil
			return
		}
		users = append(users, conflictUser)
	}

	return
}
