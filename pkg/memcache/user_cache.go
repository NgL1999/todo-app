package memcache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"todo-app/domain"

	"github.com/google/uuid"
)

type IRealStore interface {
	Get(filter map[string]any) (*domain.User, error)
}

type userCaching struct {
	store     ICache
	realStore IRealStore
	once      *sync.Once
}

func NewUserCaching(store ICache, realStore IRealStore) *userCaching {
	return &userCaching{
		store:     store,
		realStore: realStore,
		once:      new(sync.Once),
	}
}

func (uc *userCaching) Get(conditions map[string]interface{}) (*domain.User, error) {
	var ctx = context.Background()
	var user domain.User

	userId := conditions["id"].(uuid.UUID)
	key := fmt.Sprintf("user-%d", userId)

	err := uc.store.Get(ctx, key, &user)

	if err == nil && user.ID != uuid.Nil {
		return &user, nil
	}

	var userErr error

	uc.once.Do(func() {
		realUser, err := uc.realStore.Get(conditions)

		if err != nil {
			userErr = err
			log.Println(err)
			return
		}

		// Update cache
		user = *realUser
		_ = uc.store.Set(ctx, key, realUser, time.Hour*2)
	})

	if userErr != nil {
		return nil, userErr
	}

	err = uc.store.Get(ctx, key, &user)

	if err == nil && user.ID != uuid.Nil {
		return &user, nil
	}

	return nil, err
}
