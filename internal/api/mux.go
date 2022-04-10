package api

import (
	"net/http"

	"github.com/n3rdkube/user-manager/internal/api/health"
	"github.com/n3rdkube/user-manager/internal/api/user"
	ms "github.com/n3rdkube/user-manager/internal/messanger"
	"github.com/sirupsen/logrus"
)

const (
	// just to avoid a redirect to the below one
	userBaseURLWithSlash = "/users/"
	userBaseURL          = "/users"
	healthBaseURL        = "/health"
)

// NewServerMux creates a new instance of serverMux. The entrypoint of every request
func NewServerMux(consumer ms.MessageQueueRead, producer ms.MessageQueueWrite, middleware func(http.Handler) http.Handler) *http.ServeMux {
	uH := user.NewHandler(consumer, producer)
	hH := health.NewHandler()

	// https://github.com/gorilla/mux would have been another possibility
	// just sticking with golang std library
	mux := http.NewServeMux()
	mux.Handle(userBaseURL, middleware(uH))
	mux.Handle(userBaseURLWithSlash, middleware(uH))
	mux.Handle(healthBaseURL, middleware(hH))

	return mux
}

// LogHandler logs incoming request and results. It is just an example of middlewares.
// User could also pass chain of middlewares, not implemented for simplicity
func LogHandler(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Request %q received with path %q", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
	handler := http.HandlerFunc(f)
	return handler
}
