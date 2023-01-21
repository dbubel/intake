# Intake
API framework written in Go. Intake uses `httprouter` (https://github.com/julienschmidt/httprouter) under the hood. 

Intake was written to be a simplistic framework for writing API servers. It was designed to not hide whats going on 
during the lifetime of a request. In the spirit of Go, verbosity was chosen as to make the code more readable. I 
believe that with this approach APIs built using `Intake` will be more maintable and easily modified for the life of 
the application.

Sample server
```go
func simple(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	intake.Respond(w, r, http.StatusOK, []byte("testing get"))
	return
}

func main() {
	app := intake.NewDefault()
	app.AddEndpoint(http.MethodGet,"/test-get", simple)
	app.Run(&http.Server{
		Addr:           ":8000",
		Handler:        app.Router,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	})
}
```

### Endpoint groups
```go
app := intake.NewDefault()
endpoints := intake.Endpoints{
    intake.NewEndpoint(http.MethodGet,"/test-ep-one", testEndpointOne),
    intake.NewEndpoint(http.MethodGet,"/test-ep-two", testEndpointTwo),
}
app.AddEndpoints(endpoints)
```

### Middleware
Middleware is executed before the final endpoint handler in the order they are added.
```go
func someMiddleware(next intake.Handler) intake.Handler {
    return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
        // do some block of logic needed for downstream handlers
        ctx := context.WithValue(r.Context(), "key", "value")
        next(w, r.WithContext(ctx), params)
    }
}
app.AddEndpoint(http.MethodGet,"/test-in-the-middle",finalHandler,someMiddleware)
```

### Prepend and Append middleware
If you have a group of middlewares and you want to add another at the start or end of the chain, use the functions below.
This will not affect global middleware always being first in the function calls.
```go
app := intake.NewDefault()
endpoints := intake.Endpoints{
    intake.NewEndpoint(http.MethodGet,"/test-ep-one", testEndpointOne),
    intake.NewEndpoint(http.MethodGet,"/test-ep-two", testEndpointTwo),
}
endpoints.Prepend(someMiddleware)
endpoints.Append(someOtherMiddleware)
```

#### Middleware groups
Middleware groups are groups of endpoints that a middleware handler is applied to. 
The middleware is applied to ALL endpoints in the group. 
```go
app := intake.NewDefault()
endpoints := intake.Endpoints{
    intake.NewEndpoint(http.MethodGet,"/test-middleware", finalHandler),
}

endpoints.Use(someMiddleware)
```

#### Global middleware
Global middleware is applied to ALL endpoints associated with the intake app. 
```go
app := intake.NewDefault()
app.AddGlobalMiddleware(someMiddleware)
```

### Custom logger
The intake app can accept a `logrus` custom logger.
```go 
l := logrus.New()
l.SetLevel(logrus.DebugLevel)
l.SetFormatter(&logrus.JSONFormatter{})
app := intake.New(l)
```

### Request validation
`Intake` provides a basic `UnmarshalJSON` function for validating requests (github.com/go-playground/validator)
```go
type sample struct {
	Username string
	Address string
}

func SimpleHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var s sample
	if err := intake.UnmarshalJSON(r.Body, &s);err != nil {
		intake.RespondError(w,r,err, http.StatusBadRequest, "description of error here")
		return
	}
	// do other things if valid 
	intake.Respond(w, r, http.StatusOK, []byte("success"))
	return
}
```

### Request context
Usually used in middleware when doing distinct things on the request

#### Adding objects to the request
Adding the struct makes it available to downstream middleware
```go
func someMiddleware(next intake.Handler) intake.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		a := SomeStruct{
			Name:    "Tom",
			Address: "1234 Drive",
		}
		intake.AddToContext(r, "userData", a)
		next(w, r, params)
	}
}
```
#### Getting objects from the context
```go
func finalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userData SomeStruct
	intake.FromContext(r,"userData",&userData)
    //...
}
```

