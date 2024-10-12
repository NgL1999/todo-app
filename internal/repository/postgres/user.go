package postgres

import (
	"todo-app/domain"

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

func (ir *userRepo) Save(user *domain.User) error {
	if err := ir.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ir *userRepo) GetAll(users *[]domain.User) error {
	if err := ir.db.Find(users).Error; err != nil {
		return err
	}
	return nil
}

func (ir *userRepo) GetById(user *domain.User, id string) error {
	if err := ir.db.Find(user, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (ir *userRepo) UpdateById(user *domain.User) (int64, error) {
	result := ir.db.Model(user).Updates(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (ir *userRepo) DeleteById(user *domain.User) (int64, error) {
	result := ir.db.Delete(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
