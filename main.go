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
		items.GET("/all", GetAll(db))
		items.PATCH("/:id", UpdateById(db))
		// items.GET("/:id")
		// items.DELETE("/:id")
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

func GetAll(db *gorm.DB) func(c *gin.Context) {
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

func UpdateById(db *gorm.DB) func(c *gin.Context) {
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

		result := db.Model(&item).Updates(Item{Title: item.Title, Description: item.Description})

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
