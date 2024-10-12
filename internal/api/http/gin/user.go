package gin

import (
	"net/http"
	"todo-app/domain"
	"todo-app/pkg/client"
	"todo-app/pkg/tokenprovider"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	Register(data *domain.UserCreate) error
	Login(data *domain.UserLogin) (tokenprovider.Token, error)
	GetAllUsers(users *[]domain.User) error
	GetUserById(user *domain.User, id string) error
	UpdateUserById(user *domain.User) (int64, error)
	DeleteUserById(user *domain.User) (int64, error)
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(apiVersion *gin.RouterGroup, svc UserService) {
	userHandler := &userHandler{
		userService: svc,
	}

	users := apiVersion.Group("/users")
	users.POST("/register", userHandler.RegisterUserHandler)
	users.POST("/login", userHandler.LoginHandler)
	users.GET("/", userHandler.GetAllUsersHandler)
	users.GET("/:id", userHandler.GetUserByIdHandler)
	users.PATCH("/:id", userHandler.UpdateUserByIdHandler)
	users.DELETE("/:id", userHandler.DeleteUserByIdHandler)
}

func (h *userHandler) RegisterUserHandler(c *gin.Context) {
	var data domain.UserCreate

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))

		return
	}

	if err := h.userService.Register(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(data.ID))
}

func (h *userHandler) LoginHandler(c *gin.Context) {
	var data domain.UserLogin

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))

		return
	}

	token, err := h.userService.Login(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(token))
}

func (h *userHandler) GetAllUsersHandler(c *gin.Context) {
	users := []domain.User{}

	if err := h.userService.GetAllUsers(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func (h *userHandler) GetUserByIdHandler(c *gin.Context) {
	user := domain.User{}
	id := c.Param("id")

	if err := h.userService.GetUserById(&user, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (h *userHandler) UpdateUserByIdHandler(c *gin.Context) {
	user := domain.User{}

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.ID = id
	result, err := h.userService.UpdateUserById(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"updatedCount": result,
	})
}

func (h *userHandler) DeleteUserByIdHandler(c *gin.Context) {
	user := domain.User{}
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.ID = id
	result, err := h.userService.DeleteUserById(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deletedCount": result,
	})
}
