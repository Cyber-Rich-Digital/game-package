package main

import (
	"context"
	docs "cyber-api/docs"
	backend "cyber-api/handler/backend"
	frontend "cyber-api/handler/frontend"
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

	firebase, _ := initFirebase()

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	initTimeZone()
	db := initDatabase()

	r := gin.Default()
	gin.SetMode(os.Getenv("GIN_MODE"))

	// corsConfig := cors.DefaultConfig()

	// corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	// // To be able to send tokens to the server.
	// corsConfig.AllowCredentials = true
	// corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	// // corsConfig.AllowHeaders = []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"}
	// // OPTIONS method for ReactJS
	// corsConfig.AddAllowMethods("OPTIONS", "GET", "POST", "PUT", "DELETE", "PATCH")
	// corsConfig.AllowHeaders = []string{"Access-Control-Allow-Headers", "Origin", "Accept", "Content-Type", "Authorization", "authorization", "X-Requested-With", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Credentials", "X-Auth-Token", "X-Auth-Email", "X-Auth-Password", "X-Auth-Name", "X-Auth-Phone", "X-Auth-Role", "X-Auth-Website", "X-Auth-Device", "X-Auth-Tag", "X-Auth-Message", "X-Auth-Website-Id", "X-Auth-Device-Id", "X-Auth-Tag-Id", "X-Auth-Message-Id", "X-Auth-Website-Name", "X-Auth-Device-Name", "X-Auth-Tag-Name", "X-Auth-Message-Name", "X-Auth-Website-Url", "X-Auth-Device-Url", "X-Auth-Tag-Url", "X-Auth-Message-Url", "X-Auth-Website-Description", "X-Auth-Device-Description", "X-Auth-Tag-Description", "X-Auth-Message-Description", "X-Auth-Website-Image", "X-Auth-Device-Image", "X-Auth-Tag-Image", "X-Auth-Message-Image", "X-Auth-Website-Status", "X-Auth-Device-Status", "X-Auth-Tag-Status", "X-Auth-Message-Status", "X-Auth-Website-User", "X-Auth-Device-User", "X-Auth-Tag-User", "X-Auth-Message-User", "X-Auth-Website-User-Id", "X-Auth-Device-User-Id", "X-Auth-Tag-User-Id", "X-Auth-Message-User-Id", "X-Auth-Website-User-Name", "X-Auth-Device-User-Name", "X-Auth-Tag-User-Name", "X-Auth-Message-User-Name", "X-Auth-Website-User-Email", "X-Auth-Device-User-Email", "X-Auth-Tag-User-Email", "X-Auth-Message-User-Email", "X-Auth-Website-User-Password"}

	// Register the middleware
	r.Use(CORSMiddleware())

	path := "/api"
	route := r.Group(path)

	docs.SwaggerInfo.BasePath = path

	route.GET("/ping", func(c *gin.Context) {
		pingExample(c)
	})

	frontPath := "/api"
	frontRoute := r.Group(frontPath)

	frontend.TagController(frontRoute, db)
	frontend.WebsiteController(frontRoute, db)
	frontend.MessageController(frontRoute, db, firebase)
	frontend.DeviceController(frontRoute, db)

	backPath := "/api/be"
	backRoute := r.Group(backPath)
	backend.AuthController(backRoute, db)
	backend.UserController(backRoute, db)
	backend.WebsiteController(backRoute, db)
	backend.TagController(backRoute, db)
	backend.MenuController(backRoute)

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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
