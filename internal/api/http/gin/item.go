package gin

import (
	"net/http"
	"todo-app/domain"
	"todo-app/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IItemService interface {
	CreateItem(item *domain.ItemCreation) error
	GetAllItems(userID uuid.UUID, paging *client.Paging) ([]domain.Item, error)
	GetItemById(id, userID uuid.UUID) (domain.Item, error)
	UpdateItemById(id, userID uuid.UUID, item *domain.ItemUpdate) error
	DeleteItemById(id, userID uuid.UUID) error
}

type itemHandler struct {
	itemService IItemService
}

func NewItemHandler(apiVersion *gin.RouterGroup, isvc IItemService, middlewareAuth func(c *gin.Context), middlewareRateLimit func(c *gin.Context)) {
	itemHandler := &itemHandler{
		itemService: isvc,
	}

	items := apiVersion.Group("items", middlewareAuth)
	{
		items.POST("/", itemHandler.CreateItemHandler)
		items.GET("/", middlewareRateLimit, itemHandler.GetAllItemsHandler)
		items.GET("/:id", itemHandler.GetItemByIdHandler)
		items.PATCH("/:id", itemHandler.UpdateItemByIdHandler)
		items.DELETE("/:id", itemHandler.DeleteItemByIdHandler)
	}
}

// CreateItemHandler handles the creation of a new item.
//
// @Summary      Create a new item
// @Description  This endpoint allows authenticated users to create an item.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        item  body      domain.ItemCreation  true  "Item creation payload"
// @Success      200   {object}  client.SuccessRes   "Item successfully created"
// @Failure      400   {object}  client.AppError     "Bad Request"
// @Failure      401   {object}  client.AppError     "Unauthorized"
// @Failure      500   {object}  client.AppError     "Internal Server Error"
func (ih *itemHandler) CreateItemHandler(c *gin.Context) {
	var item domain.ItemCreation

	if err := c.ShouldBind(&item); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	requester := c.MustGet(client.CurrentUser).(client.Requester)
	item.UserID = requester.GetUserId()

	if err := ih.itemService.CreateItem(&item); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, client.SimpleSuccessResponse(item.ID))
}

// GetAllItemsHandler retrieves all items.
//
// @Summary      Get all items
// @Description  This endpoint retrieves a list of all items.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Success      200  {object}  clients.SuccessRes  "List of items retrieved successfully"
// @Failure      500  {object}  clients.AppError    "Internal Server Error"
// @Router       /items [get]
func (ih *itemHandler) GetAllItemsHandler(c *gin.Context) {
	var paging client.Paging
	if err := c.ShouldBind(&paging); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}
	paging.Process()

	requester := c.MustGet(client.CurrentUser).(client.Requester)

	items, err := ih.itemService.GetAllItems(requester.GetUserId(), &paging)
	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, client.NewSuccessResponse(items, paging, nil))
}

// GetItemHandler retrieves an item by its ID.
//
// @Summary      Get an item by ID
// @Description  This endpoint retrieves a single item by its unique identifier.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        id   path      string                 true  "Item ID"
// @Success      200  {object}  clients.SuccessRes     "Item retrieved successfully"
// @Failure      400  {object}  clients.AppError       "Invalid ID format or bad request"
// @Failure      404  {object}  clients.AppError       "Item not found"
// @Failure      500  {object}  clients.AppError       "Internal Server Error"
// @Router       /items/{id} [get]
func (ih *itemHandler) GetItemByIdHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	requester := c.MustGet(client.CurrentUser).(client.Requester)

	item, err := ih.itemService.GetItemById(id, requester.GetUserId())
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(item))
}

// UpdateItemHandler updates an existing item.
//
// @Summary      Update an item
// @Description  This endpoint allows updating the properties of an existing item by its ID.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "Item ID"
// @Param        item  body      domain.ItemUpdate      true  "Item update payload"
// @Success      200   {object}  clients.SuccessRes     "Item updated successfully"
// @Failure      400   {object}  clients.AppError       "Invalid input or bad request"
// @Failure      404   {object}  clients.AppError       "Item not found"
// @Failure      500   {object}  clients.AppError       "Internal Server Error"
// @Router       /items/{id} [put]
func (ih *itemHandler) UpdateItemByIdHandler(c *gin.Context) {
	var item domain.ItemUpdate

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	if err := c.ShouldBind(&item); err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	requester := c.MustGet(client.CurrentUser).(client.Requester)

	if err := ih.itemService.UpdateItemById(id, requester.GetUserId(), &item); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(true))
}

// DeleteItemHandler deletes an item by its ID.
//
// @Summary      Delete an item
// @Description  This endpoint deletes an item identified by its unique ID.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        id   path      string                 true  "Item ID"
// @Success      200  {object}  clients.SuccessRes     "Item deleted successfully"
// @Failure      400  {object}  clients.AppError       "Invalid ID format or bad request"
// @Failure      404  {object}  clients.AppError       "Item not found"
// @Failure      500  {object}  clients.AppError       "Internal Server Error"
// @Router       /items/{id} [delete]
func (ih *itemHandler) DeleteItemByIdHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(err))
		return
	}

	requester := c.MustGet(client.CurrentUser).(client.Requester)

	if err := ih.itemService.DeleteItemById(id, requester.GetUserId()); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, client.SimpleSuccessResponse(true))
}
