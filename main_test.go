package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Known(t *testing.T) {
	// Use the json-encoded Input
	reader := bytes.NewReader(knownInputJSON)

	// Create a request to pass to our handler
	req, err := http.NewRequest("POST", *pathFlag, reader)
	if err != nil {
		t.Fatalf("Request creation failed with error: %v", err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	// Serve http
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Log the response body
	t.Logf("Replied with\n%v\n", rr.Body)

	// Decode the body
	var output *Output = &Output{}
	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(output)
	if err != nil {
		t.Fatalf("Json decoding failed with error: %v", err)
	}

	// Validate the output
	err = output.Validate()
	if err != nil {
		t.Fatalf("Output isn't valid!")
	}

	// Compare the known output with this output
	ok := sameOutput(knownOutput, output)
	if !ok {
		t.Fatalf("Output isn't what is expected !")
	}
}
