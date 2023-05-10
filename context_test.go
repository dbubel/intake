package intake

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

type TestEmployee struct {
	ID     int
	Name   string
	Salary int
}
//
//
//// TestEncodeGob tests the encodeGob function
//func TestEncodeGob(t *testing.T) {
//	employee := &TestEmployee{
//		ID:     1,
//		Name:   "John Doe",
//		Salary: 50000,
//	}
//
//	_, err := encodeGob(employee)
//	if err != nil {
//		t.Errorf("encodeGob() error = %v", err)
//	}
//}
//
//// TestDecodeGob tests the decodeGob function
//func TestDecodeGob(t *testing.T) {
//	employee := &TestEmployee{
//		ID:     1,
//		Name:   "John Doe",
//		Salary: 50000,
//	}
//
//	buf, err := encodeGob(employee)
//	if err != nil {
//		t.Fatalf("encodeGob() error = %v", err)
//	}
//
//	decodedEmployee := new(TestEmployee)
//	err = decodeGob(buf, decodedEmployee)
//	if err != nil {
//		t.Errorf("decodeGob() error = %v", err)
//	}
//
//	if employee.ID != decodedEmployee.ID ||
//		employee.Name != decodedEmployee.Name ||
//		employee.Salary != decodedEmployee.Salary {
//		t.Errorf("decodeGob() got = %+v, want %+v", decodedEmployee, employee)
//	}
//}
//
//
//
//// BenchmarkEncodeGob benchmarks the encodeGob function
//func BenchmarkEncodeGob(b *testing.B) {
//	employee := &TestEmployee{
//		ID:     1,
//		Name:   "John Doe",
//		Salary: 50000,
//	}
//
//	b.ResetTimer() // resets the timer to ignore the time taken for setup
//
//	for i := 0; i < b.N; i++ {
//		_, err := encodeGob(employee)
//		if err != nil {
//			b.Errorf("encodeGob() error = %v", err)
//		}
//	}
//}
//
//// BenchmarkDecodeGob benchmarks the decodeGob function
//func BenchmarkDecodeGob(b *testing.B) {
//	employee := &TestEmployee{
//		ID:     1,
//		Name:   "John Doe",
//		Salary: 50000,
//	}
//
//	buf, err := encodeGob(employee)
//	if err != nil {
//		b.Fatalf("encodeGob() error = %v", err)
//	}
//
//	b.ResetTimer() // resets the timer to ignore the time taken for setup
//
//	for i := 0; i < b.N; i++ {
//		buf := bytes.NewBuffer(buf.Bytes())
//		decodedEmployee := new(TestEmployee)
//		err = decodeGob(buf, decodedEmployee)
//		if err != nil {
//			b.Errorf("decodeGob() error = %v", err)
//		}
//	}
//}
//
//
//// BenchmarkEncodeJson benchmarks the encodeJson function
//func BenchmarkEncodeJson(b *testing.B) {
//	employee := &TestEmployee{
//		ID:     1,
//		Name:   "John Doe",
//		Salary: 50000,
//	}
//
//	b.ResetTimer() // resets the timer to ignore the time taken for setup
//
//	for i := 0; i < b.N; i++ {
//		_, err := encodeJson(employee)
//		if err != nil {
//			b.Errorf("encodeJson() error = %v", err)
//		}
//	}
//}
//
//// BenchmarkDecodeJson benchmarks the decodeJson function
//func BenchmarkDecodeJson(b *testing.B) {
//	employee := &TestEmployee{
//		ID:     1,
//		Name:   "John Doe",
//		Salary: 50000,
//	}
//
//	buf, err := encodeJson(employee)
//	if err != nil {
//		b.Fatalf("encodeJson() error = %v", err)
//	}
//
//	b.ResetTimer() // resets the timer to ignore the time taken for setup
//
//	for i := 0; i < b.N; i++ {
//		buf := bytes.NewBuffer(buf.Bytes())
//		decodedEmployee := new(TestEmployee)
//		err = decodeJson(buf, decodedEmployee)
//		if err != nil {
//			b.Errorf("decodeJson() error = %v", err)
//		}
//	}
//}



// Benchmarks
func BenchmarkAddToContext(b *testing.B) {
	r := httptest.NewRequest("GET", "http://example.com", nil)
	key := "test"
	v := struct {
		Name string
		Age  int
	}{
		Name: "Alice",
		Age:  25,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := AddToContext(r, key, v)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFromContext(b *testing.B) {
	r := httptest.NewRequest("GET", "http://example.com", nil)
	key := "test"
	v := struct {
		Name string
		Age  int
	}{
		Name: "Alice",
		Age:  25,
	}

	err := AddToContext(r, key, v)
	if err != nil {
		b.Fatal(err)
	}

	var v2 struct {
		Name string
		Age  int
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = FromContext(r, key, &v2)
		if err != nil {
			b.Fatal(err)
		}
	}
}


// Tests
func TestAddToContext(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com", nil)
	key := "test"
	v := struct {
		Name string
		Age  int
	}{
		Name: "Alice",
		Age:  25,
	}

	err := AddToContext(r, key, v)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check if the value is correctly added
	data, ok := r.Context().Value(key).([]byte)
	if !ok {
		t.Errorf("unexpected error: unable to cast to []byte for key %s", key)
	}

	var retrievedV struct {
		Name string
		Age  int
	}

	if err := json.Unmarshal(data, &retrievedV); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if retrievedV.Name != v.Name || retrievedV.Age != v.Age {
		t.Errorf("unexpected value: got %v want %v", retrievedV, v)
	}
}

func TestFromContext(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com", nil)
	key := "test"
	v := struct {
		Name string
		Age  int
	}{
		Name: "Alice",
		Age:  25,
	}

	err := AddToContext(r, key, v)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var retrievedV struct {
		Name string
		Age  int
	}

	err = FromContext(r, key, &retrievedV)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check if the value is correctly retrieved
	if retrievedV.Name != v.Name || retrievedV.Age != v.Age {
		t.Errorf("unexpected value: got %v want %v", retrievedV, v)
	}
}



// Tests
func TestAddToContextV2(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com", nil)
	key := "test"
	v := struct {
		Name string
		Age  int
	}{
		Name: "Alice",
		Age:  25,
	}

	err := AddToContextV2(r, key, v)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check if the value is correctly added
	data, ok := r.Context().Value(key).(*bytes.Buffer)
	if !ok {
		t.Errorf("unexpected error: unable to cast to *bytes.Buffer for key %s", key)
	}

	var retrievedV struct {
		Name string
		Age  int
	}

	if err := json.NewDecoder(data).Decode(&retrievedV); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if retrievedV.Name != v.Name || retrievedV.Age != v.Age {
		t.Errorf("unexpected value: got %v want %v", retrievedV, v)
	}
}

func TestFromContextV2(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com", nil)
	key := "test"
	v := struct {
		Name string
		Age  int
	}{
		Name: "Alice",
		Age:  25,
	}

	err := AddToContextV2(r, key, v)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var retrievedV struct {
		Name string
		Age  int
	}

	err = FromContextV2(r, key, &retrievedV)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check if the value is correctly retrieved
	if retrievedV.Name != v.Name || retrievedV.Age != v.Age {
		t.Errorf("unexpected value: got %v want %v", retrievedV, v)
	}
}



// Benchmarks
func BenchmarkAddToContextV2(b *testing.B) {

	key := "test"
	v := struct {
		Name string
		Age  int
	}{
		Name: "Alice",
		Age:  25,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest("GET", "http://example.com", nil)
		err := AddToContextV2(r, key, v)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFromContextV2(b *testing.B) {
	key := "test"
	v := struct {
		Name string
		Age  int
	}{
		Name: "Alice",
		Age:  25,
	}

	var v2 struct {
		Name string
		Age  int
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest("GET", "http://example.com", nil)
		err := AddToContextV2(r, key, v)
		if err != nil {
			b.Fatal(err)
		}
		err = FromContextV2(r, key, &v2)
		if err != nil {
			b.Fatal(err)
		}
	}
}