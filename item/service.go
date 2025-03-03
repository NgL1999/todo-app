package item

import (
	"time"
	"todo-app/domain"
	"todo-app/pkg/client"

	"github.com/google/uuid"
)

//go:generate mockery --name IItemRepo
type IItemRepo interface {
	Save(item *domain.ItemCreation) error
	GetAll(filter map[string]any, paging *client.Paging) ([]domain.Item, error)
	Get(filter map[string]any) (domain.Item, error)
	Update(filter map[string]any, item *domain.ItemUpdate) error
	Delete(filter map[string]any) error
}

type itemService struct {
	itemRepo IItemRepo
}

func NewItemService(repo IItemRepo) *itemService {
	return &itemService{
		itemRepo: repo,
	}
}

func (is *itemService) Create(item *domain.ItemCreation) error {
	if err := item.Validate(); err != nil {
		return client.ErrInvalidRequest(err)
	}

	item.ID = uuid.New()
	if err := is.itemRepo.Save(item); err != nil {
		return client.ErrCannotCreateEntity(item.TableName(), err)
	}

	return nil
}

func (is *itemService) GetAll(userID uuid.UUID, paging *client.Paging) ([]domain.Item, error) {
	filter := map[string]any{"user_id": userID}
	items, err := is.itemRepo.GetAll(filter, paging)
	if err != nil {
		return nil, client.ErrCannotListEntity(domain.Item{}.TableName(), err)
	}

	return items, nil
}

func (is *itemService) GetById(id, userID uuid.UUID) (domain.Item, error) {
	item, err := is.itemRepo.Get(map[string]any{"id": id, "user_id": userID})
	if err != nil {
		return domain.Item{}, client.ErrCannotGetEntity(item.TableName(), err)
	}

	return item, nil
}

func (is *itemService) UpdateById(id, userID uuid.UUID, item *domain.ItemUpdate) error {
	item.UpdatedAt = time.Now()
	err := is.itemRepo.Update(map[string]any{"id": id, "user_id": userID}, item)
	if err != nil {
		return client.ErrCannotUpdateEntity(item.TableName(), err)
	}

	return nil
}

func (is *itemService) DeleteById(id, userID uuid.UUID) error {
	err := is.itemRepo.Delete(map[string]any{"id": id, "user_id": userID})
	if err != nil {
		return client.ErrCannotDeleteEntity(domain.Item{}.TableName(), err)
	}

	return nil
}
