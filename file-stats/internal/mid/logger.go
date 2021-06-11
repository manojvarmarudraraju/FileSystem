package mid

import (
	"context"
	"errors"
	"file-stats/internal/web"
	"log"
	"net/http"
	"time"
)

// Logger handles requests before and after serving the request handle (logging for every request)
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			//Code to be executed before the handler method
			//st:=time.Now()
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok{
				return errors.New("web values missing from context")
			}
			v.Start = time.Now()

			// Run the handler chain
			err := before(ctx, w, r)

			// Code to be executed after the handler chain
			log.Printf("(%v) %s %s (%v)", v.StatusCode, r.Method, r.URL.Path, time.Since(v.Start))

			return err
		}

		return h
	}

	return f
}