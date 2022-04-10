package user_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/n3rdkube/user-manager/internal/messanger"

	"github.com/n3rdkube/user-manager/internal/api/user"
)

func Test_Delete_Succeed(t *testing.T) {
	consumer := &messanger.MockMessagingSystem{}
	producer := &messanger.MockMessagingSystem{}

	uh := user.NewHandler(consumer, producer)
	rr := &httptest.ResponseRecorder{}
	id := uuid.New().String()
	req := craftRequest(http.MethodDelete, "/users/"+id, nil)
	uh.ServeHTTP(rr, &req)

	require.Equal(t, http.StatusAccepted, rr.Code)
	require.Len(t, producer.Received, 1)
	require.Equal(t, string(producer.Received[0].Data), fmt.Sprintf("{\"id\":\"%s\"}", id))
	require.Equal(t, producer.Received[0].MessageType, messanger.DeleteUser)

}

func Test_Delete_Fails(t *testing.T) {

	t.Run("wrongID", func(t *testing.T) {
		consumer := &messanger.MockMessagingSystem{}
		producer := &messanger.MockMessagingSystem{}

		uh := user.NewHandler(consumer, producer)
		rr := &httptest.ResponseRecorder{}

		req := craftRequest(http.MethodDelete, "/users/wrongID", nil)
		uh.ServeHTTP(rr, &req)

		require.Equal(t, http.StatusNotFound, rr.Code)
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

		req := craftRequest(http.MethodDelete, "/users/"+uuid.New().String(), nil)
		uh.ServeHTTP(rr, &req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
