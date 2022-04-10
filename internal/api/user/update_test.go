package user_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/n3rdkube/user-manager/internal/messanger"

	"github.com/n3rdkube/user-manager/internal/api/user"
)

func Test_Update_Succeed(t *testing.T) {
	consumer := &messanger.MockMessagingSystem{}
	producer := &messanger.MockMessagingSystem{}

	uh := user.NewHandler(consumer, producer)
	rr := &httptest.ResponseRecorder{}

	u := getUser()
	u.ID = uuid.New().String()
	b, err := json.Marshal(u)
	require.NoError(t, err)

	req := craftRequest(http.MethodPut, "/users/"+u.ID, b)
	uh.ServeHTTP(rr, &req)

	require.Equal(t, http.StatusAccepted, rr.Code)
	require.Len(t, producer.Received, 1)
	require.Equal(t, producer.Received[0].Data, b)
	require.Equal(t, producer.Received[0].MessageType, messanger.UpdateUser)

}

func Test_Update_Fails(t *testing.T) {

	t.Run("no_body", func(t *testing.T) {
		consumer := &messanger.MockMessagingSystem{}
		producer := &messanger.MockMessagingSystem{}

		uh := user.NewHandler(consumer, producer)
		rr := &httptest.ResponseRecorder{}

		req := craftRequest(http.MethodPut, "/users/"+uuid.New().String(), nil)
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

		req := craftRequest(http.MethodPut, "/users/"+uuid.New().String(), b)
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

		req := craftRequest(http.MethodPut, "/users/"+uuid.New().String(), b)
		uh.ServeHTTP(rr, &req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
