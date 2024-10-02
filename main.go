package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Item struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

func (Item) TableName() string {
	return "items"
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Connect database
	db, err := gorm.Open(postgres.Open(os.Getenv("CONNECTION_STRING")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(db)

	// Create APIs
	r := gin.Default()

	api := r.Group("v1")
	{
		items := api.Group("items")
		items.POST("/", CreateItem(db))
		items.GET("/all", GetAllItems(db))
		items.GET("/:id", GetItemById(db))
		items.PATCH("/:id", UpdateItemById(db))
		items.DELETE("/:id", DeleteItemById(db))
	}

	r.Run()
}

func CreateItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		item := Item{}

		if err := c.ShouldBind(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		item.ID = uuid.New()
		item.Created_at = time.Now()
		item.Updated_at = time.Now()

		if err := db.Create(&item).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data": item.ID,
		})
	}
}

func GetAllItems(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		items := []Item{}

		if err := db.Find(&items).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": items,
		})
	}
}

func GetItemById(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		item := Item{}
		id := c.Param("id")

		if err := db.Find(&item, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": item,
		})
	}
}

func UpdateItemById(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		item := Item{}

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
		item.Updated_at = time.Now()

		result := db.Model(&item).Updates(
			Item{
				Title:       item.Title,
				Description: item.Description,
				Updated_at:  item.Updated_at,
			})

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": result.Error,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"updatedCount": result.RowsAffected,
		})
	}
}

func DeleteItemById(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		items := Item{}
		id := c.Param("id")
		result := db.Delete(&items, "id = ?", id)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": result.Error,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"deletedCount": result.RowsAffected,
		})
	}
}
