package calculation

import (
	"testing"
)

func TestCalculation(t *testing.T) {
	// Test case 1: Zero income, zero deductions, zero tax
	want := 0.0
	if got := CalculateTax(0.0, 0.0, []Allowance{}); got != want {
		t.Errorf("CalculateTax() = %v, want %v", got, want)
	}

	// Test case 2: Positive income, zero deductions, some tax (replace 100.0 with expected tax for this income)
	want = 0.0 // Replace with expected tax amount
	if got := CalculateTax(150000.0, 0.0, []Allowance{}); got != want {
		t.Errorf("CalculateTax() = %v, want %v", got, want)
	}

}
