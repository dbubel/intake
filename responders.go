package intake

import (
	"encoding/gob"
	"encoding/json"
	"net/http"
)

// RespondJSON same as RespondJSON but uses a JSON streaming encoder
func RespondJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) error {
	if err := AddToContext(r, "response-code", code); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func RespondGOB(w http.ResponseWriter, r *http.Request, code int, data interface{}) error {
	if err := AddToContext(r, "response-code", code); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(code)
	return gob.NewEncoder(w).Encode(data)
}

func Respond(w http.ResponseWriter, r *http.Request, code int, data []byte) (int, error) {
	err := AddToContext(r, "response-code", code)
	if err != nil {
		return -1, err
	}

	err = AddToContext(r, "response-length", len(data))
	if err != nil {
		return -1, err
	}

	// content type is not set so attempt to set it
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(data))
	}

	w.WriteHeader(code)
	return w.Write(data)
}

type jsonStreamer struct {
	encoder *json.Encoder
}

func NewStreamingJSONWriter(buf *http.ResponseWriter) *jsonStreamer {
	return &jsonStreamer{
		encoder: json.NewEncoder(*buf),
	}
}

func (j jsonStreamer) Write(data interface{}) error {
	return j.encoder.Encode(data)
}
