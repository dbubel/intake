package intake

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestContextHelpers(t *testing.T) {
	t.Run("Test context with struct", func(t *testing.T) {
		type fake struct {
			UserName string
			Address  string
		}
		testEp := fmt.Sprintf("/%s", RandStringRunes(20))
		var app = NewDefault()
		testHandler := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
			var f fake
			FromContext(r, "key", &f)
			RespondJSON(w, r, http.StatusOK, f)
		}

		app.AddEndpoint(http.MethodGet, testEp, testHandler, wrap("key", fake{
			UserName: "tom",
			Address:  "addr",
		}))

		r := httptest.NewRequest(http.MethodGet, testEp, nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		resp, _ := ioutil.ReadAll(w.Body)
		assert.JSONEq(t, `{"Address":"addr", "UserName":"tom"}`, string(resp))
	})

	t.Run("Test context with string", func(t *testing.T) {
		testEp := fmt.Sprintf("/%s", RandStringRunes(20))
		var app = NewDefault()
		testHandler := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
			var s string
			FromContext(r, "key", &s)
			Respond(w, r, http.StatusOK, []byte(s))
		}

		app.AddEndpoint(http.MethodGet, testEp, testHandler, wrap("key", "hello"))

		r := httptest.NewRequest(http.MethodGet, testEp, nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		resp, _ := ioutil.ReadAll(w.Body)
		assert.Equal(t, "hello", string(resp))
	})
}

func wrap(key string, val interface{}) func(next Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			AddToContext(r, key, val)
			next(w, r, params)
		}
	}
}
