package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	// Create a writer
	buf := new(bytes.Buffer)

	// Create & setup the encoder
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "\t") // Print it pretty

	// Encode
	err := enc.Encode(knownInput)

	// If there's a problem, panic
	if err != nil {
		panic(fmt.Sprintf("Panic while encoding knownInputJSON ! %v", err))
	}

	// DEBUG
	knownInputJSON = buf.Bytes()
}

// knownInput is a known input
var knownInput *Input = &Input{
	CitiesCount: 4,
	CostOffers: []CostOffer{
		CostOffer{0, 1, 6},
		CostOffer{1, 2, 10},
		CostOffer{2, 1, 10},
		CostOffer{1, 3, 12},
		CostOffer{3, 2, 8},
		CostOffer{3, 0, 1},
	},
}
var knownInputJSON []byte

var knownOutput *Output = &Output{
	Feasible:  true,
	TotalCost: 15,
	DepotID:   3,
	RecommendedOffers: []CostOffer{
		CostOffer{0, 1, 6},
		CostOffer{3, 0, 1},
		CostOffer{3, 2, 8},
	},
}

// For seed
var seedIncrement int64 = 0

// helperGenInput generates a random Input
// BUG: returned input may not have each cities included in at least an edge
func helperGenInput(minCitiesCount, maxCitiesCount, minEdgesCount, maxEdgesCount, minPrice, maxPrice uint) *Input {
	// Seed the PNG
	seed := time.Now().Unix() + seedIncrement
	rand.Seed(seed)
	seedIncrement++

	// Get the amount of cities
	citiesCountIntervalLength := int(maxCitiesCount - minCitiesCount + 1)
	citiesCount := uint(rand.Intn(citiesCountIntervalLength)) + minCitiesCount

	// Get the amount of edges
	edgesCountIntervalLength := int(maxEdgesCount - minEdgesCount + 1)
	edgesCount := uint(rand.Intn(edgesCountIntervalLength)) + minEdgesCount

	// Create the price interval
	priceIntervalLength := int(maxPrice - minPrice + 1)

	// Create the input
	var input *Input = &Input{CitiesCount: citiesCount}

	// If there are no edges or no more than one city then we can return one
	// Note: If there is only one city, no edge can be created
	// Note: If there are no edges, the slice won't get populated so why bother ?
	if citiesCount <= 1 || edgesCount == 0 {
		return input
	}

	// Create the cities
	var cities []uint = make([]uint, citiesCount)
	for i := uint(0); i <= citiesCount-1; i++ {
		cities[i] = i
	}
	if uint(len(cities)) != citiesCount {
		panic("Problem while creating cities list")
	}

	// Create the cost offers
	var costOffers []CostOffer = make([]CostOffer, edgesCount)
	input.CostOffers = costOffers

	// Populate the cost offers
	// TODO: Each city must be touched by an edge
	for i := uint(0); i <= edgesCount-1; i++ {
		// Re-seed
		rand.Seed(seed + int64(i))

		var co CostOffer = CostOffer{}

		// Select the origin
		j := uint(rand.Intn(int(citiesCount - 1)))
		co.From = cities[j]

		// Select the destination
		var k uint
		for attempt := 0; true; attempt++ {
			// Reseed
			rand.Seed(seed + int64(attempt))

			// Generate a k
			k = uint(rand.Intn(int(citiesCount)))

			// If it is different then we found a winner
			if k != j {
				break
			}

			// If it's been that long it means there's a problem here
			if attempt > 12 {
				panic("WTF")
			}
		}
		co.To = cities[k]

		// Select a price
		co.Price = uint(rand.Intn(int(priceIntervalLength))) + minPrice

		// Add it to the slice
		costOffers[i] = co
	}

	return input
}

func TestHelperGenInput(t *testing.T) {
	var (
		minCitiesCount uint = 0
		maxCitiesCount uint = 12
		minEdgesCount  uint = 0
		maxEdgesCount  uint = 24
		minPrice       uint = 120
		maxPrice       uint = 720
	)
	in := helperGenInput(minCitiesCount, maxCitiesCount, minEdgesCount, maxEdgesCount, minPrice, maxPrice)
	t.Logf("Generated input: %v", in)
	if in.CitiesCount < minCitiesCount || in.CitiesCount > maxCitiesCount {
		t.Errorf("helperGenInput failed because CitiesCount is out of bonds (%d is not in [%d;%d])", in.CitiesCount, minCitiesCount, maxCitiesCount)
	}
	lenCostOffers := uint(len(in.CostOffers))
	if lenCostOffers < minEdgesCount || lenCostOffers > maxEdgesCount {
		t.Errorf("helperGenInput failed because the amount of edges is out of bounds (%d is not in [%d;%d])", lenCostOffers, minEdgesCount, maxEdgesCount)
	}
}
