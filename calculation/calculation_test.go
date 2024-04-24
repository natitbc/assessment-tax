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

}
