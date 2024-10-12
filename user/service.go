package user

import (
	"context"
	"todo-app/domain"

	"github.com/google/uuid"
)

type IUserRepo interface {
	Save(user *domain.User) error
	GetAll(users *[]domain.User) error
	GetById(user *domain.User, id string) error
	UpdateById(user *domain.User) (int64, error)
	DeleteById(user *domain.User) (int64, error)
}

type userService struct {
	userRepo IUserRepo
}

func NewUserService(repo IUserRepo) *userService {
	return &userService{
		userRepo: repo,
	}
}

func (is *userService) Register(ctx context.Context, user *domain.User) error {
	user.ID = uuid.New()
	if err := is.userRepo.Save(user); err != nil {
		return err
	}
	return nil
}

func (is *userService) GetAllUsers(ctx context.Context, users *[]domain.User) error {
	if err := is.userRepo.GetAll(users); err != nil {
		return err
	}
	return nil
}

func (is *userService) GetUserById(ctx context.Context, user *domain.User, id string) error {
	if err := is.userRepo.GetById(user, id); err != nil {
		return err
	}
	return nil
}

func (is *userService) UpdateUserById(ctx context.Context, user *domain.User) (int64, error) {
	result, err := is.userRepo.UpdateById(user)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (is *userService) DeleteUserById(ctx context.Context, user *domain.User) (int64, error) {
	result, err := is.userRepo.DeleteById(user)
	if err != nil {
		return 0, err
	}
	return result, nil
}
