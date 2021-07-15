package intake

import (
	"fmt"
	"github.com/go-playground/validator"
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
