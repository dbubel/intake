package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"io"
)

var validate = validator.New()

type Invalid struct {
	Fld  string `json:"field_name"`
	Err  string `json:"error"`
	Kind string `json:"kind"`
}

func (err Invalid) Error() string {
	return fmt.Sprintf("field %s %s type: %s", err.Fld, err.Err, err.Kind)
}

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
