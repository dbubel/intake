package intake

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AddToContext adds v to the request context. v must be json encode-able
func AddToContext(r *http.Request, key string, v interface{}) error {
	encoded, err := json.Marshal(v)
	if err != nil {
		return err
	}
	*r = *r.WithContext(context.WithValue(r.Context(), key, encoded))
	return nil
}

// FromContext adds v to the request context. v must be json decode-able
func FromContext(r *http.Request, key string, v interface{}) error {
	data, ok := r.Context().Value(key).([]byte)
	if !ok {
		return fmt.Errorf("error casting to []byte for key %s", key)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}
