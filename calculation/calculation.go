package calculation

import (
	"errors"
	"fmt"
)

type Allowance struct {
	AllowanceType string
	Amount        float64
}

func CalculateTax(totalIncome float64, wht float64, allowances []Allowance) (float64, error) {
	// extract allowance

	PERSONAL_ALLOWANCE := 60000.0

	if wht < 0.0 {
		return 0.0, errors.New("wht cannot be negative")
	}

	incomeAfterAllowance := totalIncome - PERSONAL_ALLOWANCE

	// calculate tax
	if (incomeAfterAllowance) <= 150000.0 {
		return 0.0, nil
	}

	incomeAfterAllowanceStep1 := incomeAfterAllowance - 150000.0

	if incomeAfterAllowanceStep1 > 150000.0 && incomeAfterAllowanceStep1 <= 300000.0 {
		fmt.Println(incomeAfterAllowanceStep1)

		unpaidTax := (incomeAfterAllowanceStep1 * 0.1) - wht
		fmt.Println(unpaidTax)
		return unpaidTax, nil
	}

	return incomeAfterAllowance, nil
}
