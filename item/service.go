package item

import (
	"context"
	"todo-app/domain"

	"github.com/google/uuid"
)

type ItemRepo interface {
	Save(item *domain.Item) error
}

type itemService struct {
	itemRepo ItemRepo
}

func NewItemService(repo ItemRepo) *itemService {
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
