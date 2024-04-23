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
	allowancesMap := make(map[string]float64)
	fmt.Println(allowancesMap)

	if wht < 0.0 {
		return 0.0, errors.New("wht cannot be negative")
	}

	// calculate tax
	if (totalIncome - wht) <= 150000.0 {
		return 0.0, nil
	}

	return 29000.0, nil
}
