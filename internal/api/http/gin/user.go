package gin

import (
	"context"
	"net/http"
	"time"
	"todo-app/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IUserService interface {
	Register(ctx context.Context, user *domain.User) error
	GetAllUsers(ctx context.Context, users *[]domain.User) error
	GetUserById(ctx context.Context, user *domain.User, id string) error
	UpdateUserById(ctx context.Context, user *domain.User) (int64, error)
	DeleteUserById(ctx context.Context, user *domain.User) (int64, error)
}

type userHandler struct {
	userService IUserService
}

// Constructor
func NewUserHandler(apiVersion *gin.RouterGroup, isvc IUserService) {
	userHandler := &userHandler{
		userService: isvc,
	}

	users := apiVersion.Group("users")
	{
		users.POST("/", userHandler.RegisterHandler)
		users.GET("/all", userHandler.GetAllUsersHandler)
		users.GET("/:id", userHandler.GetUserByIdHandler)
		users.PATCH("/:id", userHandler.UpdateUserByIdHandler)
		users.DELETE("/:id", userHandler.DeleteUserByIdHandler)
	}
}

func (ih *userHandler) RegisterHandler(c *gin.Context) {
	user := domain.User{}

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.ID = uuid.New()
	user.Created_at = time.Now()
	user.Updated_at = time.Now()

	if err := ih.userService.Register(c, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": user.ID,
	})
}

func (ih *userHandler) GetAllUsersHandler(c *gin.Context) {
	users := []domain.User{}

	if err := ih.userService.GetAllUsers(c, &users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func (ih *userHandler) GetUserByIdHandler(c *gin.Context) {
	user := domain.User{}
	id := c.Param("id")

	if err := ih.userService.GetUserById(c, &user, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (ih *userHandler) UpdateUserByIdHandler(c *gin.Context) {
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
	user.Updated_at = time.Now()
	result, err := ih.userService.UpdateUserById(c, &user)

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

func (ih *userHandler) DeleteUserByIdHandler(c *gin.Context) {
	user := domain.User{}
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.ID = id
	result, err := ih.userService.DeleteUserById(c, &user)

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
