package middlware

import (
	"fmt"
	"github.com/dbubel/intake"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Recover needs to be the first middleware in chain
func Recover(next intake.Handler) intake.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				intake.RespondError(w, r, fmt.Errorf("error panic"), http.StatusInternalServerError, "server recovered from a panic")
			}
		}()
		next(w, r, params)
	}
}
