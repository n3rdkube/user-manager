package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/n3rdkube/user-manager/internal/models"
	"github.com/sirupsen/logrus"
)

const (
	defaultTimeout = 5 * time.Second
)

//Close cleanUp resources
func (db StorageDB) Close() error {
	return db.DB.Close()
}

//AddUser add the user to the database
func (db *StorageDB) AddUser(data models.User) error {
	logrus.Infof("adding user to users table %q", data)

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction to create user in db: %w", err)
	}

	h := sha256.Sum256([]byte(data.Password))
	argsQuery := []interface{}{data.ID, data.NickName, data.Country, data.Email, data.FirstName, hex.EncodeToString(h[:]), data.LastName}
	_, err = tx.ExecContext(ctx, insertTableQuery, argsQuery...)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("adding user in db: %w", err)
	}

	return tx.Commit()
}

//DeleteUser delete the user from the database
func (db *StorageDB) DeleteUser(id string) error {
	logrus.Info("deleting user from users table")

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction to deleting user in db: %w", err)
	}

	r, err := tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("deleting user from db: %w", err)
	}

	rows, _ := r.RowsAffected()
	if rows == 0 {
		_ = tx.Rollback()
		return errors.New("no change made into the DB, likely userId was not matching any entry")
	}

	return tx.Commit()
}

//UpdateUser updates a user from the database
func (db *StorageDB) UpdateUser(data models.User) error {
	logrus.Info("updating user in users table")

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transhaction to update user: %w", err)
	}

	h := sha256.Sum256([]byte(data.Password))
	argsQuery := []interface{}{data.NickName, data.Country, data.Email, data.FirstName, hex.EncodeToString(h[:]), data.LastName, data.ID}
	r, err := tx.ExecContext(ctx, updateTableQuery, argsQuery...)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("updating users in DB: %w", err)
	}

	rows, _ := r.RowsAffected()
	if rows == 0 {
		_ = tx.Rollback()
		return errors.New("no change made into the DB, likely userId was not matching any entry")
	}

	return tx.Commit()
}

//ListUser list users using the list options provided, with no listOptions, all users are returned
func (db *StorageDB) ListUser(data models.ListOptions) ([]models.User, error) {
	logrus.Info("fetching user data from users table")

	list := []models.User{}
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("beginning transhaction to list users: %w", err)
	}

	listQuery, argsQuery := createListQueryWithValues(data)
	r, err := tx.Query(listQuery, argsQuery...)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("listing users from DB: %w", err)
	}
	defer r.Close()

	list, err = scanRows(r, list)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("scanning rows: %w", err)
	}

	return list, tx.Commit()
}

//scanRows check each row returned by the sql query
func scanRows(r *sql.Rows, list []models.User) ([]models.User, error) {
	for r.Next() {
		user := models.User{}
		err := r.Scan(&user.ID, &user.Country, &user.Email, &user.FirstName, &user.NickName, &user.LastName)
		if err != nil {
			return nil, fmt.Errorf("fetching data from one row: %w", err)
		}

		list = append(list, user)
	}

	return list, nil
}
