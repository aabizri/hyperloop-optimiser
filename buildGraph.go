package main

import (
	"fmt"
	"github.com/gyuho/goraph"
	"strconv"
)

// listCities lists all the cities implied in an Input
func (input *Input) listCities() ([]uint, error) {
	// List all the unique cities
	var cities []uint = make([]uint, 0)
	for _, co := range input.CostOffers {
		// Sanity check
		if co.From == co.To {
			return nil, fmt.Errorf("listCities: \"from\" (%d) is equal to \"to\" (%d)", co.From, co.To)
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
		}

		// If it is unique, add it to the list
		if fromIsUnique {
			cities = append(cities, co.From)
		}
		if toIsUnique {
			cities = append(cities, co.To)
		}
	}
	return cities, nil
}

// addEdgeToGraph adds an edge to a graph given a Costoffer and a graph to add to (duh)
func (co CostOffer) addEdgeToGraph(g goraph.Graph) error {
	err := g.AddEdge(goraph.NewNode(strconv.FormatUint(uint64(co.From), 10)).ID(), goraph.NewNode(strconv.FormatUint(uint64(co.To), 10)).ID(), float64(co.Price))
	if err != nil {
		return fmt.Errorf("addEdgeToGraph: error: %v", err)
	}
	return err
}

// buildGraph transforms an Input to a graph
func (input *Input) buildGraph() (goraph.Graph, error) {

	// Create a new graph
	graph := goraph.NewGraph()

	// Get the cities
	cities, err := input.listCities()
	if err != nil {
		return graph, fmt.Errorf("buildGraph: error: %v", err)
	}

	// Add the nodes
	for _, city := range cities {
		cityStr := strconv.FormatUint(uint64(city), 10)
		graphNode := goraph.NewNode(cityStr)
		ok := graph.AddNode(graphNode)
		if ok != true {
			logger.Printf("buildGraph: Tried to add node %s to the graph but it appears it already exists", cityStr)
		}
	}

	// Now let's deal with the edges
	for _, co := range input.CostOffers {
		err = co.addEdgeToGraph(graph)
		if err != nil {
			return graph, fmt.Errorf("buildGraph: error: %v", err)
		}
	}

	// Return the graph
	return graph, nil
}
