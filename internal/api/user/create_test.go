package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"

	"github.com/n3rdkube/user-manager/internal/models"

	"github.com/n3rdkube/user-manager/internal/messanger"

	"github.com/n3rdkube/user-manager/internal/api/user"
	"github.com/stretchr/testify/require"
)

const (
	lastName = "testName"
	email    = "testEmail"
	country  = "testCountry"
	nickname = "testNickname"
)

func Test_Create_Succeed(t *testing.T) {
	consumer := &messanger.MockMessagingSystem{}
	producer := &messanger.MockMessagingSystem{}

	uh := user.NewHandler(consumer, producer)
	rr := &httptest.ResponseRecorder{}

	u := getUser()
	u.ID = uuid.New().String()
	b, err := json.Marshal(u)
	require.NoError(t, err)

	req := craftRequest(http.MethodPost, "/users", b)
	uh.ServeHTTP(rr, &req)

	require.Equal(t, http.StatusAccepted, rr.Code)
	require.Len(t, producer.Received, 1)
	require.Equal(t, producer.Received[0].Data, b)
	require.Equal(t, producer.Received[0].MessageType, messanger.CreateUser)

}

func Test_Create_Fails(t *testing.T) {

	t.Run("no_body", func(t *testing.T) {
		consumer := &messanger.MockMessagingSystem{}
		producer := &messanger.MockMessagingSystem{}

		uh := user.NewHandler(consumer, producer)
		rr := &httptest.ResponseRecorder{}

		req := craftRequest(http.MethodPost, "/users", nil)
		uh.ServeHTTP(rr, &req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("no_userID", func(t *testing.T) {
		consumer := &messanger.MockMessagingSystem{}
		producer := &messanger.MockMessagingSystem{}

		uh := user.NewHandler(consumer, producer)
		rr := &httptest.ResponseRecorder{}

		u := getUser()
		b, err := json.Marshal(u)
		require.NoError(t, err)

		req := craftRequest(http.MethodPost, "/users", b)
		uh.ServeHTTP(rr, &req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("error_from_rabbit_mq", func(t *testing.T) {
		consumer := &messanger.MockMessagingSystem{}
		producer := &messanger.MockMessagingSystem{
			ErrorToReturn: []error{
				errors.New("generic Error"),
			},
		}

		uh := user.NewHandler(consumer, producer)
		rr := &httptest.ResponseRecorder{}

		u := getUser()
		u.ID = uuid.New().String()
		b, err := json.Marshal(u)
		require.NoError(t, err)

		req := craftRequest(http.MethodPost, "/users", b)
		uh.ServeHTTP(rr, &req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func craftRequest(method string, path string, data []byte) http.Request {
	return http.Request{
		Method: method,
		URL: &url.URL{
			Path: path,
		},
		Body: io.NopCloser(bytes.NewReader(data)),
	}
}

func getUser() models.User {
	return models.User{
		// FirstName: "", This is generated
		Country:  country,
		Email:    email,
		NickName: nickname,
		LastName: lastName,
	}
}
