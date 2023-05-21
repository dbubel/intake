package intake

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAddToContext(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	key := "exampleKey"
	value := "exampleValue"

	err := AddToContext(r, key, value)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that the value was added to the context
	data, ok := r.Context().Value(key).([]byte)
	if !ok {
		t.Errorf("Error retrieving data from context")
	}

	var resultValue string
	json.Unmarshal(data, &resultValue)
	if resultValue != value {
		t.Errorf("Expected %s but got %s", value, resultValue)
	}
}

func TestFromContext(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	key := "exampleKey"
	value := "exampleValue"

	err := AddToContext(r, key, value)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var resultValue string
	err = FromContext(r, key, &resultValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that the correct value was retrieved from the context
	if resultValue != value {
		t.Errorf("Expected %s but got %s", value, resultValue)
	}
}
