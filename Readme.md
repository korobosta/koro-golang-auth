# What is koroauth

<p align="center">
<img src="https://m-shaeri.ir/blog/wp-content/uploads/2022/04/gologin.png"  height="200" >
</p>

**koroauth** is an easy to setup professional login manager for Go web applications. It helps you protect your application resources from unattended, unauthenticated or unauthorized access. Currently it works with SQL databases authentication. It is flexible, you can use it with any user/roles table structure in database.

## How to setup

Get the package with following command :

```bash
go get github.com/korobosta/koro-golang-auth

```

## How to use

You can easily setup and customize login process with **configure()** function. You should specify following paramters to make the koroauth ready to start:

- **Login page** : path to html template. Default path is ***./template/login.html***, note that the template must be defined as ****"login"**** with ***{{define "login"}}*** at the begining line

- **Login path** : login http path. Default path is ***/login***

- **Password Column Field** : Column name that represents the password in the database. Default  is ***password***

- **Usernane Column Field** : Column name that represents the username in the database. Default is ***username***

- **User Id Column Field** : Column name that represents the user id/primary key in the database. Default is ***id***

- **BycryptCost** : The package uses bycrypt for hashing password. A cost is needed for hashing.  Default is ***14***

- **Session timeout** : Number of seconds before the session expires. Default value is 3600 seconds.

- **Password encryption** : Password encryption function to hash the password before storing it in the database. Default is ***HashPassword*** from bycrypt

- **Compare Password Funcrion** : This function compares the hashed database password and the the password the user has entered. Returns either true or false. This function takes two parameters, the first being the hashed password from the database and the second is plain entered password by the user.  Default is ***CheckPasswordHash*** from bycrypt

- **SQL connection, and SQL query to authenticate user and fetch roles** : 2 SQL queries to retrieve user and its roles by given username and password. The authentication query must return only single arbitary column, it must have a where clause with two placeholder ::username and ::password. And the query for retrieving user's roles must return only the text column of role name.

- **Wrap desired endpoints to protect** : You should wrap the endpoints you want to protect with ***koroauth.LoginRequired*** or ***koroauth.RolesRequired*** function in the main function.( see the example)

***koroauth.LoginRequired*** requires user to be authenticated for accessing the wrapped endpoint/page.

***koroauth.RolesRequired*** requires user to have specified roles in addition to be authenticated.

See the example :

```Go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/korobosta/koroauth"
	_ "github.com/go-sql-driver/mysql"
)

// static assets like CSS and JavaScript
func public() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
}

// a page in our application, it needs user only be authenticated
func securedPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi! Welcome to secured page.")
	})
}

// another page in our application, it needs user be authenticated and have ADMIN role
func securedPage2() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi! Welcome to very secured page.")
	})
}


func main() {
	// create connection to database
	db, err := sql.Open("mysql", "root:12345@tcp(127.0.0.1:6666)/mydb")
	if err != nil {
		log.Fatal(err)
	}

	// koroauth configuration
	koroauth.Configure().
		SetLoginPage("./template/login.html"). // set login page html template path
		SetSessionTimeout(3600).                 // set session expiration time in seconds
		SetBycryptCost(14).                 // set cost for encrypting the password
		SetLoginPath("/login").                // set login http path
		SetUserTableName("users").                // Table name for users
		SetPasswordColumnName("password").                // Column name of the password field
		SetUsernameColumnName("username").                // Column name of the username field
		SetUserIdColumnName("id").                // Column name of the user id field
		UserTableColumns([]string{"first_name","middle_name","last_name","email","password","username"}). // Columns in the users table

	// instantiate http server
	mux := http.NewServeMux()

	mux.Handle("/static/", public())

	// use koroauth login handler for /login endpoint
	mux.Handle("/login", koroauth.LoginHandler())

	// the pages/endpoints that we need to protect should be wrapped with koroauth.LoginRequired
	mux.Handle("/mySecuredPage", koroauth.LoginRequired(securedPage()))

	mux.Handle("/mySecuredPage2", koroauth.RolesRequired(securedPage2()),"ADMIN")

	// server configuration
	addr := ":8080"
	server := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// start listening to network
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("main: couldn't start simple server: %v\n", err)
	}
}

```

Note that the form must send form data as post to the same url (set no action attribute).

Html template for login page :

```HTML
{{define "login"}}
<html>
    <body>
        <H2>
            Login Page
        </H2>
        <form method="post">
            <!-- username input with "username" name -->
            <input type="text" name="username" />
            <input type="password" name="password" />
            <input type="submit" value="Login" />
        </form>
    </body>
</html>

{{end}}

```

You can also store data in in-memory session storage in any type during user's session with **SetSession** function, and retrieve it back by **GetSession** function.

```Go
func securedPage2() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the session data, the request parameter is *http.Request
		age, err : = koroauth.GetSession("age", request)

		// as the GetSession returns type is interface{}, we should specify the exact type of the session entry
		fmt.Printf("Your age is " + age.(int))
	})
}
```
The default EncryptFunction for hashing password is HashPassword that uses bycrypt as shown below. The function takes the plain string passowrd and returns a string hashed password

```Go
import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err.Error())
	}
	return string(bytes)
}
```

You can provide your own EncryptFunction to hash the password as shown below

```Go
func MyOwnHashPasswordFunction(password string) string {
	hashed_password = // Hash password login
	return hashed password
}

// in main.go
koroauth.EncryptFunction(MyOwnHashPasswordFunction)

```

The default ComparePasswordFunction uses bycrypt as shown below
```Go

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

```
You can write your own ComparePasswordFunction as shown below
```Go

func MyOwnCheckPasswordHash(hash, password string) bool {
	var is_authenticated bool = false

	if(password == hash){
		is_authenticated = true
	}

	return is_authenticated
}

// in main.go
koroauth.ComparePasswordFunction(MyOwnHashPasswordFunction)

```



 And with **GetCurrentUsername** you can get the current user's username.

```Go
func securedPage2() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get the current user's username, the request parameter is *http.Request
			username : = koroauth.GetCurrentUsername(request)

			fmt.Printf("Welcome " + username)
	})
}
```

To logout users direct them to your **login url + ?logout=yes** for example if your login url is **/login** your application logout url will be **/login?logout=yes**


## Todo list

- Template session messaging
- Connection to mssql, postgreSQL
