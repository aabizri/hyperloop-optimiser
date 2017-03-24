package main

import (
	"fmt"
	"github.com/aabizri/msa"
	"github.com/gyuho/goraph"
	"strconv"
)

// format transforms a graph, root ID and feasability into a ready-to-use *Output
func format(graph goraph.Graph, rootID uint, feasible bool) (*Output, error) {
	// Create a new output data structure
	var output *Output = &Output{}

	// Check feasability
	output.Feasible = feasible
	if !feasible {
		return output, nil
	}

	// Get total cost
	totalWeight, err := msa.TotalWeight(graph)
	if err != nil {
		return output, fmt.Errorf("format: error: %v", err)
	}
	output.TotalCost = uint(totalWeight)

	// Add DepotID
	output.DepotID = rootID

	// Now add the cost offers
	edges, err := msa.GetEdges(graph)
	if err != nil {
		return output, err
	}
	output.RecommendedOffers = make([]CostOffer, len(edges))
	for i, edge := range edges {
		from, err := strconv.ParseUint(edge.Source().ID().String(), 10, 64)
		if err != nil {
			logger.Printf("Error while converting graph \"source\" id %s to output \"from\" id: %v", edge.Source().ID().String(), err)
			return output, err
		}
		to, err := strconv.ParseUint(edge.Target().ID().String(), 10, 64)
		if err != nil {
			logger.Printf("Error while converting graph \"target\" id %s to output \"to\" id: %v", edge.Target().ID().String(), err)
			return output, err
		}
		co := CostOffer{
			From:  uint(from),
			To:    uint(to),
			Price: uint(edge.Weight()),
		}
		output.RecommendedOffers[i] = co
	}

	// Return
	return output, err
}
