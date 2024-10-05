package postgres

import (
	"todo-app/domain"

	"gorm.io/gorm"
)

type itemRepo struct {
	db *gorm.DB
}

func NewItemRepo(db *gorm.DB) *itemRepo {
	return &itemRepo{
		db: db,
	}
}

func (ir *itemRepo) Save(item *domain.Item) error {
	if err := ir.db.Create(item).Error; err != nil {
		return err
	}
	return nil
}