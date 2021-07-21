package intake

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"
)

func UnmarshalJSON(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}

	if fve := validate.Struct(v); fve != nil {
		for _, fe := range fve.(validator.ValidationErrors) {
			return Invalid{
				Fld:  fe.Field(),
				Err:  fe.Tag(),
				Kind: fe.Kind().String(),
			}
		}
	}
	return nil
}

func AddToContext(r *http.Request, key string, v interface{}) error {
	encoded, err := json.Marshal(v)
	if err != nil {
		return err
	}
	*r = *r.WithContext(context.WithValue(r.Context(), key, encoded))
	return nil
}

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
