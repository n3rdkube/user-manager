package storage

import (
	"fmt"

	"github.com/n3rdkube/user-manager/internal/models"
	"github.com/sirupsen/logrus"
)

const (
	insertTableQuery = `INSERT INTO users (id ,nickname, country, email, firstName, password, lastName) VALUES (?, ?, ?, ?, ?, ?, ?);`
	updateTableQuery = `UPDATE users SET nickname = ?, country = ?, email = ?, firstName = ?, password = ?, lastName = ? WHERE id = ?;`
	deleteQuery      = `DELETE FROM users WHERE id=?;`

	defaultPageNumber = 1
	defaultRowPerPage = 10
)

// createListQuery uses reflection to automatically generate the sql query
func createListQueryWithValues(listOptions models.ListOptions) (string, []interface{}) {
	query := "SELECT id, country, email, firstName, nickname, lastName from users "

	pageNumber := defaultPageNumber
	if listOptions.PageNumber > 0 {
		pageNumber = listOptions.PageNumber
	}

	rowsPerPage := defaultRowPerPage
	if listOptions.RowsPerPage != 0 {
		rowsPerPage = listOptions.RowsPerPage
	}

	query, values := addWhereClauses(listOptions, query)

	offset := (pageNumber - 1) * rowsPerPage
	query = fmt.Sprintf("%s ORDER BY email LIMIT %d OFFSET %d ;", query, rowsPerPage, offset)

	logrus.Infof("listOptions query generated '%s %s'", query, values)
	return query, values
}

// This is just a stub, not a real implementation
func addWhereClauses(listOption models.ListOptions, query string) (string, []interface{}) {
	var values []interface{}
	logicOperator := "WHERE"

	if listOption.Include.FirstName != "" {
		query = query + fmt.Sprintf("%s firstname = ? ", logicOperator)
		values = append(values, listOption.Include.FirstName)
		logicOperator = "AND"
	}

	if listOption.Include.Country != "" {
		query = query + fmt.Sprintf("%s country = ? ", logicOperator)
		values = append(values, listOption.Include.Country)
		logicOperator = "AND"
	}

	if listOption.Include.LastName != "" {
		query = query + fmt.Sprintf("%s lastName = ? ", logicOperator)
		values = append(values, listOption.Include.LastName)
	}
	return query, values
}
