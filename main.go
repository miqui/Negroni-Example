package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
	"log"
	"net/http"
    "encoding/json"
)

var db *sql.DB

func main() {

	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/negroni")
	//err := db.Connect()
	//errHandler(err)
	db.Ping()
	defer db.Close()

	mux := http.NewServeMux()
	n := negroni.Classic()

	store := sessions.NewCookieStore([]byte("ohhhsooosecret"))
	n.Use(sessions.Sessions("gloabl_session_store", store))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		SimplePage(w, r, "mainpage")
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			SimplePage(w, r, "login")
		} else if r.Method == "POST" {
			LoginPost(w, r)
		}
	})

	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			SimplePage(w, r, "signup")
		} else if r.Method == "POST" {
			SignupPost(w, r)
		}
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		SimplePage(w, r, "logout")
	})

	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		SimpleAuthenticatedPage(w, r, "home")
	})

        mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
                APIHandler(w, r)
        })


	n.UseHandler(mux)
	n.Run(":3000")

}

func errHandler(err error) {
	if err != nil {
		log.Print(err)
	}
}

func SimplePage(w http.ResponseWriter, req *http.Request, template string) {

	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)

}

func SimpleAuthenticatedPage(w http.ResponseWriter, req *http.Request, template string) {

	session := sessions.GetSession(req)
	sess := session.Get("hello")

	if sess == nil {
		http.Redirect(w, req, "/notauthenticated", 301)
	}

	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)

}

func LoginPost(w http.ResponseWriter, req *http.Request) {

	session := sessions.GetSession(req)
	session.Set("hello", "world")

	username := req.FormValue("username")
	password := req.FormValue("password")

	var (
		email string
	)

	err := db.QueryRow("SELECT email FROM users WHERE username = ? AND password = ?", username, password).Scan(&email)
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, req, "/failedquery", 301)
	}

	fmt.Println(email)

	//r := render.New(render.Options{})
	//r.HTML(w, http.StatusOK, "home", nil)

	http.Redirect(w, req, "/home", 302)

}

func SignupPost(w http.ResponseWriter, req *http.Request) {

	username := req.FormValue("username")
	password := req.FormValue("password")
	email := req.FormValue("email")

	db.Exec("INSERT INTO users (username, password, email) VAUES ('?', '?', '?')", username, email, password)

	http.Redirect(w, req, "/login", 302)

}




func APIHandler(w http.ResponseWriter, req *http.Request) {

    data, _ := json.Marshal("{'API Test':'works!'}")
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Write(data)

}

