package gin

import (
	"context"
	"net/http"
	"todo-app/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IItemService interface {
	CreateItem(ctx context.Context, item *domain.Item) error
	GetAllItems(ctx context.Context, items *[]domain.Item) error
	GetItemById(ctx context.Context, item *domain.Item, id string) error
	UpdateItemById(ctx context.Context, item *domain.Item) (int64, error)
	DeleteItemById(ctx context.Context, item *domain.Item) (int64, error)
}

type itemHandler struct {
	itemService IItemService
}

// Constructor
func NewItemHandler(apiVersion *gin.RouterGroup, isvc IItemService) {
	itemHandler := &itemHandler{
		itemService: isvc,
	}

	items := apiVersion.Group("items")
	{
		items.POST("/", itemHandler.CreateItemHandler)
		items.GET("/all", itemHandler.GetAllItemsHandler)
		items.GET("/:id", itemHandler.GetItemByIdHandler)
		items.PATCH("/:id", itemHandler.UpdateItemByIdHandler)
		items.DELETE("/:id", itemHandler.DeleteItemByIdHandler)
	}
}

// @Summary Show an item
// @Description get string by ID
// @Accept  json
// @Produce  json
// @Param   id path int true "Item ID"
// @Success 200 {object} Item
// @Failure 400 {object} HTTPError
// @Router /items [post]
func (ih *itemHandler) CreateItemHandler(c *gin.Context) {
	item := domain.Item{}

	if err := c.ShouldBind(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	item.ID = uuid.New()

	if err := ih.itemService.CreateItem(c, &item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": item.ID,
	})
}

func (ih *itemHandler) GetAllItemsHandler(c *gin.Context) {
	items := []domain.Item{}

	if err := ih.itemService.GetAllItems(c, &items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

func (ih *itemHandler) GetItemByIdHandler(c *gin.Context) {
	item := domain.Item{}
	id := c.Param("id")

	if err := ih.itemService.GetItemById(c, &item, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

func (ih *itemHandler) UpdateItemByIdHandler(c *gin.Context) {
	item := domain.Item{}

	if err := c.ShouldBind(&item); err != nil {
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

	item.ID = id
	result, err := ih.itemService.UpdateItemById(c, &item)

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

func (ih *itemHandler) DeleteItemByIdHandler(c *gin.Context) {
	item := domain.Item{}
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	item.ID = id
	result, err := ih.itemService.DeleteItemById(c, &item)

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
