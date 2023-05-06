package intake

import (
	"encoding/json"
	"io"
)

type StreamingEndpoint struct {
	Endpoint
	Encoder *io.Writer
}

func NewStreamingJsonEndpoint(method, path string, encoder io.Writer, endpointHandler Handler, mid ...MiddleWare) StreamingEndpoint {

	return StreamingEndpoint{
		Endpoint: NewEndpoint(method,path,endpointHandler,mid...),
		Encoder:  json.NewEncoder(encoder),
	}
}
