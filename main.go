package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// "github.com/natitbc/assessment-tax/calculation"
	"github.com/natitbc/assessment-tax/calculation"
)

type TaxLevel struct {
	level string
	tax   float64
}
type Tax struct {
	Tax      float64    `json:"tax"`
	TaxLevel []TaxLevel `json:"taxLevel"`
	// err error
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type userData struct {
	TotalIncome float64     `json:"totalIncome"`
	Wht         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Err struct {
	Message string `json:"message"`
}

var responseTax = []Tax{
	{
		Tax: 0.0,
		TaxLevel: []TaxLevel{
			{level: "0-150,000", tax: 0.0},
			{level: "150,001-500,000", tax: 0.0},
			{level: "500,001-1,000,000", tax: 0.0},
			{level: "1,000,001-2,000,000", tax: 0.0},
			{level: "2,000,001 ขึ้นไป", tax: 0.0},
		},
	},
}

func createTaxHandler(c echo.Context) error {
	var data userData

	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	totalincome := data.TotalIncome
	wht := data.Wht
	allowancesdata := data.Allowances
	fmt.Println(totalincome, wht, allowancesdata)

	tax, CalculatedTaxLevel, _ := calculation.CalculateTax(totalincome, wht, []calculation.Allowance{
		{AllowanceType: "donation", Amount: allowancesdata[0].Amount},
	})
	responseTax[0].Tax = tax
	fmt.Println("--CalculatedTaxLevel : ", CalculatedTaxLevel)

	for i := 0; i < len(CalculatedTaxLevel); i++ {
		// fmt.Println("TaxLevel : ", CalculatedTaxLevel[i].Tax)
		responseTax[0].TaxLevel[i].tax = CalculatedTaxLevel[i].Tax
	}
	fmt.Println("responseTax : ", responseTax)

	fmt.Println("tax data : ", tax)

	return c.JSON(http.StatusOK, responseTax)
}

func getTaxHandler(c echo.Context) error {
	fmt.Print("tax : % #v\n", responseTax)
	return c.JSON(http.StatusOK, responseTax)
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

	e.POST("/tax/calculation", createTaxHandler)
	e.GET("/tax/calculation", getTaxHandler)

	log.Fatal(e.Start(":8080"))
}
