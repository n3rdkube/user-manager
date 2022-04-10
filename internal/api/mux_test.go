package api_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/n3rdkube/user-manager/internal/api"
	"github.com/n3rdkube/user-manager/internal/messanger"
	"github.com/stretchr/testify/require"
)

// Here we test merely the mux code doing a sanity check
func Test_Mux_Redirects_Fail_With(t *testing.T) {
	mq := &messanger.MockMessagingSystem{}
	mux := api.NewServerMux(mq, mq, api.LogHandler)

	cases := map[string]http.Request{
		"random_string": {
			URL: &url.URL{
				Path: "/test",
			},
		},
		"root": {
			URL: &url.URL{
				Path: "/",
			},
		},
		"singular": {
			URL: &url.URL{
				Path: "/user",
			},
		},
	}
	for testCaseName, req := range cases {
		req := req
		rr := &httptest.ResponseRecorder{}

		t.Run(testCaseName, func(t *testing.T) {
			mux.ServeHTTP(rr, &req)
			require.Equal(t, http.StatusNotFound, rr.Code)
		})
	}
}

// Here we test merely the mux code doing a sanity check
func Test_Mux_Redirects_Succeed(t *testing.T) {
	mq := &messanger.MockMessagingSystem{}
	mux := api.NewServerMux(mq, mq, dumpRouteHandler)

	cases := map[string]http.Request{
		"users": {
			URL: &url.URL{
				Path: "/users",
			},
		},
		"users_with_slash": {
			URL: &url.URL{
				Path: "/users/",
			},
		},
		"health": {
			URL: &url.URL{
				Scheme: "http",
				Path:   "/health",
			},
		},
		"users_with_slash_and_path": {
			URL: &url.URL{
				Scheme: "http",
				Path:   "/users/whatever",
			},
		},
	}

	for testCaseName, req := range cases {
		req := req
		rr := &httptest.ResponseRecorder{}

		t.Run(testCaseName, func(t *testing.T) {
			mux.ServeHTTP(rr, &req)
			require.Equal(t, http.StatusOK, rr.Code)
		})
	}
}

// it does nothing, it simply answers with ok avoiding to call following handlers
func dumpRouteHandler(_ http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	handler := http.HandlerFunc(f)
	return handler
}
