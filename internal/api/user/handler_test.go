package user_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/n3rdkube/user-manager/internal/api/user"
	"github.com/stretchr/testify/require"
)

func Test_UserHandler_fail(t *testing.T) {

	hh := user.NewHandler(nil, nil)

	cases := map[string]http.Request{
		"wrong_method_connect": {
			Method: http.MethodConnect,
			URL: &url.URL{
				Path: "/users",
			},
		},
		"wrong_method_on_path": {
			Method: http.MethodPut,
			URL: &url.URL{
				Path: "/users",
			},
		},
		"wrong_url": {
			Method: http.MethodGet,
			URL: &url.URL{
				Path: "/users/asd",
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
