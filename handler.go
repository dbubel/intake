package intake

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handler a compatible type for httprouter
type Handler func(w http.ResponseWriter, r *http.Request, params httprouter.Params)
