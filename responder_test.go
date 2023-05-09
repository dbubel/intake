package intake
//import (
//
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//)
//
//// Sample data structure to be used in the benchmark test
//type SampleData struct {
//	Name  string `json:"name"`
//	Value int    `json:"value"`
//	Para  string `json:"para"`
//}
//
//func TestRespondJSON(t *testing.T) {
//	data := map[string]interface{}{
//		"message": "Hello, world!",
//	}
//
//	handler := func(w http.ResponseWriter, r *http.Request) {
//		_, err := RespondJSON(w, r, http.StatusOK, data)
//		if err != nil {
//			t.Errorf("Expected no error, got: %v", err)
//		}
//	}
//
//	testServer := httptest.NewServer(http.HandlerFunc(handler))
//	defer testServer.Close()
//
//	resp, err := http.Get(testServer.URL)
//	if err != nil {
//		t.Fatalf("Expected no error making request, got: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected status code %d, got: %d", http.StatusOK, resp.StatusCode)
//	}
//
//	var responseData map[string]interface{}
//	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
//		t.Fatalf("Error decoding JSON response: %v", err)
//	}
//	if responseData["message"] != data["message"] {
//		t.Errorf("Expected message %q, got: %q", data["message"], responseData["message"])
//	}
//}
//
//func TestRespondJSONSyncPool(t *testing.T) {
//	data := map[string]interface{}{
//		"message": "Hello, world!",
//	}
//
//	handler := func(w http.ResponseWriter, r *http.Request) {
//		_, err := RespondJSONSyncPool(w, r, http.StatusOK, data)
//		if err != nil {
//			t.Errorf("Expected no error, got: %v", err)
//		}
//	}
//
//	testServer := httptest.NewServer(http.HandlerFunc(handler))
//	defer testServer.Close()
//
//	resp, err := http.Get(testServer.URL)
//	if err != nil {
//		t.Fatalf("Expected no error making request, got: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected status code %d, got: %d", http.StatusOK, resp.StatusCode)
//	}
//
//	var responseData map[string]interface{}
//	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
//		t.Fatalf("Error decoding JSON response: %v", err)
//	}
//	if responseData["message"] != data["message"] {
//		t.Errorf("Expected message %q, got: %q", data["message"], responseData["message"])
//	}
//}
//
//func TestUnmarshalJSONSync(t *testing.T) {
//	type Person struct {
//		Name   string `json:"name"`
//		Age    int    `json:"age"`
//		Gender string `json:"gender"`
//	}
//
//	tests := []struct {
//		name          string
//		inputJSON     string
//		expectedError bool
//		expected      Person
//	}{
//		{
//			name:          "Valid JSON input",
//			inputJSON:     `{"name":"John Doe","age":30,"gender":"male"}`,
//			expectedError: false,
//			expected: Person{
//				Name:   "John Doe",
//				Age:    30,
//				Gender: "male",
//			},
//		},
//		{
//			name:          "Invalid JSON input",
//			inputJSON:     `{"name":"John Doe","age":"30","gender":"male"}`,
//			expectedError: true,
//		},
//		{
//			name:          "Empty JSON input",
//			inputJSON:     `{}`,
//			expectedError: false,
//			expected:      Person{},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			inputReader := strings.NewReader(test.inputJSON)
//			var person Person
//			err := UnmarshalJSONSync(inputReader, &person)
//
//			if test.expectedError && err == nil {
//				t.Error("Expected an error, but didn't get one")
//			} else if !test.expectedError && err != nil {
//				t.Errorf("Unexpected error: %v", err)
//			}
//
//			if !test.expectedError && person != test.expected {
//				t.Errorf("Expected person %v, but got %v", test.expected, person)
//			}
//		})
//	}
//}
//
//func TestUnmarshalJSONStream(t *testing.T) {
//	type Person struct {
//		Name   string `json:"name"`
//		Age    int    `json:"age"`
//		Gender string `json:"gender"`
//	}
//
//	tests := []struct {
//		name          string
//		inputJSON     string
//		expectedError bool
//		expected      Person
//	}{
//		{
//			name:          "Valid JSON input",
//			inputJSON:     `{"name":"John Doe","age":30,"gender":"male"}`,
//			expectedError: false,
//			expected: Person{
//				Name:   "John Doe",
//				Age:    30,
//				Gender: "male",
//			},
//		},
//		{
//			name:          "Invalid JSON input",
//			inputJSON:     `{"name":"John Doe","age":"30","gender":"male"}`,
//			expectedError: true,
//		},
//		{
//			name:          "Empty JSON input",
//			inputJSON:     `{}`,
//			expectedError: false,
//			expected:      Person{},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			inputReader := strings.NewReader(test.inputJSON)
//			var person Person
//			err := UnmarshalJSONStream(inputReader, &person)
//
//			if test.expectedError && err == nil {
//				t.Error("Expected an error, but didn't get one")
//			} else if !test.expectedError && err != nil {
//				t.Errorf("Unexpected error: %v", err)
//			}
//
//			if !test.expectedError && person != test.expected {
//				t.Errorf("Expected person %v, but got %v", test.expected, person)
//			}
//		})
//	}
//}
//
//func TestUnmarshalJSON(t *testing.T) {
//	type Person struct {
//		Name   string `json:"name"`
//		Age    int    `json:"age"`
//		Gender string `json:"gender"`
//	}
//
//	tests := []struct {
//		name          string
//		inputJSON     string
//		expectedError bool
//		expected      Person
//	}{
//		{
//			name:          "Valid JSON input",
//			inputJSON:     `{"name":"John Doe","age":30,"gender":"male"}`,
//			expectedError: false,
//			expected: Person{
//				Name:   "John Doe",
//				Age:    30,
//				Gender: "male",
//			},
//		},
//		{
//			name:          "Invalid JSON input",
//			inputJSON:     `{"name":"John Doe","age":"30","gender":"male"}`,
//			expectedError: true,
//		},
//		{
//			name:          "Empty JSON input",
//			inputJSON:     `{}`,
//			expectedError: false,
//			expected:      Person{},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			inputReader := strings.NewReader(test.inputJSON)
//			var person Person
//			err := UnmarshalJSON(inputReader, &person)
//
//			if test.expectedError && err == nil {
//				t.Error("Expected an error, but didn't get one")
//			} else if !test.expectedError && err != nil {
//				t.Errorf("Unexpected error: %v", err)
//			}
//
//			if !test.expectedError && person != test.expected {
//				t.Errorf("Expected person %v, but got %v", test.expected, person)
//			}
//		})
//	}
//}
