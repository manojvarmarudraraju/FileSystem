package web

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"time"
)
// App is the entry point for all webapplicaitons.
type App struct {
	mux *chi.Mux
	log *log.Logger
	mw []Middleware
}

//ctxKey represents the type of value for the context key
type ctxKey int
//keyValues is how request values are store/retrieved.
const KeyValues ctxKey = 1

//type to carry information about the request.
type Values struct {
	StatusCode int
	Start time.Time
}
//handler type is the signature that all handlers will apply.
type Handler func(context.Context, http.ResponseWriter, *http.Request)error

// NewApp knows how to construct internal state for an App.
func NewApp(logger *log.Logger, mw ...Middleware) *App  {
	return &App{
		mux: GetMuxWithoutCors(),
		log: logger,
		mw: mw,
	}
}

//configure cors.

func GetMuxWithoutCors()*chi.Mux  {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	}))
	return mux
}

//handle connects a pattern to the handler function.
func (a *App) Handle(method, pattern string, h Handler, m ...Middleware )  { //fn http.HandlerFunc
	// First wrap handler specific middleware around this handler.
	h = wrapMiddleware(m, h)

	// Add the application's general middleware to the handler chain.
	//Once we get the method, pattern and handler we should wrap the handler with the middleware functions.
	h = wrapMiddleware(a.mw,h)
	fn:= func(w http.ResponseWriter, r * http.Request) {

		//place to populate the required into request.
		v:=Values{
			Start:      time.Now(),
		}

		//we should be not basic type are key in context to put values, so build your type
		ctx := context.WithValue(r.Context(), KeyValues,&v)

		if err:= h(ctx,w,r); err!=nil{
			// replace with reusable code or remove the business logic from handler
			/*resp:=ErrorResponse{Error:err.Error()}
		if err := Respond(w,resp,http.StatusInternalServerError); err!=nil{
			a.log.Println(err)
		}*/
			a.log.Printf("Error: Unhandled error %v",err)
			//built a new middle ware to handle application errors.
			/*if err:=RespondError(w,err); err!=nil{
				a.log.Printf("Error: %v",err)
			}*/
		}

	}

	a.mux.MethodFunc(method, pattern, fn)
}

func (a * App) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	a.mux.ServeHTTP(w,r)
}
