package gin

import (
	"net/http"
	"todo-app/domain"
	"todo-app/pkg/client"
	"todo-app/pkg/tokenprovider"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	Register(data *domain.UserCreate) error
	Login(data *domain.UserLogin) (tokenprovider.Token, error)
}

type userHandler struct {
	userService IUserService
}

func NewUserHandler(apiVersion *gin.RouterGroup, svc IUserService) {
	userHandler := &userHandler{
		userService: svc,
	}

	users := apiVersion.Group("users")
	{
		users.POST("/register", userHandler.RegisterHandler)
		users.POST("/login", userHandler.LoginHandler)
	}
}

func (uh *userHandler) RegisterHandler(c *gin.Context) {
	var data domain.UserCreate

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	if err := uh.userService.Register(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(data.ID))
}

func (uh *userHandler) LoginHandler(c *gin.Context) {
	var data domain.UserLogin

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	token, err := uh.userService.Login(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(token))
}

// func (h *userHandler) GetAllUsersHandler(c *gin.Context) {
// 	users := []domain.User{}

// 	if err := h.userService.GetAllUsers(&users); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": users,
// 	})
// }

// func (h *userHandler) GetUserByIdHandler(c *gin.Context) {
// 	user := domain.User{}
// 	id := c.Param("id")

// 	if err := h.userService.GetUserById(&user, id); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": user,
// 	})
// }

// func (h *userHandler) UpdateUserByIdHandler(c *gin.Context) {
// 	user := domain.User{}

// 	if err := c.ShouldBind(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	id, err := uuid.Parse(c.Param("id"))

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	user.ID = id
// 	result, err := h.userService.UpdateUserById(&user)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"updatedCount": result,
// 	})
// }

// func (h *userHandler) DeleteUserByIdHandler(c *gin.Context) {
// 	user := domain.User{}
// 	id, err := uuid.Parse(c.Param("id"))

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	user.ID = id
// 	result, err := h.userService.DeleteUserById(&user)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"deletedCount": result,
// 	})
// }
