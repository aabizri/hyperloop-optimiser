// TODO: Add feasability
package main

import (
	"encoding/json"
	"fmt"
	"github.com/aabizri/msa"
	"github.com/gyuho/goraph"
	"log"
	"net/http"
	"strconv"
)

type Input struct {
	CitiesCount uint        `json:"citiesCount"`
	CostOffers  []CostOffer `json:"costOffers"`
}

type CostOffer struct {
	From  uint `json:"from"`
	To    uint `json:"to"`
	Price uint `json:"price"`
}

type Output struct {
	Feasible          bool        `json:"feasible"`
	TotalCost         uint        `json:"totalCost"`
	DepotID           uint        `json:"depotId"`
	RecommendedOffers []CostOffer `json:"recommendedOffers"`
}

type ValidationError struct {
	in               Input
	CitiesCountIsNil bool
}

func (err ValidationError) Error() string {
	return "Validation error : no cities indicated !"
}

func (input *Input) Validate() error {
	var err ValidationError
	if input.CitiesCount == 0 {
		err.CitiesCountIsNil = true
		return err
	}
	return nil
}

func (input *Input) buildGraph() (goraph.Graph, error) {

	// Create a new graph
	graph := goraph.NewGraph()

	// List all the unique cities
	var cities []uint = make([]uint, input.CitiesCount)
	for _, co := range input.CostOffers {
		// Sanity check
		if co.From == co.To {
			panic(fmt.Sprintf("From (%d) is equal to to (%d), PANIC", co.From, co.To))
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
			log.Printf("Tried to add node %s to the graph but it appears it already exists !", cityStr)
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

func createOutput(graph goraph.Graph, rootID uint, feasible bool) (*Output, error) {
	// Create a new output data structure
	var output *Output = &Output{}

	// Check feasability
	if !feasible {
		output.Feasible = false
		return output, nil
	} else {
		output.Feasible = true
	}

	// Get total cost
	totalWeight, err := msa.TotalWeight(graph)
	if err != nil {
		return output, err
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
			log.Printf("Error while converting graph \"source\" id %s to output \"from\" id: %v", edge.Source().ID().String(), err)
			return output, err
		}
		to, err := strconv.ParseUint(edge.Target().ID().String(), 10, 64)
		if err != nil {
			log.Printf("Error while converting graph \"target\" id %s to output \"to\" id: %v", edge.Target().ID().String(), err)
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

func Handler(w http.ResponseWriter, r *http.Request) {
	// Log
	log.Printf("Received request:\n%v\n", r.Body)

	// First let's parse the input
	log.Print("Decoding...")
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var input *Input = &Input{}
	err := dec.Decode(input)
	if err != nil {
		msg := fmt.Sprintf("Error: while decoding json: %v", err)
		log.Print(msg)
		fmt.Fprint(w, msg)
		return
	}
	log.Print("DONE\n")
	log.Print(*input)

	// Let's validate it
	err = input.Validate()
	if err != nil {
		msg := fmt.Sprintf("Error: Validation error: %v", err)
		log.Print(msg)
		fmt.Fprint(w, msg)
		return
	}

	// Now let's build a graph with that
	log.Print("Building graph...")
	graph, err := input.buildGraph()
	if err != nil {
		fmt.Fprintf(w, "Error: buildGraph: %v", err)
		return
	}
	log.Print("DONE\n")

	// Call MSA
	log.Print("Calling MSA...")
	minimal, root, err := msa.MSAAllRoots(graph)
	if err != nil {
		fmt.Fprintf(w, "Error: MSAAllRoots: %v", err)
		return
	}
	log.Print("DONE\n")

	// Build output
	log.Print("Building output...")
	id, err := strconv.ParseUint(root.String(), 10, 64)
	if err != nil {
		log.Printf("Error: %v", err)
		fmt.Fprintf(w, "Error: converting uint: %v", err)
		return
	}
	output, err := createOutput(minimal, uint(id), true)
	if err != nil {
		fmt.Fprintf(w, "Error: createOutput: %v", err)
		return
	}
	log.Print("DONE\n")

	// Encode it
	log.Print("Encoding..")
	enc := json.NewEncoder(w)
	err = enc.Encode(output)
	if err != nil {
		fmt.Fprintf(w, "Error: json encode: %v", err)
		return
	}
	log.Print("DONE\n")

	// Return
	return
}

func main() {
	http.HandleFunc("/", Handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
