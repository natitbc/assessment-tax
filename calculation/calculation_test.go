package calculation

import (
	"testing"
)

func TestCalculation(t *testing.T) {
	want := 0.0

	if got := CalculateTax(150000.0, 0.0, []Allowance{}); got != want {
		t.Errorf("CalculateTax() = %v, want %v", got, want)
	}

}
