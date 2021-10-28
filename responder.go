package intake

import (
	"encoding/json"
	"net/http"
)

func RespondError(w http.ResponseWriter, r *http.Request, err error, code int, description ...string) (int, error) {
	return RespondJSON(w, r, code, map[string]interface{}{"error": err.Error(), "description": description})
}

func RespondJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) (int, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		resp, _ := json.Marshal(map[string]string{
			"error":       err.Error(),
			"description": "error marshalling response to JSON",
		})
		return Respond(w, r, http.StatusInternalServerError, resp)
	}
	w.Header().Set("Content-Type", "application/json")
	return Respond(w, r, code, jsonData)
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
