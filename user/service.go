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
	Get(filter map[string]any) (*domain.User, error)
	GetAll(filter map[string]any, paging *client.Paging) ([]domain.User, error)
	Update(filter map[string]any, user *domain.UserUpdate) error
	Delete(filter map[string]any) error
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

	user, err := us.userRepo.Get(map[string]any{"email": data.Email})
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
	user, err := us.userRepo.Get(map[string]interface{}{"email": data.Email})
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

func (us *userService) GetAll(paging *client.Paging) ([]domain.User, error) {
	users, err := us.userRepo.GetAll(nil, paging)
	if err != nil {
		return nil, client.ErrCannotListEntity(domain.User{}.TableName(), err)
	}

	return users, nil
}

func (us *userService) GetById(id uuid.UUID) (*domain.User, error) {
	user, err := us.userRepo.Get(map[string]any{"id": id})
	if err != nil {
		return nil, client.ErrCannotGetEntity(user.TableName(), err)
	}

	return user, nil
}

func (us *userService) UpdateById(id uuid.UUID, user *domain.UserUpdate) error {
	// user.UpdatedAt = time.Now()
	err := us.userRepo.Update(map[string]any{"id": id}, user)
	if err != nil {
		return client.ErrCannotUpdateEntity(user.TableName(), err)
	}

	return nil
}

func (us *userService) DeleteById(id uuid.UUID) error {
	err := us.userRepo.Delete(map[string]any{"id": id})
	if err != nil {
		return client.ErrCannotDeleteEntity(domain.User{}.TableName(), err)
	}

	return nil
}
