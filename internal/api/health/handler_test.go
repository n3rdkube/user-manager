package health_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/n3rdkube/user-manager/internal/api/health"
	"github.com/stretchr/testify/require"
)

func Test_HealthHandler_Success(t *testing.T) {
	hh := health.NewHandler()

	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/health",
		},
	}

	rr := &httptest.ResponseRecorder{}
	hh.ServeHTTP(rr, &req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func Test_HealthHandler_fail(t *testing.T) {

	hh := health.NewHandler()

	cases := map[string]http.Request{
		"wrong_method": {
			Method: http.MethodPost,
			URL: &url.URL{
				Path: "/health",
			},
		},
		"wrong_url": {
			Method: http.MethodGet,
			URL: &url.URL{
				Path: "/health/asd",
			},
		},
	}

	for testCaseName, req := range cases {
		req := req
		rr := &httptest.ResponseRecorder{}

		t.Run(testCaseName, func(t *testing.T) {
			hh.ServeHTTP(rr, &req)
			require.Equal(t, http.StatusNotFound, rr.Code)
		})
	}
}
