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

	"github.com/n3rdkube/user-manager/internal/models"

	"github.com/stretchr/testify/require"

	"github.com/n3rdkube/user-manager/internal/messanger"

	"github.com/n3rdkube/user-manager/internal/api/user"
)

func Test_List_Succeed(t *testing.T) {
	u := getUser()
	b, err := json.Marshal(models.ListUsers{u})
	require.NoError(t, err)

	t.Run("with_no_list_options", func(t *testing.T) {
		consumer := &messanger.MockMessagingSystem{
			MessageToReturn: []messanger.Message{
				{
					Data:        b,
					MessageType: messanger.ListUserAnswer,
				},
			},
		}
		producer := &messanger.MockMessagingSystem{}

		uh := user.NewHandler(consumer, producer)
		body := bytes.Buffer{}
		rr := &httptest.ResponseRecorder{Body: &body}

		req := craftRequest(http.MethodGet, "/users", nil)
		uh.ServeHTTP(rr, &req)
		require.Equal(t, http.StatusOK, rr.Code)

		// testing that body is what we expect
		bodyBytes, err := io.ReadAll(&body)
		require.NoError(t, err)
		var list models.ListUsers
		err = json.Unmarshal(bodyBytes, &list)
		require.NoError(t, err)
		require.Equal(t, list[0], u)
	})
	t.Run("and_list_options_are_correctly_generated", func(t *testing.T) {
		consumer := &messanger.MockMessagingSystem{
			MessageToReturn: []messanger.Message{
				{
					Data:        b,
					MessageType: messanger.ListUserAnswer,
				},
			},
		}
		producer := &messanger.MockMessagingSystem{}

		uh := user.NewHandler(consumer, producer)
		body := bytes.Buffer{}
		rr := &httptest.ResponseRecorder{Body: &body}

		req := craftRequest(http.MethodGet, "/users", nil)
		req.URL, err = url.Parse("http://localhost/users?include.firstname=paolo&page.number=1&page.rows=3")
		require.NoError(t, err)

		//testing that the list options request is what we expect
		uh.ServeHTTP(rr, &req)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Len(t, producer.Received, 1)
		require.Equal(t, string(producer.Received[0].Data), "{\"include\":{\"first_name\":\"paolo\"},\"page_number\":1,\"rows_per_page\":3}")
	})

}

func Test_List_Fails(t *testing.T) {

	t.Run("error_from_rabbit_mq", func(t *testing.T) {
		consumer := &messanger.MockMessagingSystem{}
		producer := &messanger.MockMessagingSystem{
			ErrorToReturn: []error{
				errors.New("generic Error"),
			},
		}

		uh := user.NewHandler(consumer, producer)
		rr := &httptest.ResponseRecorder{}

		req := craftRequest(http.MethodGet, "/users", nil)
		uh.ServeHTTP(rr, &req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("error_from_rabbit_mq_callback", func(t *testing.T) {
		consumer := &messanger.MockMessagingSystem{
			ErrorToReturn: []error{
				errors.New("generic Error"),
			},
		}
		producer := &messanger.MockMessagingSystem{}

		uh := user.NewHandler(consumer, producer)
		rr := &httptest.ResponseRecorder{}

		req := craftRequest(http.MethodGet, "/users", nil)
		uh.ServeHTTP(rr, &req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
