package domain

import (
	"time"
	"todo-app/pkg/client"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID     `json:"id"`
	Email      string        `json:"email"`
	Password   string        `json:"password"`
	FirstName  string        `json:"first_name"`
	LastName   string        `json:"last_name"`
	Phone      string        `json:"phone"`
	Role       int8          `json:"role"`
	Salt       string        `json:"salt"`
	Status     client.Status `json:"status"`
	Created_at time.Time     `json:"created_at"`
	Updated_at time.Time     `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// type UserUpdate struct {
// 	ID         uuid.UUID `json:"id"`
// 	Email      string    `json:"email"`
// 	Password   string    `json:"password"`
// 	Updated_at time.Time `json:"updated_at"`
// }

// func (UserUpdate) TableName() string {
// 	return User{}.TableName()
// }
