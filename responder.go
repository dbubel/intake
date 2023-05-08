package intake

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// RespondJSON used for responding with a JSON payload
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

// RespondJSONStream same as RespondJSON but uses a JSON streaming encoder
func RespondJSONStream(w http.ResponseWriter, r *http.Request, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

// UnmarshalJSON uses a sync.Pool to decode payloads to v.
func UnmarshalJSONSync(r io.Reader, v interface{}) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	if err := json.NewDecoder(buf).Decode(v); err != nil {
		return err
	}

	buf.Reset()
	bufferPool.Put(buf)
	return nil
}

func UnmarshalJSONStream(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func UnmarshalJSON(r io.Reader, v interface{}) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

// RespondJSONSyncPool takes advantage of a sync.Pool and may provide slightly lower memory
// consumption.
func RespondJSONSyncPool(w http.ResponseWriter, r *http.Request, code int, data interface{}) (int, error) {
	buf := bufferPool.Get().(*bytes.Buffer)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(&data)

	w.Header().Set("Content-Type", "application/json")
	respCode, err := Respond(w, r, code, buf.Bytes())

	// free up the sync pool mem
	buf.Reset()
	bufferPool.Put(buf)
	return respCode, err
}

func RespondXML(w http.ResponseWriter, r *http.Request, code int, data interface{}) (int, error) {
	jsonData, err := xml.Marshal(data)
	if err != nil {
		resp, _ := xml.Marshal(map[string]string{
			"error":       err.Error(),
			"description": "error marshalling response to XML",
		})
		return Respond(w, r, http.StatusInternalServerError, resp)
	}
	w.Header().Set("Content-Type", "application/xml")
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
