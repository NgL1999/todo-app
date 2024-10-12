package postgres

import (
	"errors"
	"todo-app/domain"
	"todo-app/pkg/client"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Save(user *domain.UserCreate) error {
	if err := r.db.Create(&user).Error; err != nil {
		return client.ErrDB(err)
	}

	return nil
}

func (r *userRepo) GetUser(conditions map[string]any) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where(conditions).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, client.ErrRecordNotFound
		}

		return nil, client.ErrDB(err)
	}

	return &user, nil
}

func (r *userRepo) GetAllUsers(users *[]domain.User) error {
	if err := r.db.Find(users).Error; err != nil {
		return err
	}
	return nil
}

func (ir *userRepo) GetUserById(user *domain.User, id string) error {
	if err := ir.db.Find(user, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (ir *userRepo) UpdateUserById(user *domain.User) (int64, error) {
	result := ir.db.Model(user).Updates(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (ir *userRepo) DeleteUserById(user *domain.User) (int64, error) {
	result := ir.db.Delete(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
