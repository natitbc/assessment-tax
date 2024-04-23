package calculation

import (
	"testing"
)

func TestCalculation(t *testing.T) {

	_, err := CalculateTax(150000.0, -1.0, []Allowance{})
	if err == nil {
		t.Errorf("Expected error for negative wht, got nil")
	}

}
