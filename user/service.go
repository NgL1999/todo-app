package user

import (
	"errors"
	"todo-app/domain"
	"todo-app/pkg/client"
	"todo-app/pkg/tokenprovider"
	"todo-app/pkg/util"

	"github.com/google/uuid"
)

type IUserRepo interface {
	Save(user *domain.UserCreate) error
	GetUser(conditions map[string]any) (*domain.User, error)
}

type IHasher interface {
	Hash(data string) string
}

type userService struct {
	userRepo      IUserRepo
	hasher        IHasher
	tokenProvider tokenprovider.Provider
	expiry        int
}

func NewUserService(repo IUserRepo, hasher IHasher, tokenProvider tokenprovider.Provider, expiry int) *userService {
	return &userService{
		userRepo:      repo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
		expiry:        expiry,
	}
}

func (us *userService) Register(data *domain.UserCreate) error {
	if err := data.Validate(); err != nil {
		return client.ErrInvalidRequest(err)
	}

	user, err := us.userRepo.GetUser(map[string]any{"email": data.Email})
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
	data.Password = us.hasher.Hash(data.Password + salt)
	data.Salt = salt
	data.Role = 1

	if err := us.userRepo.Save(data); err != nil {
		return client.ErrCannotCreateEntity(data.TableName(), err)
	}

	return nil
}

func (us *userService) Login(data *domain.UserLogin) (tokenprovider.Token, error) {
	user, err := us.userRepo.GetUser(map[string]interface{}{"email": data.Email})
	if err != nil {
		return nil, domain.ErrEmailOrPasswordInvalid
	}

	passHashed := us.hasher.Hash(data.Password + user.Salt)

	if user.Password != passHashed {
		return nil, domain.ErrEmailOrPasswordInvalid
	}

	payload := &client.TokenPayload{
		UID:   user.ID,
		URole: user.Role.String(),
	}

	accessToken, err := us.tokenProvider.Generate(payload, us.expiry)
	if err != nil {
		return nil, client.ErrInternal(err)
	}

	return accessToken, nil
}

// func (s *userService) GetAllUsers(users *[]domain.User) error {
// 	if err := s.userRepo.GetAllUsers(users); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (is *userService) GetUserById(user *domain.User, id string) error {
// 	if err := is.userRepo.GetUserById(user, id); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (is *userService) UpdateUserById(user *domain.User) (int64, error) {
// 	result, err := is.userRepo.UpdateUserById(user)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return result, nil
// }

// func (is *userService) DeleteUserById(user *domain.User) (int64, error) {
// 	result, err := is.userRepo.DeleteUserById(user)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return result, nil
// }