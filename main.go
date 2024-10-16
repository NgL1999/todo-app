package main

import (
	"log"
	"os"
	"time"
	"todo-app/docs"
	restApi "todo-app/internal/api/http/gin"
	"todo-app/internal/api/http/gin/middleware"
	pgRepo "todo-app/internal/repository/postgres"
	"todo-app/item"
	"todo-app/pkg/memcache"
	"todo-app/pkg/tokenprovider/jwt"
	"todo-app/pkg/util"
	"todo-app/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
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

	// ─── Utils ───────────────────────────────────────────────────────────
	hasher := util.NewMd5Hash()
	tokenProvider := jwt.NewJWTProvider(os.Getenv("SECRET_KEY"))
	tokenExpire := 60 * 60 * 24 * 30

	// ─── Swagger ─────────────────────────────────────────────────────────
	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ─── Error Handler ───────────────────────────────────────────────────
	r.Use(middleware.Recover())

	// ─── Repos ───────────────────────────────────────────────────────────
	userRepo := pgRepo.NewUserRepo(db)
	itemRepo := pgRepo.NewItemRepo(db)

	// ─── Services ────────────────────────────────────────────────────────
	userService := user.NewUserService(userRepo, hasher, tokenProvider, tokenExpire)
	itemService := item.NewItemService(itemRepo)

	// ─── Base Api ────────────────────────────────────────────────────────
	api := r.Group("v1")
	
	// ─── Middlewares ─────────────────────────────────────────────────────
	// Auth
	authCache := memcache.NewUserCaching(memcache.NewRedisCache(), userRepo)
	middlewareAuth := middleware.RequiredAuth(tokenProvider, authCache)

	// Cache
	limiterRate := limiter.Rate{
		Period: 5 * time.Second,
		Limit:  3,
	}
	store := memory.NewStore()
	limiter := limiter.New(store, limiterRate)
	middlewareRateLimit := middleware.RateLimiter(limiter)

	// ─── Handlers ───────────────────────────────────────────────────────────
	restApi.NewUserHandler(api, userService, middlewareAuth)
	restApi.NewItemHandler(api, itemService, middlewareAuth, middlewareRateLimit)

	r.Run()
}
