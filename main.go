package main

import (
	"log"
	"os"
	docs "todo-app/docs"
	restApi "todo-app/internal/api/http/gin"
	pgRepo "todo-app/internal/repository/postgres"
	"todo-app/item"
	"todo-app/pkg/tokenprovider/jwt"
	"todo-app/pkg/util"
	"todo-app/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r := gin.Default()

	hasher := util.NewMd5Hash()
	tokenProvider := jwt.NewJWTProvider()

	// Swagger
	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Create APIs
	api := r.Group("v1")
	{
		itemRepo := pgRepo.NewItemRepo(db)
		itemService := item.NewItemService(itemRepo)
		restApi.NewItemHandler(api, itemService)
		userRepo := pgRepo.NewUserRepo(db)
		userService := user.NewUserService(userRepo, hasher, tokenProvider, 60*60*24*30)
		restApi.NewUserHandler(api, userService)
	}
	r.Run()
}
