package main

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {

	uname := r.FormValue("username")
	pword := r.FormValue("password")

	hash := getHash(db, uname)
	if comparePWD(hash, []byte(pword)) {

		exptime := time.Now().Unix() + 3600
		r1 := rand.New(rand.NewSource(exptime))
		Cookie := sha1.Sum([]byte(string(rune(r1.Intn(100))) + hash))

		ret := "Welcome " + uname + "!"

		addCookie(w, "session", fmt.Sprintf("%x", Cookie), 60*time.Minute)
		setCookie(db, uname, fmt.Sprintf("%x", Cookie), exptime)

		fmt.Fprintf(w, ret)

	} else {

		fmt.Fprintf(w, "Inalid Login!")

	}

}

func addCookie(w http.ResponseWriter, name, value string, ttl time.Duration) {
	expire := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

// untested

func changePassword(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session")
	time := time.Now().Unix()

	if checksession(db, cookie.Value, time) {

		newpassword := r.FormValue("password")
		newpassword2 := r.FormValue("confirmpassword")
		oldpassword := r.FormValue("oldpassword")

		if newpassword == newpassword2 {
			if len(newpassword) >= 10 {
				hash := getHashbySession(db, cookie.Value)
				if comparePWD(hash, []byte(oldpassword)) {

					newhash := hashpwd([]byte(newpassword))

					if updateUserHash(db, cookie.Value, newhash) {

						fmt.Fprintf(w, "Successfully updated password")

					} else {

						fmt.Fprintf(w, "Unable to update password :( Check logs!")

					}

				} else {

					fmt.Fprintf(w, "Old password mismatch!")

				}
			} else {

				fmt.Fprintf(w, "New password needs to be 10 chars or more")

			}
		} else {

			fmt.Fprintf(w, "Passwords dont match")

		}
	}
}

func adduser(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session")
	time := time.Now().Unix()

	if checksession(db, cookie.Value, time) {

		if getRoleBySession(db, cookie.Value) == "admin" {

			username := r.FormValue("username")
			password := r.FormValue("password")
			role := r.FormValue("role")

			hash := hashpwd([]byte(password))
			if !userExist(db, username) {

				if addDBUser(db, username, role, hash) {

					fmt.Fprintf(w, "Successfully added user!")

				} else {

					fmt.Fprintf(w, "Unable to add user check logs!")

				}

			} else {

				fmt.Fprintf(w, "User with that username already exists!")

			}

		} else {

			fmt.Fprintf(w, "Find a priv esc bug and this would work!")

		}

	}
}
