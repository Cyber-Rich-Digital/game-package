package main

import (
	docs "cyber-api/docs"
	backend "cyber-api/handler/backend"
	frontend "cyber-api/handler/frontend"
	"cyber-api/middleware"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	initTimeZone()
	db := initDatabase()

	r := gin.New()
	gin.SetMode(os.Getenv("GIN_MODE"))

	path := "/api"
	route := r.Group(path)

	docs.SwaggerInfo.BasePath = path
	docs.SwaggerInfo.Title = "Cyber API"

	route.GET("/ping", middleware.Authorize, func(c *gin.Context) {
		pingExample(c)
	})

	frontPath := "/api"
	frontRoute := r.Group(frontPath)

	frontend.PromotionController(frontRoute, db)
	frontend.AuthController(frontRoute, db)
	frontend.UserController(frontRoute, db)

	backPath := "/api/be"
	backRoute := r.Group(backPath)
	backend.PromotionController(backRoute, db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	err := r.Run(port)
	if err != nil {
		panic(err)
	}
}

type ping struct {
	Message string `json:"message" example:"pong" `
}

// @BasePath /ping
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object} ping
// @Router /ping [get]
func pingExample(c *gin.Context) {
	c.JSON(200, ping{Message: "pong"})
}

func initTimeZone() {

	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict

	println("Time now", time.Now().Format("2006-01-02 15:04:05"))
}

func initDatabase() *gorm.DB {

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(err)
	}

	println("Database is connected")

	return db
}
