package calculation

import (
	"errors"
	"fmt"
	"math"
)

type Allowance struct {
	AllowanceType string
	Amount        float64
}

type TaxLevel struct {
	Level  string
	Amount float64
}

func CalculateTax(totalIncome float64, wht float64, allowances []Allowance) (float64, error) {
	// extract allowance
	var donation float64
	var kReceipt float64

	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			donation = allowance.Amount

		} else if allowance.AllowanceType == "k-receipt" {
			kReceipt = allowance.Amount
		}
	}

	PERSONAL_ALLOWANCE := 60000.0

	if kReceipt < 0.0 {
		return 0.0, errors.New("k-receipt cannot be negative")
	}

	if kReceipt >= 50000.0 {
		kReceipt = 50000.0
	}

	if donation >= 100000.0 {
		donation = 100000.0
	}

	if wht < 0.0 {
		return 0.0, errors.New("wht cannot be negative")
	}

	incomeAfterAllowance := totalIncome - PERSONAL_ALLOWANCE - donation - kReceipt
	fmt.Println("incomeAfterAllowance")
	fmt.Println(incomeAfterAllowance)

	// calculate tax
	if (incomeAfterAllowance) <= 150000.0 {
		return 0.0, nil
	}

	if incomeAfterAllowance > 150000.0 && incomeAfterAllowance <= 500000.0 {
		fmt.Print("--level 1--")
		incomeAfterAllowanceStep1 := incomeAfterAllowance - 150000.0

		unpaidTax := (incomeAfterAllowanceStep1 * 0.1) - wht
		fmt.Println(unpaidTax)
		roundedTax := math.Trunc(unpaidTax*1e10) / 1e10
		return roundedTax, nil
	}

	if incomeAfterAllowance > 500000.0 && incomeAfterAllowance <= 1000000.0 {
		fmt.Print("--level 2--")
		incomeAfterAllowanceStep2 := incomeAfterAllowance - 300000.0
		unpaidTax := (incomeAfterAllowanceStep2 * 0.2) - wht
		roundedTax := math.Trunc(unpaidTax*1e10) / 1e10
		return roundedTax, nil
	}

	return 0.0, nil
}
