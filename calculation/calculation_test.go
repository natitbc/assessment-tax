package calculation

import (
	"testing"
)

func TestCalculation(t *testing.T) {

	_, err := CalculateTax(150000.0, -1.0, []Allowance{})
	if err == nil {
		t.Errorf("Expected error for negative wht, got nil")
	}

	// test got tax 29000 when input 500000, 0, []
	want := 29000.0
	got, _ := CalculateTax(500000.0, 0.0, []Allowance{})

	if got != want {
		t.Errorf("Expected %f, got %f", want, got)
	}

}
