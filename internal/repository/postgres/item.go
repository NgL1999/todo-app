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

func (ir *itemRepo) GetAll(items *[]domain.Item) error {
	if err := ir.db.Find(items).Error; err != nil {
		return err
	}
	return nil
}

func (ir *itemRepo) GetById(item *domain.Item, id string) error {
	if err := ir.db.Find(item, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (ir *itemRepo) UpdateById(item *domain.Item) (int64, error) {
	result := ir.db.Model(item).Updates(item)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (ir *itemRepo) DeleteById(item *domain.Item) (int64, error) {
	result := ir.db.Delete(item)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
