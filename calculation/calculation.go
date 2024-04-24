package calculation

import (
	"errors"
	"math"
)

type Allowance struct {
	AllowanceType string
	Amount        float64
}

func CalculateTax(totalIncome float64, wht float64, allowances []Allowance) (float64, error) {
	// extract allowance
	var donation float64

	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			donation = allowance.Amount
			break
		}
	}

	PERSONAL_ALLOWANCE := 60000.0

	if donation > 100000.0 {
		donation = 100000.0
	}

	if wht < 0.0 {
		return 0.0, errors.New("wht cannot be negative")
	}

	incomeAfterAllowance := totalIncome - PERSONAL_ALLOWANCE - donation

	// calculate tax
	if (incomeAfterAllowance) <= 150000.0 {
		return 0.0, nil
	}

	incomeAfterAllowanceStep1 := incomeAfterAllowance - 150000.0

	if incomeAfterAllowanceStep1 > 150000.0 && incomeAfterAllowanceStep1 <= 300000.0 {

		unpaidTax := (incomeAfterAllowanceStep1 * 0.1) - wht
		roundedTax := math.Trunc(unpaidTax*1e10) / 1e10
		return roundedTax, nil
	}

	return 0.0, nil
}
