package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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

type TaxData struct {
	TotalIncome float64
	Wht         float64
	Allowances  []calculation.Allowance `json:"allowances"`
}

type TaxResult struct {
	Tax      float64
	TaxLevel []calculation.TaxLevel
}

func CalculateTaxes(data []TaxData) ([]TaxResult, error) {
	var results []TaxResult

	for _, entry := range data {
		tax, taxLevel, err := calculation.CalculateTax(entry.TotalIncome, entry.Wht, entry.Allowances)
		if err != nil {
			return nil, err
		}
		result := TaxResult{
			Tax:      tax,
			TaxLevel: taxLevel,
		}
		results = append(results, result)
	}
	return results, nil
}

func parseCSV(data []byte, taxData *[]TaxData) error {
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row (assuming the first row contains headers)
	headers := records[0]
	for _, record := range records[1:] {
		if len(record) != len(headers) {
			return fmt.Errorf("invalid record length: expected %d, got %d", len(headers), len(record))
		}

		// Convert record values to TaxData struct
		var rowData TaxData
		rowData.TotalIncome, _ = strconv.ParseFloat(record[0], 64)
		rowData.Wht, _ = strconv.ParseFloat(record[1], 64)

		// Parse allowances (assuming columns 2 onwards contain allowances)
		rowData.Allowances = make([]calculation.Allowance, 0, len(record)-2) // Pre-allocate slice for allowances
		for i := 2; i < len(record); i++ {
			allowanceType := record[i]
			allowanceAmount, err := strconv.ParseFloat(record[i+1], 64) // Assuming allowance amount follows allowance type
			if err != nil {
				return fmt.Errorf("invalid allowance format: %w", err)
			}
			rowData.Allowances = append(rowData.Allowances, calculation.Allowance{AllowanceType: allowanceType, Amount: allowanceAmount})
			i++ // Skip allowance amount after parsing type and amount
		}

		*taxData = append(*taxData, rowData)
	}

	return nil
}

// func multiTaxHandler(c echo.Context) error {

// 	file, err := c.FormFile("taxFile")
// 	if err != nil {
// 		if err == http.ErrMissingFile {
// 			return c.JSON(http.StatusBadRequest, Err{Message: "No file uploaded"})
// 		}
// 		fmt.Println("Error retrieving uploaded file:", err) // Log error details
// 		return c.JSON(http.StatusInternalServerError, Err{Message: "Error retrieving uploaded file"})
// 	}

// 	// If file is uploaded successfully, log some details
// 	fmt.Println("Uploaded filename:", file.Filename)
// 	fmt.Println("Uploaded file size:", file.Size)

// 	fmt.Print(file)

// 	// Opean uploaded file
// 	reader, err := file.Open()
// 	fmt.Print("reader :", reader)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, Err{Message: "Error opening uploaded file"})
// 	}
// 	defer reader.Close()

// 	// read csv data
// 	var data bytes.Buffer
// 	_, err = io.Copy(&data, reader)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, Err{Message: "Error reading uploaded file"})
// 	}

// 	// Parse CSV
// 	var taxData []TaxData
// 	err = parseCSV(data.Bytes(), &taxData)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, Err{Message: "Error parsing uploaded file"})
// 	}

// 	// Process tax data
// 	results, err := CalculateTaxes([]TaxData{taxData[0]})
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, Err{Message: "Error calculating tax"})
// 	}

//		return c.JSON(http.StatusOK, results)
//	}
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

	fd, error := os.Open("taxes.csv")

	if error != nil {
		fmt.Println(error)
	}

	fmt.Println("Successfully opened the CSV file")
	defer fd.Close()

	// read csv data
	fileReader := csv.NewReader(fd)
	records, error := fileReader.ReadAll()

	if error != nil {
		fmt.Println(error)
	}

	fmt.Println(records)
	// 	// Parse CSV
	// 	var taxData []TaxData
	// 	err = parseCSV(data.Bytes(), &taxData)
	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, Err{Message: "Error parsing uploaded file"})
	// 	}

	// 	// Process tax data
	// 	results, err := CalculateTaxes([]TaxData{taxData[0]})
	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, Err{Message: "Error calculating tax"})
	// 	}

	return c.JSON(http.StatusOK, records)
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

	return c.JSON(http.StatusOK, map[string]interface{}{"personalDeduction": req.PersonalDeduction})

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

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == adminUsername && password == adminPassword {
			return true, nil
		}
		return false, nil
	}))

	adminGroup := e.Group("/admin")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	adminGroup.POST("/deductions/personal", func(c echo.Context) error {
		return setDeductionsHandler(c, &calculation.Config{})
	})

	e.POST("/tax/calculation", createTaxHandler)
	e.POST("/tax/calculations/upload-csv", upload)
	e.GET("/tax/calculation", getTaxHandler)

	log.Fatal(e.Start(":8080"))
}
