package main

import (
	"fmt"
	"github.com/gyuho/goraph"
	"strconv"
)

func (input *Input) buildGraph() (goraph.Graph, error) {

	// Create a new graph
	graph := goraph.NewGraph()

	// List all the unique cities
	var cities []uint = make([]uint, input.CitiesCount)
	for _, co := range input.CostOffers {
		// Sanity check
		if co.From == co.To {
			return nil, fmt.Errorf("From (%d) is equal to to (%d)", co.From, co.To)
		}
		var (
			fromIsUnique bool = true
			toIsUnique   bool = true
		)
		// Iterate through all the cities already listed
		for _, city := range cities {
			if co.From == city {
				fromIsUnique = false
			}
			if co.To == city {
				toIsUnique = false
			}
			if fromIsUnique && toIsUnique {
				break
			}
		}

		// If it is unique, add it to the list
		if fromIsUnique {
			cities = append(cities, co.From)
		}
		if toIsUnique {
			cities = append(cities, co.To)
		}
	}

	// Add the nodes
	for _, city := range cities {
		cityStr := strconv.FormatUint(uint64(city), 10)
		graphNode := goraph.NewNode(cityStr)
		ok := graph.AddNode(graphNode)
		if ok != true {
			logger.Printf("Tried to add node %s to the graph but it appears it already exists", cityStr)
		}
	}

	// Now let's deal with the edges
	for _, co := range input.CostOffers {
		err := graph.AddEdge(goraph.NewNode(strconv.FormatUint(uint64(co.From), 10)).ID(), goraph.NewNode(strconv.FormatUint(uint64(co.To), 10)).ID(), float64(co.Price))
		if err != nil {
			return graph, err
		}
	}

	// Return the graph
	return graph, nil
}
