package gin

import (
	"net/http"
	"todo-app/domain"
	"todo-app/pkg/client"
	"todo-app/pkg/tokenprovider"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IUserService interface {
	Register(data *domain.UserCreate) error
	Login(data *domain.UserLogin) (tokenprovider.Token, error)
	GetAll(paging *client.Paging) ([]domain.User, error)
	GetById(id uuid.UUID) (*domain.User, error)
	UpdateById(id uuid.UUID, user *domain.UserUpdate) error
	DeleteById(id uuid.UUID) error
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
		users.GET("/", userHandler.GetAllHandler)
		users.GET("/:id", userHandler.GetByIdHandler)
		users.PATCH("/:id", userHandler.UpdateByIdHandler)
		users.DELETE("/:id", userHandler.DeleteByIdHandler)
	}
}

func (uh *userHandler) RegisterHandler(c *gin.Context) {
	var data domain.UserCreate

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	if err := uh.userService.Register(&data); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(data.ID))
}

// LoginHandler retrieves an user by its ID.
//
// @Summary      Get an user by ID
// @Description  This endpoint retrieves a single user by its unique identifier.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string                 true  "User ID"
// @Success      200  {object}  client.successRes     "User retrieved successfully"
// @Failure      400  {object}  client.AppError       "Invalid ID format or bad request"
// @Failure      404  {object}  client.AppError       "User not found"
// @Failure      500  {object}  client.AppError       "Internal Server Error"
// @Router       /users/{id} [get]
func (uh *userHandler) LoginHandler(c *gin.Context) {
	var data domain.UserLogin

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	token, err := uh.userService.Login(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(token))
}

// GetAllHandler retrieves all users.
//
// @Summary      Get all users
// @Description  This endpoint retrieves a list of all users.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  client.successRes  "List of users retrieved successfully"
// @Failure      500  {object}  client.AppError    "Internal Server Error"
// @Router       /users [get]
func (uh *userHandler) GetAllHandler(c *gin.Context) {
	var paging client.Paging
	var users []domain.User
	var err error

	if err := c.ShouldBind(&paging); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}
	paging.Process()

	requester := c.MustGet(client.CurrentUser).(client.Requester)
	if requester.GetRole() == domain.RoleAdmin.String() {
		users, err = uh.userService.GetAll(&paging)
	} else {
		c.JSON(http.StatusBadRequest, client.ErrNoPermission(err))
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, client.NewSuccessResponse(users, paging, nil))
}

// GetHandler retrieves an user by its ID.
//
// @Summary      Get an user by ID
// @Description  This endpoint retrieves a single user by its unique identifier.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string                 true  "User ID"
// @Success      200  {object}  client.successRes     "User retrieved successfully"
// @Failure      400  {object}  client.AppError       "Invalid ID format or bad request"
// @Failure      404  {object}  client.AppError       "User not found"
// @Failure      500  {object}  client.AppError       "Internal Server Error"
// @Router       /users/{id} [get]
func (uh *userHandler) GetByIdHandler(c *gin.Context) {
	var user *domain.User
	var err error

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	requester := c.MustGet(client.CurrentUser).(client.Requester)
	if requester.GetRole() == domain.RoleAdmin.String() {
		user, err = uh.userService.GetById(id)
	} else {
		user, err = uh.userService.GetById(requester.GetUserId())
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrCannotGetEntity(user.TableName(), err))
		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(user))
}

// UpdateHandler updates an existing user.
//
// @Summary      Update an user
// @Description  This endpoint allows updating the properties of an existing user by its ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "User ID"
// @Param        user  body      domain.UserUpdate      true  "User update payload"
// @Success      200   {object}  client.successRes     "User updated successfully"
// @Failure      400   {object}  client.AppError       "Invalid input or bad request"
// @Failure      404   {object}  client.AppError       "User not found"
// @Failure      500   {object}  client.AppError       "Internal Server Error"
// @Router       /users/{id} [put]
func (uh *userHandler) UpdateByIdHandler(c *gin.Context) {
	var user domain.UserUpdate

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	requester := c.MustGet(client.CurrentUser).(client.Requester)
	if requester.GetRole() == domain.RoleAdmin.String() {
		err = uh.userService.UpdateById(id, &user)
	} else {
		err = uh.userService.UpdateById(requester.GetUserId(), &user)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrCannotUpdateEntity(user.TableName(), err))
		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(true))
}

// DeleteHandler deletes an user by its ID.
//
// @Summary      Delete an user
// @Description  This endpoint deletes an user identified by its unique ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string                 true  "User ID"
// @Success      200  {object}  client.successRes     "User deleted successfully"
// @Failure      400  {object}  client.AppError       "Invalid ID format or bad request"
// @Failure      404  {object}  client.AppError       "User not found"
// @Failure      500  {object}  client.AppError       "Internal Server Error"
// @Router       /users/{id} [delete]
func (uh *userHandler) DeleteByIdHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	requester := c.MustGet(client.CurrentUser).(client.Requester)
	if requester.GetRole() == domain.RoleAdmin.String() {
		err = uh.userService.DeleteById(id)
	} else {
		err = uh.userService.DeleteById(requester.GetUserId())
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrCannotDeleteEntity(domain.User{}.TableName(), err))
		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(true))
}
