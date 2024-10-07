package main

import (
	"log"
	"os"
	restApi "todo-app/internal/api/http/gin"
	pgRepo "todo-app/internal/repository/postgres"
	"todo-app/item"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
		itemRepo := pgRepo.NewItemRepo(db)
		itemService := item.NewItemService(itemRepo)
		restApi.NewItemHandler(api, itemService)
		// items.DELETE("/:id", DeleteItemById(db))
	}

	r.Run()
}


// func DeleteItemById(db *gorm.DB) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		items := Item{}
// 		id := c.Param("id")
// 		result := db.Delete(&items, "id = ?", id)

// 		if result.Error != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": result.Error,
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"deletedCount": result.RowsAffected,
// 		})
// 	}
// }
