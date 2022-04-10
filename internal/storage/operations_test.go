package storage_test

import (
	"testing"

	"github.com/n3rdkube/user-manager/internal/storage"

	"github.com/google/uuid"
	"github.com/n3rdkube/user-manager/internal/models"
	"github.com/stretchr/testify/require"
)

const (
	firstName = "testName"
	lastName  = "lastName"
	email     = "testEmail"
	country   = "testCountry"
	nickname  = "testNickname"
)

func TestCompleteFlow(t *testing.T) {

	db, err := storage.NewStorageDBInMemory("TestCompleteFlow")
	require.NoError(t, err)
	defer db.Close()

	//cannot update, no user present
	user := getUserTest()
	err = db.UpdateUser(user)
	require.Error(t, err)

	//create a user
	err = db.AddUser(user)
	require.NoError(t, err)

	//updating the user
	user.Country = "updated"
	err = db.UpdateUser(user)
	require.NoError(t, err)

	//list the user
	list, err := db.ListUser(models.ListOptions{})
	require.NoError(t, err)

	require.Len(t, list, 1)
	require.Equal(t, firstName, list[0].FirstName)

	// delete do not fail
	err = db.DeleteUser(user.ID)
	require.NoError(t, err)
}

func TestList(t *testing.T) {
	db, err := storage.NewStorageDBInMemory("TestList")
	require.NoError(t, err)
	user := getUserTest()

	err = db.AddUser(user)
	require.NoError(t, err)

	user.NickName = "differentNickname"
	err = db.AddUser(user)
	require.NoError(t, err)

	list, _ := db.ListUser(models.ListOptions{})
	require.Len(t, list, 2)

	t.Run("pagination", func(t *testing.T) {
		list, _ = db.ListUser(models.ListOptions{
			PageNumber:  1,
			RowsPerPage: 2,
		})
		require.Len(t, list, 2)

		list, _ = db.ListUser(models.ListOptions{
			PageNumber:  1,
			RowsPerPage: 1,
		})
		require.Len(t, list, 1)
		require.Equal(t, nickname, list[0].NickName)

		list, _ = db.ListUser(models.ListOptions{
			PageNumber:  2,
			RowsPerPage: 1,
		})
		require.Len(t, list, 1)
		require.Equal(t, "differentNickname", list[0].NickName)
	})
}

func TestListOptions(t *testing.T) {
	db, err := storage.NewStorageDBInMemory("TestListOptions")
	require.NoError(t, err)
	user := getUserTest()

	err = db.AddUser(user)
	require.NoError(t, err)

	list, _ := db.ListUser(models.ListOptions{})
	require.Len(t, list, 1)

	list, _ = db.ListUser(models.ListOptions{
		Include: models.User{
			FirstName: firstName,
			Country:   country,
			LastName:  lastName,
		},
	})
	require.Len(t, list, 0)

	list, _ = db.ListUser(models.ListOptions{
		Include: models.User{
			FirstName: "differentFirstName",
		},
	})
	require.Len(t, list, 0)

	list, _ = db.ListUser(models.ListOptions{
		Include: models.User{
			Country: "differentCountry",
		},
	})
	require.Len(t, list, 0)

	list, _ = db.ListUser(models.ListOptions{
		Include: models.User{
			LastName: "differentLastName",
		},
	})
	require.Len(t, list, 0)
}

func TestApisDelete(t *testing.T) {
	db, err := storage.NewStorageDBInMemory("TestApisDelete")
	require.NoError(t, err)

	user := getUserTest()

	//first delete fail since no user exist
	err = db.DeleteUser(user.ID)
	require.Error(t, err)

	//create a user
	err = db.AddUser(user)
	require.NoError(t, err)

	//second delete do not fail
	err = db.DeleteUser(user.ID)
	require.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	db, err := storage.NewStorageDBInMemory("TestUpdate")
	require.NoError(t, err)

	defer db.Close()

	//cannot update, no user present
	user := getUserTest()
	err = db.UpdateUser(user)
	require.Error(t, err)

	//create a user
	err = db.AddUser(user)
	require.NoError(t, err)

	//list the user
	list, err := db.ListUser(models.ListOptions{})
	require.NoError(t, err)

	require.Len(t, list, 1)
	require.Equal(t, country, list[0].Country)

	//updating the user
	user.Country = "updated"
	err = db.UpdateUser(user)
	require.NoError(t, err)

	//list the user
	list, err = db.ListUser(models.ListOptions{})
	require.NoError(t, err)

	require.Len(t, list, 1)
	require.Equal(t, "updated", list[0].Country)

}

func getUserTest() models.User {
	return models.User{
		Country:   country,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		NickName:  nickname,
		ID:        uuid.New().String(),
	}
}
