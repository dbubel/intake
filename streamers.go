package intake

import (
	"encoding/json"
	"net/http"
)

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
