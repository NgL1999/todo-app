package user

import (
	"errors"
	"todo-app/domain"
	"todo-app/pkg/client"
	"todo-app/pkg/tokenprovider"
	"todo-app/pkg/util"

	"github.com/google/uuid"
)

type UserRepo interface {
	Save(user *domain.UserCreate) error
	GetUser(conditions map[string]any) (*domain.User, error)
	GetAllUsers(users *[]domain.User) error
	GetUserById(user *domain.User, id string) error
	UpdateUserById(user *domain.User) (int64, error)
	DeleteUserById(user *domain.User) (int64, error)
}

type Hasher interface {
	Hash(data string) string
}

type userService struct {
	userRepo      UserRepo
	hasher        Hasher
	tokenProvider tokenprovider.Provider
	expiry        int
}

func NewUserService(repo UserRepo, hasher Hasher, tokenProvider tokenprovider.Provider, expiry int) *userService {
	return &userService{
		userRepo:      repo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
		expiry:        expiry,
	}
}

func (s *userService) Register(data *domain.UserCreate) error {
	if err := data.Validate(); err != nil {
		return client.ErrInvalidRequest(err)
	}

	user, err := s.userRepo.GetUser(map[string]any{"email": data.Email})
	if err != nil {
		if !errors.Is(err, client.ErrRecordNotFound) {
			return err
		}
	}

	if user != nil {
		return domain.ErrEmailExisted
	}

	salt := util.GenSalt(50)

	data.ID = uuid.New()
	data.Password = s.hasher.Hash(data.Password + salt)
	data.Salt = salt
	data.Role = 1

	if err := s.userRepo.Save(data); err != nil {
		return client.ErrCannotCreateEntity(data.TableName(), err)
	}

	return nil
}

func (s *userService) Login(data *domain.UserLogin) (tokenprovider.Token, error) {
	user, err := s.userRepo.GetUser(map[string]interface{}{"email": data.Email})
	if err != nil {
		return nil, domain.ErrEmailOrPasswordInvalid
	}

	passHashed := s.hasher.Hash(data.Password + user.Salt)

	if user.Password != passHashed {
		return nil, domain.ErrEmailOrPasswordInvalid
	}

	payload := &client.TokenPayload{
		UID:   user.ID,
		URole: user.Role.String(),
	}

	accessToken, err := s.tokenProvider.Generate(payload, s.expiry)
	if err != nil {
		return nil, client.ErrInternal(err)
	}

	return accessToken, nil
}

func (s *userService) GetAllUsers(users *[]domain.User) error {
	if err := s.userRepo.GetAllUsers(users); err != nil {
		return err
	}
	return nil
}

func (is *userService) GetUserById(user *domain.User, id string) error {
	if err := is.userRepo.GetUserById(user, id); err != nil {
		return err
	}
	return nil
}

func (is *userService) UpdateUserById(user *domain.User) (int64, error) {
	result, err := is.userRepo.UpdateUserById(user)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (is *userService) DeleteUserById(user *domain.User) (int64, error) {
	result, err := is.userRepo.DeleteUserById(user)
	if err != nil {
		return 0, err
	}
	return result, nil
}
