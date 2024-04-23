package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Tax struct {
	Tax float64 `json:"tax"`
}

type Err struct {
	Message string `json:"message"`
}

var users = []User{
	{ID: 1, Name: "AnuchitO", Age: 20},
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

func createUserHandler(c echo.Context) error {
	u := User{}
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	users = append(users, u)

	fmt.Println("id : % #v\n", u)

	return c.JSON(http.StatusCreated, u)
}

func getTaxHandler(c echo.Context) error {
	fmt.Print("tax : % #v\n", tax)
	return c.JSON(http.StatusOK, tax)
}

func getUsersHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

func main() {
	e := echo.New()

	//Set up environment variables
	envVars := map[string]string{
		"PORT":         "8080",
		"DATABASE_URL": "host={REPLACE_ME} port=5432 ADMIN_USERNAME={adminTax} ADMIN_PASSWORD={admin!} dbname={45678} sslmode=disable",
	}

	for name, value := range envVars {
		err := os.Setenv(name, value)
		if err != nil {
			fmt.Printf("Failed to set env var %s: %s\n", name, err)
		}
	}

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv("ADMIN_USERNAME") || os.Getenv("ADMIN_PASSWORD") == password {
			return true, nil
		}
		return false, nil
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/users", createUserHandler)
	e.POST("/tax", createTaxHandler)
	e.GET("/users", getUsersHandler)
	e.GET("/tax", getTaxHandler)

	log.Fatal(e.Start(":8080"))
}
