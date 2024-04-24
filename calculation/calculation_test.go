package calculation

import (
	"math"
	"testing"
)

func TestCalculation(t *testing.T) {

	const epsilon = 0.001 // Adjust tolerance as needed

	_, _, err := CalculateTax(150000.0, -1.0, []Allowance{})
	if err == nil {
		t.Errorf("Expected error for negative wht, got nil")
	}

	// test got tax 29000 when input 500000, 0, []
	want := 29000.0
	got, _, _ := CalculateTax(500000.0, 0.0, []Allowance{})

	if got != want {
		t.Errorf("Expected %f, got %f", want, got)
	}

	// test got tax 0 when input 150000, 0, []
	want = 0.0
	got, _, _ = CalculateTax(150000.0, 0.0, []Allowance{})
	if got != want {
		t.Errorf("Expected %f, got %f", want, got)
	}

	// test got tax 4000.0 when input 500000.0, 25000.0, []

	want = 4000.0
	got, _, _ = CalculateTax(500000.0, 25000.0, []Allowance{})
	if math.Abs(got-want) > epsilon {
		t.Errorf("Expected %f, got %f", want, got)
	}

	want = 19000.0
	got, _, _ = CalculateTax(500000.0, 0.0, []Allowance{
		{AllowanceType: "donation", Amount: 200000.0},
	})
	if math.Abs(got-want) > epsilon {
		t.Errorf("Expected %f, got %f", want, got)
	}

	want = 14000.0
	got, _, _ = CalculateTax(500000.0, 0.0, []Allowance{
		{AllowanceType: "k-receipt", Amount: 200000.0},
		{AllowanceType: "donation", Amount: 100000.0},
	})
	if math.Abs(got-want) > epsilon {
		t.Errorf("Expected %f, got %f", want, got)
	}

	expectedTaxLevels := []TaxLevel{
		{level: "0-150,000", tax: 0.0},
		{level: "150,001-500,000", tax: 19000.0},
		{level: "500,001-1,000,000", tax: 0.0},
		{level: "1,000,001-2,000,000", tax: 0.0},
		{level: "2,000,001 ขึ้นไป", tax: 0.0},
	}

	_, actualTaxLevels, _ := CalculateTax(500000.0, 0.0, []Allowance{
		{AllowanceType: "donation", Amount: 100000.0},
	})

	if len(expectedTaxLevels) != len(actualTaxLevels) {
		t.Errorf("Expected %d tax levels, got %d", len(expectedTaxLevels), len(actualTaxLevels))
	}

	for i := 0; i < len(expectedTaxLevels); i++ {
		if expectedTaxLevels[i].level != actualTaxLevels[i].level {
			t.Errorf("Expected level %s, got %s", expectedTaxLevels[i].level, actualTaxLevels[i].level)
		}
		if math.Abs(expectedTaxLevels[i].tax-actualTaxLevels[i].tax) > epsilon {
			t.Errorf("Expected tax %f, got %f", expectedTaxLevels[i].tax, actualTaxLevels[i].tax)
		}
	}

}
