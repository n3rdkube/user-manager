package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/n3rdkube/user-manager/internal/models"
	_ "github.com/proullon/ramsql/driver"
	"github.com/sirupsen/logrus"
)

const (
	mysqlDriver = "mysql"
	ramDriver   = "ramsql"

	defaultMaxLifetime = time.Minute * 3
	defaultMaxOpenConn = 10
	defaultMaxIdleConn = 10

	createTableQuery = `CREATE TABLE IF NOT EXISTS users(id VARCHAR(100) primary key not null, nickname text, country text,email text,  firstName text, password text, lastName text);`
)

//Storage interface is the contract that a struct should meet in order to be used as storage for the API
// At some point I can think of breaking down this interface down into several smaller interfaces
type Storage interface {
	AddUser(data models.User) error
	DeleteUser(id string) error
	UpdateUser(data models.User) error
	ListUser(data models.ListOptions) ([]models.User, error)
}

// StorageDB implements Storage interface
type StorageDB struct {
	*sql.DB
}

//NewStorageDBInMemory create a new instance of a temporary db as a using an in memory database, it is not the same, but it works to run
// small unit tests
func NewStorageDBInMemory(name string) (*StorageDB, error) {

	db, err := sql.Open(ramDriver, name)
	if err != nil {
		return nil, fmt.Errorf("creating mock db: %w", err)
	}

	err = createUsersTable(db)
	if err != nil {
		return nil, fmt.Errorf("creating usertable: %w", err)
	}

	return &StorageDB{db}, nil
}

//NewMysqlStorageDB create a new instance of the db with a mysql driver
func NewMysqlStorageDB(connURL string) (*StorageDB, error) {
	logrus.Info("Creating an instance of a mysql db to start opening connections")
	db, err := sql.Open(mysqlDriver, connURL)
	if err != nil {
		return nil, fmt.Errorf("opening sql connection: %w", err)
	}

	db.SetConnMaxLifetime(defaultMaxLifetime)
	db.SetMaxOpenConns(defaultMaxOpenConn)
	db.SetMaxIdleConns(defaultMaxIdleConn)

	err = createUsersTable(db)
	if err != nil {
		return nil, fmt.Errorf("creating usertable: %w", err)
	}

	return &StorageDB{db}, nil
}

//CreateUsersTable create the users tabled in the database
func createUsersTable(db *sql.DB) error {
	logrus.Info("Creating the table users")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning db transaction to create table: %w", err)
	}

	_, err = db.ExecContext(ctx, createTableQuery)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("creating table table: %w", err)
	}

	return tx.Commit()
}
