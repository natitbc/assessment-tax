package calculation

import "errors"

type Allowance struct {
	AllowanceType string
	Amount        float64
}

func CalculateTax(totalIncome float64, wht float64, allowances []Allowance) (float64, error) {
	if wht < 0.0 {
		return 0.0, errors.New("wht cannot be negative")
	}

	return 0.0, nil
}
