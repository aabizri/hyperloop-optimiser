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
	portFlag = flag.Uint("p", 8080, "Port to listen to")    // Indicates the port to listen to
	pathFlag = flag.String("path", "/", "Path to endpoint") // Indicates the path of endpoint
)

func init() {
	flag.Parse()
}

/*
Input is what the requester sends, it is the equivalent of that JSON object:
	{
		"citiesCount" : 4,
		"costOffers" : [
			{
				"from" : 0,
				"to" : 1,
				"price" : 6
			},{
				"from" : 1,
				"to" : 2,
				"price" : 10
			}
		]
	}
*/
type Input struct {
	CitiesCount uint        `json:"citiesCount"`
	CostOffers  []CostOffer `json:"costOffers"`
}

// Validate validates the input
func (input *Input) Validate() error {
	if input.CitiesCount == 0 {
		return fmt.Errorf("Cities count isn't defined")
	}
	return nil
}

/*
CostOffer is a build offer for a hyperloop line, it is the equivalent of that JSON object:
	{
		"from" : 0,
		"to" : 1,
		"price" : 6
	}
*/
type CostOffer struct {
	From  uint `json:"from"`
	To    uint `json:"to"`
	Price uint `json:"price"`
}

// Validate validates a CostOffer
func (co CostOffer) Validate() error {
	if co.From == co.To {
		return fmt.Errorf("Origin and destination of line are the same !")
	}
	if co.Price == 0 {
		return fmt.Errorf("Price is 0")
	}
	return nil
}

/*
Output is the structure to be json encoded and sent back to the requester.
It is the equivalent of that JSON object:
	{
		"feasible" : true,
		"totalCost" : 15,
		"depotId" : 3,
		"recommendedOffers" : [
		{
			"from" : 0,
			"to" : 1,
			"price" : 6
		},{
			"from" : 3,
			"to" : 0,
			"price" : 1
		}
		]
	}
*/
type Output struct {
	Feasible          bool        `json:"feasible"`
	TotalCost         uint        `json:"totalCost,omitempty"`
	DepotID           uint        `json:"depotId,omitempty"`
	RecommendedOffers []CostOffer `json:"recommendedOffers,omitempty"`
}

// Validate validates an Output
func (o *Output) Validate() error {
	if len(o.RecommendedOffers) == 0 {
		return fmt.Errorf("No recommended offers given")
	}
	return nil
}

// Handler handles http requests
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
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, msg)
		return
	}

	// Now let's build a graph with that
	logger.Print("Building graph...")
	graph, err := input.buildGraph()
	if err != nil {
		msg := fmt.Sprintf("Error: buildGraph: %v", err)
		logger.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, msg)
		return
	}
	logger.Print("DONE\n")

	// Call MSA
	logger.Print("Calling MSA...")
	feasible, minimal, root, err := msa.MSAAllRoots(graph)
	if err != nil {
		msg := fmt.Sprintf("Error: MSAAllRoots: %v", err)
		logger.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, msg)
		return
	}
	logger.Print("DONE\n")

	// Build output
	logger.Print("Building output...")
	var id uint64
	if root != nil {
		id, err = strconv.ParseUint(root.String(), 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Error: converting uint: %v", err)
			logger.Print(msg)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, msg)
			return
		}
	}

	// Format the graph to an output
	output, err := format(minimal, uint(id), feasible)
	if err != nil {
		msg := fmt.Sprintf("Error: format: %v", err)
		logger.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, msg)
		return
	}
	logger.Print("DONE\n")

	// Validate the output
	// An error isn't considered a complete failure
	err = output.Validate()
	if err != nil {
		logger.Printf("Output is invalid ! Error : %v", err)
	}

	// Encode it
	logger.Print("Encoding..")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t") // Print it pretty
	err = enc.Encode(output)
	if err != nil {
		msg := fmt.Sprintf("Error: Json encoding: %v", err)
		logger.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, msg)
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
