package main

import (
	"encoding/json"
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

type TaxResponse struct {
	Tax      float64                `json:"tax"`
	TaxLevel []calculation.TaxLevel `json:"taxLevel"` // Use the actual type from the calculation package
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

	fmt.Print(allowancesdata)

	tax, CalculatedTaxLevel, _ := calculation.CalculateTax(totalincome, wht, []calculation.Allowance{
		{AllowanceType: "donation", Amount: allowancesdata[0].Amount},
	})

	responseTax := &TaxResponse{
		Tax:      tax,
		TaxLevel: CalculatedTaxLevel,
	}

	return c.JSON(http.StatusOK, responseTax)
}

func getTaxHandler(c echo.Context) error {
	fmt.Print("tax : % #v\n", responseTax)
	return c.JSON(http.StatusOK, responseTax)
}

type UpdatePersonalDeductionRequest struct {
	PersonalDeduction float64 `json:"amount"`
}

func setDeductionsHandler(c echo.Context, config *calculation.Config) error {
	// if isAdmin {
	// 	return c.JSON(http.StatusUnauthorized, Err{Message: "Unauthorized"})
	// }

	var req UpdatePersonalDeductionRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if req.PersonalDeduction < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid personal deduction amount"})
	}

	config.PersonalDeduction = req.PersonalDeduction

	data, err := json.Marshal(config)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	err = os.WriteFile("calculation/config.json", data, 0644)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Personal deduction updated successfully"})

}

func main() {
	e := echo.New()

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading environment variables:", err)
	}

	// Use os.Getenv to access environment variables
	// adminUsername := os.Getenv("ADMIN_USERNAME")
	// adminPassword := os.Getenv("ADMIN_PASSWORD")

	// fmt.Println(os.Getenv("ADMIN_USERNAME"))

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "adminTax" && password == "admin!" {
			return true, nil
		}
		return false, nil
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/admin/deductions/personal", func(c echo.Context) error {
		return setDeductionsHandler(c, &calculation.Config{})
	})

	e.POST("/tax/calculation", createTaxHandler)
	e.GET("/tax/calculation", getTaxHandler)

	log.Fatal(e.Start(":8080"))
}
