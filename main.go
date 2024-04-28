package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// "github.com/natitbc/assessment-tax/calculation"
	"github.com/natitbc/assessment-tax/calculation"
	"github.com/natitbc/assessment-tax/config"
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
	Tax       float64                `json:"tax"`
	TaxRefund float64                `json:"taxRefund"`
	TaxLevel  []calculation.TaxLevel `json:"taxLevel"` // Use the actual type from the calculation package
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

	var tax float64
	var CalculatedTaxLevel []calculation.TaxLevel

	if wht < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "wht cannot be negative"})
	}
	if wht > totalincome {
		return c.JSON(http.StatusBadRequest, Err{Message: "wht cannot be greater than total income"})
	}

	//check has allowancedata[0].AllowanceType
	if len(allowancesdata) == 0 {
		fmt.Println("Tax with No allowances")
		tax, CalculatedTaxLevel, _ = calculation.CalculateTax(totalincome, wht, []calculation.Allowance{})

	} else {
		fmt.Println("Tax with allowances")
		countAllowance := len(allowancesdata)

		if countAllowance > 2 {
			return c.JSON(http.StatusBadRequest, Err{Message: "max 2 allowances"})
		}
		for i := 0; i < countAllowance; i++ {
			if allowancesdata[i].Amount < 0 {
				return c.JSON(http.StatusBadRequest, Err{Message: "amount must be positive"})
			}

			if allowancesdata[i].AllowanceType == "donation" || allowancesdata[i].AllowanceType == "k-receipt" {
				tax, CalculatedTaxLevel, _ = calculation.CalculateTax(totalincome, wht, []calculation.Allowance{
					{AllowanceType: allowancesdata[i].AllowanceType, Amount: allowancesdata[i].Amount},
				})
			} else {
				return c.JSON(http.StatusBadRequest, Err{Message: "allowanceType must be donation or k-receipt"})
			}
		}
	}

	TaxRefund := 0.0
	if tax < 0 {
		TaxRefund = tax * -1
	}

	responseTax := &TaxResponse{
		Tax:       tax,
		TaxRefund: TaxRefund,
		TaxLevel:  CalculatedTaxLevel,
	}

	return c.JSON(http.StatusOK, responseTax)
}

type TaxData struct {
	TotalIncome float64 `csv:"totalIncome"`
	Wht         float64 `csv:"wht"`
	Donation    float64 `csv:"donation"`
}

type TaxResult struct {
	TotalIncome float64
	Tax         float64
}

func upload(c echo.Context) error {
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("taxFile")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	fmt.Println("Successfully opened the CSV file")

	defer src.Close()

	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	fd, error := os.OpenFile("taxes.csv", os.O_RDWR, 0644)

	if error != nil {
		fmt.Println(error)
	}

	fmt.Println("Successfully opened the CSV file")
	defer fd.Close()

	// read csv data
	fileBytes, err := os.ReadFile("taxes.csv")
	if err != nil {
		fmt.Println(err)
	}

	var TaxData []TaxData

	gocsv.UnmarshalBytes(fileBytes, &TaxData)
	fmt.Println(TaxData)

	var results []TaxResult

	for _, entry := range TaxData {
		tax, _, err := calculation.CalculateTax(entry.TotalIncome, entry.Wht, []calculation.Allowance{
			{AllowanceType: "donation", Amount: entry.Donation},
		})
		if err != nil {
			return err
		}
		results = append(results, TaxResult{
			TotalIncome: entry.TotalIncome,
			Tax:         tax,
		})
	}

	fmt.Println(results)

	response := struct {
		Taxes []TaxResult `json:"taxes"`
	}{
		Taxes: results,
	}
	fmt.Println(response)

	return c.JSON(http.StatusOK, response)
}

func getTaxHandler(c echo.Context) error {
	fmt.Print("tax : % #v\n", responseTax)
	return c.JSON(http.StatusOK, responseTax)
}

type UpdateKReceiptDeductionRequest struct {
	KReceiptDeduction float64 `json:"amount"`
}

func setDeductionsHandler(c echo.Context, config *config.Config) error {

	var req UpdateKReceiptDeductionRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if req.KReceiptDeduction < 10000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid personal deduction amount"})
	}

	if req.KReceiptDeduction > 100000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Admin can only set personal deduction up to 100,000"})

	}
	config.PersonalDeduction = req.KReceiptDeduction

	data, err := json.Marshal(config)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	err = os.WriteFile("config/config.json", data, 0644)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"personalDeduction": req.KReceiptDeduction})

}

func setKreceiptDeductionsHandler(c echo.Context, config *config.Config) error {

	var req UpdateKReceiptDeductionRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if req.KReceiptDeduction < 10000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid personal deduction amount"})
	}

	if req.KReceiptDeduction > 100000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Admin can only set personal deduction up to 100,000"})

	}
	config.KReceiptDeduction = req.KReceiptDeduction

	data, err := json.Marshal(config)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	err = os.WriteFile("config/config.json", data, 0644)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"kReceipt": req.KReceiptDeduction})

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

	// fmt.Println(os.Getenv("ADMIN_USERNAME"))

	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if username == adminUsername && password == adminPassword {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))

	adminGroup := e.Group("/admin")

	adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == adminUsername && password == adminPassword {
			return true, nil
		}
		return false, nil
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	adminGroup.POST("/deductions/personal", func(c echo.Context) error {
		return setDeductionsHandler(c, &config.Config{})
	})

	adminGroup.POST("/deductions/k-receipt", func(c echo.Context) error {
		return setKreceiptDeductionsHandler(c, &config.Config{})
	})

	e.POST("/tax/calculation", createTaxHandler)
	e.POST("/tax/calculations/upload-csv", upload)
	e.GET("/tax/calculation", getTaxHandler)

	log.Fatal(e.Start(":8080"))
}
