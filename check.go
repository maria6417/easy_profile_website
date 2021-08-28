package main

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func getCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	c, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		uuid := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: uuid.String(),
		}
		http.SetCookie(w, c)
	}
	return c
}

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	cookie := getCookie(w, r)
	uuid := cookie.Value
	currentUserId := sessionMap[uuid]
	if currentUserId == "" {
		return false
	} else {
		return true
	}
}

func verifyUser(w http.ResponseWriter, r *http.Request) bool {
	if doesUserExist(w, r) {
		// if user exists, check if password is correct
		pass := []byte(r.FormValue("password"))
		correctPass := userInfo[r.FormValue("username")].Password
		if checkPasswordHash(pass, correctPass) {
			// if is correct returns true
			return true
		} else {
			return false
		}
	} else {
		// if user doesnt exist, just return false.
		return false
	}
}

func doesUserExist(w http.ResponseWriter, r *http.Request) bool {
	username := r.FormValue("username")
	if _, ok := userInfo[username]; ok {
		// if there is corresponding record in userInfo with that username, returns true.
		return true
	} else {
		return false
	}
}

func addUser(w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	password, err := hashPassword(r.FormValue("password"))
	if err != nil {
		return err
	}
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	u1 := user{
		username, []byte(password), fname, lname,
	}
	userInfo[username] = u1
	return err
}

func hashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return bytes, err
}

func checkPasswordHash(password, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}
