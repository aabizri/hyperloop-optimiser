package main

import (
	"testing"
)

const iterations int = 360

// Test listCities under a diversity of inputs
// BUG: As helperGenInput doesn't produce cost offers including all cities, it currently fails
func TestListCities_Rand(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	for i := 0; i <= iterations; i++ {
		input := helperGenInput(0, 32, 0, 64, 120, 360)
		if _, err := input.listCities(); err != nil {
			t.Error(err)
		}
	}
}

// Test listCities under a known input
func TestListCities_Known(t *testing.T) {
	// Known cities list
	knownCitiesList := []uint{0, 1, 2, 3}
	funcCitiesList, err := knownInput.listCities()
	if err != nil {
		t.Fatal(err)
	}

	// Compare one to one
	if len(knownCitiesList) != len(funcCitiesList) {
		t.Fatalf("Length of known cities list (%d) differs from that of func (%d)", len(knownCitiesList), len(funcCitiesList))
	}
}
