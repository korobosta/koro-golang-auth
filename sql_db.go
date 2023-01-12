package koroauth

import (
	"database/sql"
	"fmt"
)

type DataBaseInterface interface {
	AuthenticateUser(username string) (bool, map[string]interface{})
	RetriveRoles(username string) (bool, []string)
}

type SqlDataBase struct {
	*sql.DB
	AuthenticationSqlQuery string
	RolesSqlQuery          string
}

func (db *SqlDataBase) AuthenticateUser(username string) (bool, map[string]interface{}) {

	var col_string string = ""
	values := make([]string, len(config.UserTableColumns))
	pointers := make([]interface{}, len(config.UserTableColumns))

	for i, _ := range values {
		pointers[i] = &values[i]
	}

	for _, v := range config.UserTableColumns {
		if col_string == "" {
			col_string = col_string + v
		} else {
			col_string = col_string + "," + v
		}
	}

	var query string = "SELECT " + col_string + " from " + config.UserTableName + " WHERE username = ?"

	err := db.QueryRow(query, username).Scan(pointers...)

	if err != nil {
		fmt.Printf(err.Error())
	}

	result := make(map[string]interface{})
	for i, val := range values {
		result[config.UserTableColumns[i]] = val
	}

	return (err == nil), result

}

func (db *SqlDataBase) RetriveRoles(username string) (bool, []string) {

	if db.RolesSqlQuery == "" {
		return true, []string{}
	}

	rows, err := db.Query(db.RolesSqlQuery, username)

	if err != nil {
		fmt.Printf(err.Error())
		return true, []string{}
	}

	var roles []string
	for rows.Next() {
		var role string
		err := rows.Scan(&role)
		if err != nil {
			fmt.Println(err)
			return true, []string{}
		}

		roles = append(roles, role)
	}

	return (err == nil), roles

}
