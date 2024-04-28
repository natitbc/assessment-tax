package calculation

import (
	"fmt"
	"math"
	"testing"
)

func TestCalculation(t *testing.T) {

	const epsilon = 0.001 // Adjust tolerance as needed

	//case1

	_, _, err := CalculateTax(150000.0, -1.0, []Allowance{})
	if err == nil {
		t.Errorf("Expected error for negative wht, got nil")
	}

	// case2
	// test got tax 29000 when input 500000, 0, []
	want := 29000.0
	got, _, _ := CalculateTax(500000.0, 0.0, []Allowance{})

	if got != want {
		t.Errorf("Expected %f, got %f", want, got)
	}

	// case3
	// test got tax 0 when input 150000, 0, []
	want = 0.0
	got, _, _ = CalculateTax(150000.0, 0.0, []Allowance{})
	if got != want {
		t.Errorf("Expected %f, got %f", want, got)
	}

	// case4
	// test got tax 4000.0 when input 500000.0, 25000.0, []
	fmt.Println(">>>>>> Test case 4")
	want = 4000.0
	got, _, _ = CalculateTax(500000.0, 25000.0, []Allowance{})
	if math.Abs(got-want) > epsilon {
		t.Errorf("Expected %f, got %f", want, got)
	}

	// case5
	want = 19000.0
	got, _, _ = CalculateTax(500000.0, 0.0, []Allowance{
		{AllowanceType: "donation", Amount: 200000.0},
	})
	if math.Abs(got-want) > epsilon {
		t.Errorf("Expected %f, got %f", want, got)
	}

	// case6
	want = 14000.0
	got, _, _ = CalculateTax(500000.0, 0.0, []Allowance{
		{AllowanceType: "k-receipt", Amount: 200000.0},
		{AllowanceType: "donation", Amount: 100000.0},
	})
	if math.Abs(got-want) > epsilon {
		t.Errorf("Expected %f, got %f", want, got)
	}

	// case7 check for negative wht
	want = 0.0
	_, _, err = CalculateTax(500000.0, 600000.0, []Allowance{})
	if err == nil {
		t.Errorf("Expected error for negative wht, got nil")
	}

	// case7 check for negative wht
	want = 0.0
	_, _, err = CalculateTax(500000.0, -1.0, []Allowance{})
	if err == nil {
		t.Errorf("Expected error for negative wht, got nil")
	}

	// case7
	expectedTaxLevels := []TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 19000.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}

	_, actualTaxLevels, _ := CalculateTax(500000.0, 0.0, []Allowance{
		{AllowanceType: "donation", Amount: 100000.0},
	})

	if len(expectedTaxLevels) != len(actualTaxLevels) {
		t.Errorf("Expected %d tax levels, got %d", len(expectedTaxLevels), len(actualTaxLevels))
	}

	for i := 0; i < len(expectedTaxLevels); i++ {
		if expectedTaxLevels[i].Level != actualTaxLevels[i].Level {
			t.Errorf("Expected level %s, got %s", expectedTaxLevels[i].Level, actualTaxLevels[i].Level)
		}
		if math.Abs(expectedTaxLevels[i].Tax-actualTaxLevels[i].Tax) > epsilon {
			t.Errorf("Expected tax %f, got %f", expectedTaxLevels[i].Tax, actualTaxLevels[i].Tax)
		}
	}

	// case8
	expectedTaxLevels2 := []TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 14000.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}

	_, actualTaxLevels2, _ := CalculateTax(500000.0, 0.0, []Allowance{
		{AllowanceType: "donation", Amount: 100000.0},
		{AllowanceType: "k-receipt", Amount: 200000.0},
	})

	if len(expectedTaxLevels2) != len(actualTaxLevels2) {
		t.Errorf("Expected %d tax levels, got %d", len(expectedTaxLevels2), len(actualTaxLevels2))
	}

	for i := 0; i < len(expectedTaxLevels2); i++ {
		if expectedTaxLevels2[i].Level != actualTaxLevels2[i].Level {
			t.Errorf("Expected level %s, got %s", expectedTaxLevels2[i].Level, actualTaxLevels2[i].Level)
		}
		if math.Abs(expectedTaxLevels2[i].Tax-actualTaxLevels2[i].Tax) > epsilon {
			t.Errorf("Expected tax %f, got %f", expectedTaxLevels2[i].Tax, actualTaxLevels2[i].Tax)
		}
	}

}
