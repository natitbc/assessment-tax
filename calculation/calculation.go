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
	Level string
	Tax   float64
}

func CalculateTax(totalIncome float64, wht float64, allowances []Allowance) (float64, []TaxLevel, error) {
	// extract allowance
	var donation float64
	var kReceipt float64

	TaxLevels := []TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 0.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}

	PERSONAL_ALLOWANCE := 60000.0
	kReceipt = 0.0

	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			donation = allowance.Amount

		} else if allowance.AllowanceType == "k-receipt" {
			kReceipt = allowance.Amount
		} else {
			return 0.0, TaxLevels, errors.New("invalid allowance type")
		}
	}

	if kReceipt < 0.0 {
		return 0.0, TaxLevels, errors.New("k-receipt cannot be negative")
	}

	if kReceipt >= 50000.0 {
		kReceipt = 50000.0
	}

	if donation < 0.0 {
		return 0.0, TaxLevels, errors.New("donation cannot be negative")
	}

	if donation >= 100000.0 {
		donation = 100000.0
	}

	if wht < 0.0 {
		return 0.0, TaxLevels, errors.New("wht cannot be negative")
	}

	if wht > totalIncome {
		return 0.0, TaxLevels, errors.New("wht cannot be greater than total income")
	}
	fmt.Println("------Conditions------")
	fmt.Println("totalIncome: ", totalIncome)
	fmt.Println("wht: ", wht)
	fmt.Println("PERSONAL_ALLOWANCE: ", PERSONAL_ALLOWANCE)
	fmt.Println("kReceipt: ", kReceipt)
	fmt.Println("donation: ", donation)

	incomeAfterAllowance := totalIncome - PERSONAL_ALLOWANCE - donation - kReceipt
	incomeAfterAllowanceStep1 := incomeAfterAllowance - 150000.0
	incomeAfterAllowanceStep2 := incomeAfterAllowanceStep1 - 500000.0
	incomeAfterAllowanceStep3 := incomeAfterAllowanceStep2 - 1000000.0
	incomeAfterAllowanceStep4 := incomeAfterAllowanceStep3 - 2000000.0

	fmt.Println("------------")

	totalTax := 0.0

	// calculate tax
	if (incomeAfterAllowance) <= 150000.0 {
		return 0.0, TaxLevels, nil
	}

	if incomeAfterAllowanceStep1 > 0 {
		unpaidTax := (incomeAfterAllowanceStep1 * 0.1)
		roundedTax := math.Trunc(unpaidTax*1e10) / 1e10
		TaxLevels[1].Tax = roundedTax
		totalTax += roundedTax
		fmt.Println("taxlevel1: ", TaxLevels[1].Tax)
	}

	if incomeAfterAllowanceStep2 > 0 {
		fmt.Println("incomeAfterAllowanceStep2: ", incomeAfterAllowanceStep2)
		unpaidTax := (incomeAfterAllowanceStep2 * 0.15)
		roundedTax := math.Trunc(unpaidTax*1e10) / 1e10
		TaxLevels[2].Tax = roundedTax
		totalTax += roundedTax
		fmt.Println("taxlevel2: ", TaxLevels[2].Tax)
	}

	if incomeAfterAllowanceStep3 > 0 {
		unpaidTax := (incomeAfterAllowanceStep3 * 0.2)
		roundedTax := math.Trunc(unpaidTax*1e10) / 1e10
		TaxLevels[3].Tax = roundedTax
		totalTax += roundedTax
	}

	if incomeAfterAllowanceStep4 > 0 {
		unpaidTax := (incomeAfterAllowanceStep3 * 0.35)
		roundedTax := math.Trunc(unpaidTax*1e10) / 1e10
		TaxLevels[4].Tax = roundedTax
		totalTax += roundedTax
	}

	roundedTax := math.Trunc(totalTax*1e10) / 1e10
	roundedTax = roundedTax - wht

	fmt.Println("+++++++++++")
	fmt.Println("roundedTax: ", roundedTax)
	fmt.Println("totalTax: ", totalTax)
	fmt.Println("wht: ", wht)
	fmt.Println("TaxLevel: ", TaxLevels)

	return roundedTax, TaxLevels, nil

}
