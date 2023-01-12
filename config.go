package koroauth

import (
	"database/sql"
	"strings"
)

type Config struct {
	LoginPage               string
	SessionTimeout          int
	LoginPath               string
	SqlDataBaseModel        SqlDataBase
	EncryptFunction         EncryptFunction
	ComparePasswordFunction ComparePasswordFunction
	UserTableColumns        []string
	UserTableName           string
	PasswordColumnName      string
	UsernameColumnName      string
	UserIdColumnName        string
	BycryptCost      int
}
type EncryptFunction func(string) string

type ComparePasswordFunction func(string, string) bool

var config Config

func Configure() *Config {

	config.LoginPage = "./templates/login.html"
	config.SessionTimeout = 120
	config.LoginPath = "/login"
	config.ComparePasswordFunction = CheckPasswordHash
	config.EncryptFunction = HashPassword
	config.UserTableColumns = []string{"fName", "lName", "userId", "email", "password", "isActive", "username", "pfNumber"}
	config.UserTableName = "users"
	config.PasswordColumnName = "password"
	config.UsernameColumnName = "username"
	config.UserIdColumnName = "userId"
	config.BycryptCost =14
	return &config
}

func (config *Config) SetPasswordColumnName(passwordColumnName string) *Config {
	config.PasswordColumnName = passwordColumnName
	return config
}

func (config *Config) SetBycryptCost(bycryptCost int) *Config {
	config.BycryptCost = bycryptCost
	return config
}

func (config *Config) SetUsernameColumnName(usernameColumnName string) *Config {
	config.UsernameColumnName = usernameColumnName
	return config
}

func (config *Config) SetUserIdColumnName(userIdColumnName string) *Config {
	config.UserIdColumnName = userIdColumnName
	return config
}

func (config *Config) SetLoginPage(loginPage string) *Config {
	config.LoginPage = loginPage
	return config
}

func (config *Config) SetUserTableName(userTableName string) *Config {
	config.UserTableName = userTableName
	return config
}

func (config *Config) SetUserTableColumns(tableColumns []string) *Config {
	config.UserTableColumns = tableColumns
	return config
}

func (config *Config) SetSessionTimeout(sessionTimeout int) *Config {
	config.SessionTimeout = sessionTimeout
	return config
}

func (config *Config) SetLoginPath(loginPath string) *Config {
	config.LoginPath = loginPath
	return config
}

func (config *Config) SetComparePasswordFunction(ComparePasswordFunc ComparePasswordFunction) *Config {
	config.ComparePasswordFunction = ComparePasswordFunc
	return config
}

func (config *Config) SetPasswordEncryption(EncryptFunc EncryptFunction) *Config {
	config.EncryptFunction = EncryptFunc
	return config
}

func (config *Config) AuthenticateBySqlQuery(db *sql.DB, authenticateQuery string, rolesQuery string) *Config {

	authenticateQuery = strings.Replace(authenticateQuery, "::username", "?", 1)
	authenticateQuery = strings.Replace(authenticateQuery, "::password", "?", 1)

	rolesQuery = strings.Replace(rolesQuery, "::username", "?", 1)

	config.SqlDataBaseModel = SqlDataBase{db, authenticateQuery, rolesQuery}
	return config
}

func (config *Config) GetDBType() string {
	if config.SqlDataBaseModel != (SqlDataBase{}) {
		return "sql"
	}
	return "noDB"
}
