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
	}

	r.Run()
}
