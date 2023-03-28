package main

import (
	"context"
	docs "cybergame-api/docs"
	handler "cybergame-api/handler"
	"cybergame-api/middleware"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"google.golang.org/api/option"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	// firebase, _ := initFirebase()

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	initTimeZone()
	db := initDatabase()

	r := gin.Default()
	gin.SetMode(os.Getenv("GIN_MODE"))

	// Register the middleware
	r.Use(middleware.CORSMiddleware())

	path := "/api"
	route := r.Group(path)

	docs.SwaggerInfo.BasePath = path

	route.GET("/ping", func(c *gin.Context) {
		pingExample(c)
	})

	backRoute := r.Group(path)
	handler.AuthController(backRoute, db)
	handler.UserController(backRoute, db)
	handler.WebsiteController(backRoute, db)
	handler.TagController(backRoute, db)
	handler.MenuController(backRoute)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	err := r.Run(port)
	if err != nil {
		panic(err)
	}
}

func initFirebase() (*firebase.App, context.Context) {

	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("firebase_account_key.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	firebase, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("error initializing app: %v")
	}

	log.Println("Firebase initialized")

	return firebase, ctx
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

	ict, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		panic(err)
	}

	time.Local = ict

	println("Time now", time.Now().Format("2006-01-02 15:04:05"))
}

func initDatabase() *gorm.DB {

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(err)
	}

	// _, offset := time.Now().Zone()
	// tz := fmt.Sprintf("%+03d:00", (offset / 3600))
	// println(fmt.Sprintf("set time_zone = '%s';", tz))
	// if err := db.Exec(fmt.Sprintf("set time_zone = '%s';", tz)).Error; err != nil {
	// 	println(fmt.Sprintf("\033[31m%s\033[0m ", "Error: "+err.Error()))
	// 	panic(err)
	// }

	println("Database is connected")

	return db
}
