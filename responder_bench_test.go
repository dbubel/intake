package intake

import (
	"bytes"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Person struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

func BenchmarkRespondJSON(b *testing.B) {
	// Create a sample http.Request
	req, _ := http.NewRequest("GET", "/benchmark", nil)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a new ResponseRecorder for each iteration
		// Create a sample data instance
		data := &SampleData{
			Name:  gofakeit.BeerName(),
			Value: 42,
			Para:  gofakeit.Paragraph(1, 3, 10, ""),
		}
		w := httptest.NewRecorder()

		// Call the RespondJSON function
		RespondJSON(w, req, http.StatusOK, data)

		// Check the status code
		if w.Code != http.StatusOK {
			b.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Check the content type
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			b.Errorf("Expected content type 'application/json', got '%s'", contentType)
		}

		// Unmarshal the response JSON
		var responseData SampleData
		if err := json.Unmarshal(w.Body.Bytes(), &responseData); err != nil {
			b.Errorf("Error unmarshalling response JSON: %v", err)
		}

		// Check the response data
		if responseData.Name != data.Name || responseData.Value != data.Value {
			b.Errorf("Expected response data %+v, got %+v", data, responseData)
		}
	}
}
func BenchmarkRespondJSONSyncPool(b *testing.B) {
	// Create a sample http.Request
	req, _ := http.NewRequest("GET", "/benchmark", nil)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a sample data instance
		data := &SampleData{
			Name:  gofakeit.BeerName(),
			Value: 42,
			Para:  gofakeit.Paragraph(1, 3, 10, ""),
		}

		// Create a new ResponseRecorder for each iteration
		w := httptest.NewRecorder()

		// Call the RespondJSON function
		RespondJSONSyncPool(w, req, http.StatusOK, data)

		// Check the status code
		if w.Code != http.StatusOK {
			b.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Check the content type
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			b.Errorf("Expected content type 'application/json', got '%s'", contentType)
		}

		// Unmarshal the response JSON
		var responseData SampleData
		if err := json.Unmarshal(w.Body.Bytes(), &responseData); err != nil {
			b.Errorf("Error unmarshalling response JSON: %v", err)
		}

		// Check the response data
		if responseData.Name != data.Name || responseData.Value != data.Value {
			b.Errorf("Expected response data %+v, got %+v", data, responseData)
		}
	}
}
func BenchmarkRespondJSONStream(b *testing.B) {
	// Create a sample http.Request
	req, _ := http.NewRequest("GET", "/benchmark", nil)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a sample data instance
		data := &SampleData{
			Name:  gofakeit.BeerName(),
			Value: 42,
			Para:  gofakeit.Paragraph(1, 3, 10, ""),
		}

		// Create a new ResponseRecorder for each iteration
		w := httptest.NewRecorder()

		// Call the RespondJSON function

		RespondJSONStream(w, req, http.StatusOK, data)

		// Check the status code
		if w.Code != http.StatusOK {
			b.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Check the content type
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			b.Errorf("Expected content type 'application/json', got '%s'", contentType)
		}

		// Unmarshal the response JSON
		var responseData SampleData
		if err := json.Unmarshal(w.Body.Bytes(), &responseData); err != nil {
			b.Errorf("Error unmarshalling response JSON: %v", err)
		}

		// Check the response data
		if responseData.Name != data.Name || responseData.Value != data.Value {
			b.Errorf("Expected response data %+v, got %+v", data, responseData)
		}
	}
}
func BenchmarkJSONEncodeSyncPoolDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		person := Person{Name: gofakeit.Animal(), Age: 30, Gender: gofakeit.BeerName()}
		buf := bufferPool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			bufferPool.Put(buf)
		}()
		encoder := json.NewEncoder(buf)
		_ = encoder.Encode(&person)
		_ = buf.Bytes()
	}
}
func BenchmarkJSONEncodeSyncPoolNoDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		person := Person{Name: gofakeit.Animal(), Age: 30, Gender: gofakeit.BeerName()}
		buf := bufferPool.Get().(*bytes.Buffer)
		encoder := json.NewEncoder(buf)
		_ = encoder.Encode(&person)
		_ = buf.Bytes()
		buf.Reset()
		bufferPool.Put(buf)
	}
}
func BenchmarkJSONEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		person := Person{Name: gofakeit.Animal(), Age: 30, Gender: gofakeit.BeerName()}
		buf := bytes.Buffer{}
		encoder := json.NewEncoder(&buf)
		_ = encoder.Encode(&person)
		_ = buf.Bytes()
	}
}
func BenchmarkJSONMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		person := Person{Name: gofakeit.Animal(), Age: 30, Gender: gofakeit.BeerName()}
		buf, _ := json.Marshal(person)
		_ = buf
	}
}

func BenchmarkUnmarshalJSONSync(b *testing.B) {
	type Person struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
	}

	sampleJSON := `{"name":"John Doe","age":30,"gender":"male"}`

	for i := 0; i < b.N; i++ {
		inputReader := strings.NewReader(sampleJSON)
		var person Person
		if err := UnmarshalJSONSync(inputReader, &person); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalJSONStream(b *testing.B) {
	type Person struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
	}

	sampleJSON := `{"name":"John Doe","age":30,"gender":"male"}`

	for i := 0; i < b.N; i++ {
		inputReader := strings.NewReader(sampleJSON)
		var person Person
		if err := UnmarshalJSONStream(inputReader, &person); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	type Person struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
	}

	sampleJSON := `{"name":"John Doe","age":30,"gender":"male"}`

	for i := 0; i < b.N; i++ {
		inputReader := strings.NewReader(sampleJSON)
		var person Person
		if err := UnmarshalJSON(inputReader, &person); err != nil {
			b.Fatal(err)
		}
	}
}
