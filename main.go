package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name string `json:"name"`
	Price int `json:"price"`
}

var DB *gorm.DB

const DEFAULT_PORT = "1323"

func main() {
	var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("error: ", err)
	}

	DB.AutoMigrate(&Item{})

	log.Println("Connected to the database")

	app := echo.New()

	app.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "hello",
		})
	})

	app.GET("/items", func(c echo.Context) error {
		var items []Item
		
		if err := DB.Find(&items).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}

		return c.JSON(http.StatusOK, map[string]any{
			"message": "success",
			"data": items,
		})
	})

	app.POST("/items", func(c echo.Context) error {
		var item Item
		if err := c.Bind(&item); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}

		if err := DB.Create(&item).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}

		return c.JSON(http.StatusOK, map[string]any{
			"message": "item created!",
			"data": item,
		})
	})

	var port string = os.Getenv("PORT")

	if port == "" {
		port = DEFAULT_PORT
	}

	var appPort string = fmt.Sprintf(":%s", port)

	app.Logger.Fatal(app.Start(appPort))
}