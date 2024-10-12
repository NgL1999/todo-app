package domain

import (
	"time"
	"todo-app/pkg/client"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      client.Status `json:"status"`
	Created_at  time.Time      `json:"created_at"`
	Updated_at  time.Time      `json:"updated_at"`
}

func (Item) TableName() string {
	return "items"
}

type ItemUpdate struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Updated_at  time.Time `json:"updated_at"`
}

func (ItemUpdate) TableName() string {
	return Item{}.TableName()
}
