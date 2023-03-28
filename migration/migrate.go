package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type migration struct {
	ID        int `gorm:"primaryKey" autoIncrement:"true"`
	Version   int
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func main() {

	if err := godotenv.Load("migration.env"); err != nil {
		panic(err)
	}

	db := initDatabase()

	// check table migration if not exist then create table migration
	checkTableMigrate(db)

	// get last version from table migration
	var lastVersion migration
	if err := db.First(&lastVersion).Error; err != nil {
		panic(err)
	}

	// read sql files from directory migration
	files, err := os.ReadDir("./migration")
	if err != nil {
		panic(err)
	}

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		panic(err)
	}

	// execute sql files
	for _, file := range files {

		if file.IsDir() {
			continue
		}

		fileName := strings.Split(file.Name(), ".")[0]

		if fileName == "migrate" {
			continue
		}

		fileNameInt, _ := strconv.Atoi(fileName)

		// check if file version is greater than last version
		if fileNameInt > lastVersion.Version {

			// read sql file
			sql, err := os.ReadFile("./migration/" + file.Name())
			if err != nil {
				// print color red
				println(fmt.Sprintf("\033[31m%s\033[0m ", "Error: "+err.Error()))
				panic(err)
			}

			list := strings.Split(string(sql), ";")

			for _, query := range list {

				if strings.TrimSpace(query) == "" {
					continue
				}

				// execute sql file
				if err := tx.Exec(query).Error; err != nil {
					// print color red
					println(fmt.Sprintf("\033[31m%s\033[0m ", "Error: "+err.Error()))
					// You Cant rollback on Create/Alter
					// tx.Rollback()
					panic(err)
				}
			}

			// print color green
			println(fmt.Sprintf("\033[32m%s\033[0m ", "Query file: "+file.Name()))

			// update last version
			if err := tx.Model(&migration{}).Where("id = ?", lastVersion.ID).Update("version", fileNameInt).Error; err != nil {
				// print color red
				println(fmt.Sprintf("\033[31m%s\033[0m ", "Error: "+err.Error()))
				tx.Rollback()
				panic(err)
			}

		}

	}

	tx.Commit()

	println(fmt.Sprintf("\033[32m%s\033[0m ", "Migration is done"))

}

func checkTableMigrate(db *gorm.DB) {

	// insert first record to table migration
	if !db.Migrator().HasTable(&migration{}) {
		if err := db.Migrator().CreateTable(&migration{}); err != nil {
			panic(err)
		}

		if err := db.Create(&migration{Version: 0}).Error; err != nil {
			panic(err)
		}
	}
}

func initDatabase() *gorm.DB {

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(err)
	}

	fmt.Println("Database is connected")

	return db
}
