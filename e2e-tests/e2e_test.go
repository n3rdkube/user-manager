//go:build e2e

package e2e_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/dchest/uniuri"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/n3rdkube/user-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testURL = "http://localhost:35307"

const (
	lastName = "testName"
	email    = "testEmail"
	country  = "testCountry"
	nickname = "testNickname"
)

func Test_Server_manages_wrong_urls(t *testing.T) {
	client := resty.New().SetBaseURL(testURL)
	waitServerReady(t, client)

	resp, err := client.R().Get("/notExisting")

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode())
}

func Test_Server(t *testing.T) {
	client := resty.New().SetBaseURL(testURL)
	waitServerReady(t, client)

	u := getUser()
	u.ID = uuid.New().String()
	rand := uniuri.New()
	u.FirstName = rand

	t.Run("create_user", func(t *testing.T) {
		b, err := json.Marshal(u)
		require.NoError(t, err)

		resp, err := client.R().SetBody(b).Post("/users")
		require.NoError(t, err)
		require.Equal(t, http.StatusAccepted, resp.StatusCode())
		require.Equal(t, []byte{}, resp.Body())

		time.Sleep(time.Second)
		t.Run("list_with_options", func(t *testing.T) {
			resp, err = client.R().SetQueryParam("include.firstname", rand).Get("/users")
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode())

			var list models.ListUsers
			err = json.Unmarshal(resp.Body(), &list)
			require.NoError(t, err)
			require.Equal(t, models.ListUsers{u}, list)
		})
	})

	u.NickName = "newNickname"

	t.Run("update_user", func(t *testing.T) {
		b, err := json.Marshal(u)
		require.NoError(t, err)

		resp, err := client.R().SetBody(b).Put(path.Join("/users/", u.ID))
		require.NoError(t, err)
		require.Equal(t, http.StatusAccepted, resp.StatusCode())
		require.Equal(t, []byte{}, resp.Body())

		time.Sleep(time.Second)
		t.Run("list_with_options", func(t *testing.T) {
			resp, err = client.R().SetQueryParam("include.firstname", rand).Get("/users")
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode())

			var list models.ListUsers
			err = json.Unmarshal(resp.Body(), &list)
			require.NoError(t, err)
			require.Equal(t, models.ListUsers{u}, list)
		})
	})

	t.Run("delete_user", func(t *testing.T) {
		resp, err := client.R().Delete(path.Join("/users/", u.ID))
		require.NoError(t, err)
		require.Equal(t, http.StatusAccepted, resp.StatusCode())

		t.Run("list_with_options", func(t *testing.T) {
			resp, err = client.R().SetQueryParam("include.firstname", rand).Get("/users")
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode())

			var list models.ListUsers
			err = json.Unmarshal(resp.Body(), &list)
			require.NoError(t, err)
			require.Len(t, list, 0)
		})
	})
}

func Test_With_Parallel_Requests(t *testing.T) {
	client := resty.New().SetBaseURL(testURL)
	waitServerReady(t, client)

	country := uniuri.New()

	for i := 0; i < 10; i++ {
		go t.Run("create_user_succeed", func(t *testing.T) {
			u := getUser()
			u.Country = country
			u.ID = uuid.New().String()
			u.FirstName = uniuri.New()

			b, err := json.Marshal(u)
			require.NoError(t, err)

			resp, err := client.R().SetBody(b).Post("/users")
			require.NoError(t, err)
			require.Equal(t, http.StatusAccepted, resp.StatusCode())
			require.Equal(t, []byte{}, resp.Body())

		})

		t.Run("list_succeed", func(t *testing.T) {
			resp, err := client.R().Get("/users")
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode())

			var list models.ListUsers
			err = json.Unmarshal(resp.Body(), &list)
			require.NoError(t, err)
		})
	}

	time.Sleep(2 * time.Second)

	t.Run("all_user_were_created", func(t *testing.T) {
		resp, err := client.R().SetQueryParam("include.country", country).Get("/users")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())

		var list models.ListUsers
		err = json.Unmarshal(resp.Body(), &list)
		require.NoError(t, err)
		require.Len(t, list, 10)
	})

}

func Test_Pagination(t *testing.T) {
	client := resty.New().SetBaseURL(testURL)
	waitServerReady(t, client)

	country := uniuri.New()

	for i := 0; i < 10; i++ {
		go t.Run("create_user_succeed", func(t *testing.T) {
			u := getUser()
			u.Country = country
			u.ID = uuid.New().String()
			u.FirstName = uniuri.New()

			b, err := json.Marshal(u)
			require.NoError(t, err)

			resp, err := client.R().SetBody(b).Post("/users")
			require.NoError(t, err)
			require.Equal(t, http.StatusAccepted, resp.StatusCode())
			require.Equal(t, []byte{}, resp.Body())
		})
	}

	time.Sleep(3 * time.Second)

	cases := map[string]struct {
		numberResults int
		page          string
		rows          string
	}{
		"allData": {
			numberResults: 10,
			page:          "1",
			rows:          "100",
		},
		"noData": {
			numberResults: 0,
			page:          "2",
			rows:          "100",
		},
		"fullFirstPage": {
			numberResults: 3,
			page:          "1",
			rows:          "3",
		},
		"fullSecondPage": {
			numberResults: 3,
			page:          "2",
			rows:          "3",
		},
		"fullThirdPage": {
			numberResults: 3,
			page:          "3",
			rows:          "3",
		},
		"PartialFourthPage": {
			numberResults: 1,
			page:          "4",
			rows:          "3",
		},
	}

	for testCaseName, testData := range cases {
		testData := testData
		t.Run(testCaseName, func(t *testing.T) {
			resp, err := client.R().
				SetQueryParam("include.country", country).
				SetQueryParam("page.number", testData.page).
				SetQueryParam("page.rows", testData.rows).
				Get("/users")
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode())

			var list models.ListUsers
			err = json.Unmarshal(resp.Body(), &list)
			assert.NoError(t, err)
			assert.Len(t, list, testData.numberResults)
		})
	}
}

func waitServerReady(t *testing.T, client *resty.Client) {
	t.Run("is_ready", func(t *testing.T) {

		err := retry.Do(
			func() error {
				resp, err := client.R().Get("/health")
				if err != nil {
					return err
				}

				if resp.StatusCode() != http.StatusOK {
					return fmt.Errorf("waiting for manager to be ready")
				}

				return nil
			},
			retry.Delay(time.Second*2),
		)
		assert.NoError(t, err)
	})
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
