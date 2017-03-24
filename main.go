package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aabizri/msa"
	"net/http"
	"strconv"
)

var (
	portFlag = flag.Uint("p", 8080, "Port to listen to")
	pathFlag = flag.String("path", "/", "Path to endpoint")
)

func init() {
	flag.Parse()
}

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
	TotalCost         uint        `json:"totalCost,omitempty"`
	DepotID           uint        `json:"depotId,omitempty"`
	RecommendedOffers []CostOffer `json:"recommendedOffers,omitempty"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Log the request
	err := logRequest(r)
	if err != nil {
		logger.Printf("Handler: request logging failed: %v", err)
	}

	// First let's parse the input
	logger.Print("Decoding...")
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var input *Input = &Input{}
	err = dec.Decode(input)
	if err != nil {
		msg := fmt.Sprintf("Error: while decoding json: %v", err)
		logger.Print(msg)
		fmt.Fprint(w, msg)
		return
	}
	logger.Print("DONE\n")
	logger.Print(*input)

	// Let's validate it
	err = input.Validate()
	if err != nil {
		msg := fmt.Sprintf("Error: Validation error: %v", err)
		logger.Print(msg)
		fmt.Fprint(w, msg)
		return
	}

	// Now let's build a graph with that
	logger.Print("Building graph...")
	graph, err := input.buildGraph()
	if err != nil {
		fmt.Fprintf(w, "Error: buildGraph: %v", err)
		return
	}
	logger.Print("DONE\n")

	// Call MSA
	logger.Print("Calling MSA...")
	feasible, minimal, root, err := msa.MSAAllRoots(graph)
	if err != nil {
		fmt.Fprintf(w, "Error: MSAAllRoots: %v", err)
		return
	}
	logger.Print("DONE\n")

	// Build output
	logger.Print("Building output...")
	var id uint64
	if root != nil {
		id, err = strconv.ParseUint(root.String(), 10, 64)
		if err != nil {
			logger.Printf("Error: %v", err)
			fmt.Fprintf(w, "Error: converting uint: %v", err)
			return
		}
	}

	output, err := format(minimal, uint(id), feasible)
	if err != nil {
		fmt.Fprintf(w, "Error: createOutput: %v", err)
		return
	}
	logger.Print("DONE\n")

	// Encode it
	logger.Print("Encoding..")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t") // Print it pretty
	err = enc.Encode(output)
	if err != nil {
		fmt.Fprintf(w, "Error: json encode: %v", err)
		return
	}
	logger.Print("DONE\n")

	// Return
	return
}

func main() {
	http.HandleFunc(*pathFlag, Handler)

	portStr := ":" + strconv.FormatUint(uint64(*portFlag), 10)
	logger.Fatal(http.ListenAndServe(portStr, nil))
}
