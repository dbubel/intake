// Package intake provides HTTP response utilities for common content types.
// This file contains helper functions for responding to HTTP requests with
// different content types such as JSON, XML, and raw bytes.
package intake

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// RespondJSON writes a JSON response with the specified HTTP status code.
// It automatically sets the Content-Type header to "application/json" and
// marshals the provided data into JSON format.
//
// Parameters:
//   - w: The HTTP response writer to write the response to
//   - r: The HTTP request that triggered this response
//   - code: The HTTP status code to send
//   - data: The data to marshal as JSON
//
// Returns:
//   - An error if JSON marshaling fails, nil otherwise
func RespondJSON(w http.ResponseWriter, r *http.Request, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

// RespondXML writes an XML response with the specified HTTP status code.
// It automatically sets the Content-Type header to "application/xml" and
// marshals the provided data into XML format.
//
// Parameters:
//   - w: The HTTP response writer to write the response to
//   - r: The HTTP request that triggered this response
//   - code: The HTTP status code to send
//   - data: The data to marshal as XML
//
// Returns:
//   - An error if XML marshaling fails, nil otherwise
func RespondXML(w http.ResponseWriter, r *http.Request, code int, data any) error {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	return xml.NewEncoder(w).Encode(data)
}

// Respond writes raw bytes as an HTTP response with the specified status code.
// If no Content-Type header is set, it attempts to detect the content type from
// the data using http.DetectContentType.
//
// Parameters:
//   - w: The HTTP response writer to write the response to
//   - r: The HTTP request that triggered this response
//   - code: The HTTP status code to send
//   - data: The raw bytes to write as the response body
//
// Returns:
//   - The number of bytes written to the response
//   - An error if writing to the response fails, nil otherwise
func Respond(w http.ResponseWriter, r *http.Request, code int, data []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(data))
	}

	w.WriteHeader(code)
	return w.Write(data)
}
