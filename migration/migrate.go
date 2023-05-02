package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"sort"
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

var db *gorm.DB
var lastVersion migration
var files []fs.DirEntry
var sortList []fs.DirEntry

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("Please input command up or down")
	}

	if err := godotenv.Load("migration.env"); err != nil {
		panic(err)
	}

	db = initDatabase()

	// check table migration if not exist then create table migration
	checkTableMigrate(db)

	// get last version from table migration

	if err := db.First(&lastVersion).Error; err != nil {
		panic(err)
	}

	// read sql files from directory migration
	var fileErr error
	files, fileErr = os.ReadDir("./migration")
	if fileErr != nil {
		panic(fileErr)
	}

	if len(args) == 1 && args[0] == "down" {
		log.Fatal("Please input version")
	}

	if len(args) == 2 && args[0] == "down" {

		version, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal("Please input number")
		}

		if version < 1 {
			log.Fatal("Please input version greater than 0")
		}
	}

	if args[0] == "up" {
		Up()
	} else if args[0] == "down" {
		Down(args[1])
	} else {
		log.Fatal("Command not found")
	}
}

func Up() {

	// sort file name in array of files
	sort.Slice(files, func(i, j int) bool {

		nameI := strings.Split(files[i].Name(), "-")[0]
		nameJ := strings.Split(files[j].Name(), "-")[0]

		intI, _ := strconv.Atoi(nameI)
		intJ, _ := strconv.Atoi(nameJ)

		return intI < intJ
	})

	for _, file := range files {
		if strings.Contains(file.Name(), "up") && !strings.Contains(file.Name(), "migrate") {
			sortList = append(sortList, file)
		}
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
	for _, file := range sortList {

		getVersion := strings.Split(file.Name(), "-")[0]

		version, _ := strconv.Atoi(getVersion)

		// check if file version is greater than last version
		if version > lastVersion.Version {

			// read sql file
			sql, err := os.ReadFile("./migration/" + file.Name())
			if err != nil {
				// print color red
				println(fmt.Sprintf("\033[31m%s\033[0m ", "Error: "+err.Error()))
				panic(err)
			}

			list := strings.Split(string(sql), ";")

			for _, query := range list {

				if query == "" {
					continue
				}

				// execute sql file
				if err := tx.Exec(query).Error; err != nil {
					// print color red
					println(fmt.Sprintf("\033[31m%s\033[0m ", "Error: "+err.Error()))
					tx.Rollback()
					panic(err)
				}
			}

			// print color green
			println(fmt.Sprintf("\033[32m%s\033[0m ", "Query file: "+file.Name()))

			// update last version
			if err := tx.Model(&migration{}).Where("id = ?", lastVersion.ID).Update("version", version).Error; err != nil {
				// print color red
				println(fmt.Sprintf("\033[31m%s\033[0m ", "Error: "+err.Error()))
				tx.Rollback()
				panic(err)
			}
		}
	}

	tx.Commit()

	println(fmt.Sprintf("\033[32m%s\033[0m ", "Up version completed"))
}

func Down(target string) {

	// sort file name in array of files
	sort.Slice(files, func(i, j int) bool {

		nameI := strings.Split(files[i].Name(), "-")[0]
		nameJ := strings.Split(files[j].Name(), "-")[0]

		intI, _ := strconv.Atoi(nameI)
		intJ, _ := strconv.Atoi(nameJ)

		return intI > intJ
	})

	for _, file := range files {
		if strings.Contains(file.Name(), "down") {
			sortList = append(sortList, file)
		}
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
	for _, file := range sortList {

		getVersion := strings.Split(file.Name(), "-")[0]
		version, _ := strconv.Atoi(getVersion)
		target, _ := strconv.Atoi(target)

		if version < target {
			break
		}

		// check if file version is greater than last version
		if version <= lastVersion.Version {

			// read sql file
			sql, err := os.ReadFile("./migration/" + file.Name())
			if err != nil {
				// print color red
				println(fmt.Sprintf("\033[31m%s\033[0m ", "Error: "+err.Error()))
				panic(err)
			}

			list := strings.Split(string(sql), ";")

			for _, query := range list {

				if query == "" {
					continue
				}

				// execute sql file
				if err := tx.Exec(query).Error; err != nil {
					// print color red
					println(fmt.Sprintf("\033[31m%s\033[0m ", "Error: "+err.Error()))
					tx.Rollback()
					panic(err)
				}
			}

			// print color green
			println(fmt.Sprintf("\033[32m%s\033[0m ", "Query file: "+file.Name()))

			// update last version
			if err := tx.Model(&migration{}).Where("id = ?", lastVersion.ID).Update("version", version-1).Error; err != nil {
				// print color red
				println(fmt.Sprintf("\033[31m%s\033[0m ", "Error: "+err.Error()))
				tx.Rollback()
				panic(err)
			}
		}
	}

	tx.Commit()

	println(fmt.Sprintf("\033[32m%s\033[0m ", "Down version completed"))
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
