package intake

import (
	"bytes"
	"context"
	"encoding/gob"
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
	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	if err := e.Encode(v); err != nil {
		return err
	}
	*r = *r.WithContext(context.WithValue(r.Context(), key, b))
	return nil
}

func FromContext(r *http.Request, key string, v interface{}) error {
	data, ok := r.Context().Value(key).(bytes.Buffer)
	if !ok {
		return fmt.Errorf("error casting from context for (%s)", key)
	}
	d := gob.NewDecoder(&data)
	if err := d.Decode(v); err != nil {
		return err
	}
	return nil
}
