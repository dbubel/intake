package middlware

import (
	"context"
	"github.com/dbubel/intake"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

// Timeout if added will created a context that will cancel after t. This cancel
// will affect all downstream uses of the context
func Timeout(t time.Duration) func(handler intake.Handler) intake.Handler {
	return func(next intake.Handler) intake.Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			// Create a context that is both manually cancellable and will signal
			// cancel at the specified duration.
			ctx, cancel := context.WithTimeout(r.Context(), t)
			defer cancel()
			next(w, r.WithContext(ctx), params)
		}
	}
}
