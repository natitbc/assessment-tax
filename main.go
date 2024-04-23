package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Tax struct {
	Tax float64 `json:"tax"`
}

type Err struct {
	Message string `json:"message"`
}

var tax = []Tax{
	{Tax: 29000.0},
}

func createTaxHandler(c echo.Context) error {
	t := Tax{}
	err := c.Bind(&t)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	tax = append(tax, t)
	fmt.Println("id : % #v\n", t)
	return c.JSON(http.StatusCreated, t)
}

func getTaxHandler(c echo.Context) error {
	fmt.Print("tax : % #v\n", tax)
	return c.JSON(http.StatusOK, tax)
}

func main() {
	e := echo.New()

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading environment variables:", err)
	}

	// Use os.Getenv to access environment variables
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	fmt.Println(os.Getenv("ADMIN_USERNAME"))

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == adminUsername && password == adminPassword {
			return true, nil
		}
		return false, nil
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/tax", createTaxHandler)
	e.GET("/tax", getTaxHandler)

	log.Fatal(e.Start(":8080"))
}
