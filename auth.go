package korologin

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// func GetUserName(request *http.Request) (userName string) {
// 	if cookie, err := request.Cookie("session"); err == nil {
// 		cookieValue := make(map[string]string)
// 		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
// 			userName = cookieValue["name"]
// 		}
// 	}
// 	return userName
// }

func GetCurrentUsername(request *http.Request) string {
	username, ok := GetSession(username_session_key, request)
	if ok {
		return username.(string)
	} else {
		log.Println("getCurrentUsername error")
		return ""
	}
}

func GetCurrentUserRoles(request *http.Request) []string {
	roles, ok := GetSession(roles_session_key, request)
	if ok {
		return roles.([]string)
	} else {
		log.Println("GetCurrentUserRoles error")
		return []string{}
	}
}

func HasRole(role string, request *http.Request) bool {
	roles, ok := GetSession(roles_session_key, request)
	if ok {
		if rolesContains(roles.([]string), role) > -1 {
			return true
		}
		return false
	} else {
		log.Println("HasRole error")
		return false
	}
}

func GetDataReturnedByAuthQuery(request *http.Request) interface{} {
	data, ok := GetSession(auth_query_result_session_key, request)
	if ok {
		return data
	} else {
		log.Println("GetDataReturnedByAuthQuery error")
		return nil
	}
}

func doLogin(response http.ResponseWriter, request *http.Request, db DataBaseInterface) {
	username := request.FormValue("username")
	password := request.FormValue("password")

	password = config.EncryptFunction(password)

	redirectPath := request.URL.Query().Get("redirect")
	redirectTarget := config.LoginPath + "?wrong=yes&redirect=" + redirectPath
	if username != "" && password != "" {

		ok, data := db.AuthenticateUser(username)
		match := CheckPasswordHash(password, data["password"])
		if ok {

			sessionId := generateSessionId(username)
			setSessionId(sessionId, response)
			setSessionBySessionId(sessionId, auth_query_result_session_key, data, request)
			setSessionBySessionId(sessionId, username_session_key, username, request)

			if config.SqlDataBaseModel.RolesSqlQuery != "" {
				_, roles := db.RetriveRoles(username)
				setSessionBySessionId(sessionId, roles_session_key, roles, request)
			}

			if redirectPath != "" {
				redirectTarget = redirectPath
			} else {
				redirectTarget = "/"
			}

		}

	}
	http.Redirect(response, request, redirectTarget, 302)
}

func loginView(response http.ResponseWriter, request *http.Request) {
	// cookie handling
	var Logintemplates = template.Must(template.ParseFiles(
		config.LoginPage,
	))
	err := Logintemplates.ExecuteTemplate(response, "login", nil)

	if err != nil {
		http.Error(response, fmt.Sprintf("login: couldn't parse template: %v", err), http.StatusInternalServerError)
		return
	}

}

func LoginHandler() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		if request.Method == "GET" {
			logout := request.URL.Query().Get("logout")
			if logout != "" {
				clearSession(response, request)
			}
			loginView(response, request)
		} else if request.Method == "POST" {
			var db DataBaseInterface
			if config.GetDBType() == "sql" {
				db = &config.SqlDataBaseModel
			}
			doLogin(response, request, db)
		}
	})
}

func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		id := GetSessionId(request)
		if id != "" {
			next.ServeHTTP(response, request)
		} else {

			http.Redirect(response, request, config.LoginPath+"?redirect="+request.URL.Path, 302)
		}

	})
}

func RolesRequired(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		id := GetSessionId(request)
		if id != "" {

			for _, role := range roles {

				if !HasRole(role, request) {
					response.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintln(response, "Unauthorized")
					return
				}
			}

			next.ServeHTTP(response, request)
		} else {

			http.Redirect(response, request, config.LoginPath+"?redirect="+request.URL.Path, 302)
		}

	})
}
