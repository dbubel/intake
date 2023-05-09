package intake

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// AddToContext adds v to the request context.
func AddToContext(r *http.Request, key string, v interface{}) error {
	buf := bufferPool.Get().(*bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(v); err != nil {
		return err
	}

	*r = *r.WithContext(context.WithValue(r.Context(), key, buf.Bytes()))

	buf.Reset()
	bufferPool.Put(buf)
	return nil
}

// FromContext adds v to the request context. v must be json decode-able
func FromContext(r *http.Request, key string, v interface{}) error {
	data, ok := r.Context().Value(key).([]byte)
	buf := bufferPool.Get().(*bytes.Buffer)

	if err := gob.NewDecoder(data).Decode(v); err != nil {
		return err
	}
	return nil
}
