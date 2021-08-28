package main

import (
	"html/template"
	"io"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

type user struct {
	UserName string
	Password []byte
	Fname    string
	Lname    string
	Admin    bool
}

var sessionMap = map[string]string{} //key: uuid, value: username
var userInfo = map[string]user{}     //key: username, value: userinfo

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/main", mainPage)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/main", http.StatusSeeOther)
		return
	}
	// go to login page and make them login.
	err := tpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/main", http.StatusSeeOther)
		return
	}
	// go to signup page and make them login.
	err := tpl.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(w, r)
	submitType := r.FormValue("submitType")
	if r.Method == http.MethodPost && submitType == "login" {
		// check if username and password is correct
		if !verifyUser(w, r) {
			err := tpl.ExecuteTemplate(w, "login.html", "username or password is not correct.")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		// register sessionId to sessionMap form value to userInfo
		username := r.FormValue("username")
		sessionMap[cookie.Value] = username
	}
	if r.Method == http.MethodPost && submitType == "signup" {
		// check is username is already taken
		if doesUserExist(w, r) {
			// user exists so cannot signup.
			err := tpl.ExecuteTemplate(w, "signup.html", "username already taken.")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		err := addUser(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		username := r.FormValue("username")
		sessionMap[cookie.Value] = username
	}

	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	username := sessionMap[cookie.Value]
	userData := userInfo[username]
	err := tpl.ExecuteTemplate(w, "main.html", userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func bar(w http.ResponseWriter, r *http.Request) {
	// bar can only be accessed by admins.
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	cookie := getCookie(w, r)
	currentUserId := sessionMap[cookie.Value]
	userData := userInfo[currentUserId]
	if !userData.Admin {
		io.WriteString(w, "Access Not Permitted. You have to be Administrator to access bar.")
		return
	}
	tpl.ExecuteTemplate(w, "bar.html", userData)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	cookie := getCookie(w, r)
	delete(sessionMap, cookie.Value)

	err := tpl.ExecuteTemplate(w, "logout.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
