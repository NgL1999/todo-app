package gin

import (
	"context"
	"net/http"
	"time"
	"todo-app/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IItemService interface {
	CreateItem(ctx context.Context, item *domain.Item) error
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
	}
}

func (ih *itemHandler) CreateItemHandler(c *gin.Context) {
		item := domain.Item{}

		if err := c.ShouldBind(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		item.ID = uuid.New()
		item.Created_at = time.Now()
		item.Updated_at = time.Now()

		// if err := db.Create(&item).Error; err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		c.JSON(http.StatusCreated, gin.H{
			"data": item.ID,
		})
	}
