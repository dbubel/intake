// Package intake provides HTTP response utilities for common content types.
package intake

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// RespondJSON writes a JSON response with the specified HTTP status code.
// It automatically sets the Content-Type header to "application/json" and
// marshals the provided data into JSON format. If marshaling fails, the error
// is returned to the caller.
func RespondJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

// RespondXML writes an XML response with the specified HTTP status code.
// It automatically sets the Content-Type header to "application/xml" and
// marshals the provided data into XML format. If marshaling fails, the error
// is returned to the caller.
func RespondXML(w http.ResponseWriter, r *http.Request, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	return xml.NewEncoder(w).Encode(data)
}

// Respond writes raw bytes as an HTTP response with the specified status code.
// If no Content-Type header is set, it attempts to detect the content type from
// the data. Returns the number of bytes written and any error that occurred.
func Respond(w http.ResponseWriter, r *http.Request, code int, data []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(data))
	}

	w.WriteHeader(code)
	return w.Write(data)
}
