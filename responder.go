package intake

import (
	"encoding/json"
	"encoding/xml"
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
	return Respond(w, r, code, jsonData)
}

func RespondXML(w http.ResponseWriter, r *http.Request, code int, data interface{}) (int, error) {
	jsonData, err := xml.Marshal(data)
	if err != nil {
		resp, _ := xml.Marshal(map[string]string{
			"error":       err.Error(),
			"description": "error marshalling response to JSON",
		})
		return Respond(w, r, http.StatusInternalServerError, resp)

	}
	return Respond(w, r, code, jsonData)
}

func Respond(w http.ResponseWriter, r *http.Request, code int, data []byte) (int, error) {
	err := AddToContext(r, "response-code", code)
	if err != nil {
		return 0, err
	}

	err = AddToContext(r, "response-length", len(data))
	if err != nil {
		return 0, err
	}
	contentType := http.DetectContentType(data)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	return w.Write(data)
}
