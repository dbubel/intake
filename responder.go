package intake

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func RespondXML(w http.ResponseWriter, r *http.Request, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return xml.NewEncoder(w).Encode(data)
}

func Respond(w http.ResponseWriter, r *http.Request, code int, data []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(data))
	}

	w.WriteHeader(code)
	return w.Write(data)
}
