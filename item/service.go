package item

import (
	"context"
	"todo-app/domain"

	"github.com/google/uuid"
)

type IItemRepo interface {
	Save(item *domain.Item) error
	GetAll(items *[]domain.Item) error
	GetById(item *domain.Item, id string) error
	UpdateById(item *domain.ItemUpdate) (int64, error)
	DeleteById(item *domain.Item) (int64, error)
}

type itemService struct {
	itemRepo IItemRepo
}

func NewItemService(repo IItemRepo) *itemService {
	return &itemService{
		itemRepo: repo,
	}
}

func (is *itemService) CreateItem(ctx context.Context, item *domain.Item) error {
	item.ID = uuid.New()
	if err := is.itemRepo.Save(item); err != nil {
		return err
	}
	return nil
}

func (is *itemService) GetAllItems(ctx context.Context, items *[]domain.Item) error {
	if err := is.itemRepo.GetAll(items); err != nil {
		return err
	}
	return nil
}

func (is *itemService) GetItemById(ctx context.Context, item *domain.Item, id string) error {
	if err := is.itemRepo.GetById(item, id); err != nil {
		return err
	}
	return nil
}

func (is *itemService) UpdateItemById(ctx context.Context, item *domain.ItemUpdate) (int64, error) {
	result, err := is.itemRepo.UpdateById(item)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (is *itemService) DeleteItemById(ctx context.Context, item *domain.Item) (int64, error) {
	result, err := is.itemRepo.DeleteById(item)
	if err != nil {
		return 0, err
	}
	return result, nil
}
