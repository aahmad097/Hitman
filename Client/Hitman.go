package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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
                                 
	- Client - 

`

	usage = `
 Required Params:

 -host          Faceoff Operator API Host
 -port          Faceoff Operator API Port (default 80)
 -u             Username to use
 -p             Password to use

 Optional Params:
 
 -h             Print this help menu
 -ssl           Use ssl to interact with the server
	`
)

type FlagOptions struct { // option var decleration

	host     string
	port     int
	username string
	password string

	help bool
	ssl  bool
}

func options() *FlagOptions {

	host := flag.String("host", "", "Faceoff Operator API Host")
	port := flag.Int("port", 80, "Faceoff Operator API Port")
	username := flag.String("u", "", "Operator Account Username")
	password := flag.String("p", "", "Operator Account Password")

	help := flag.Bool("h", false, "Help Menu")
	ssl := flag.Bool("ssl", false, "Use ssl to interact with the server")

	flag.Parse()

	return &FlagOptions{

		host:     *host,
		port:     *port,
		username: *username,
		password: *password,

		help: *help,
		ssl:  *ssl,
	}

}

type session struct {
	url      string
	host     string
	username string
	jar      *Jar
	client   http.Client
}

var sess session

func main() {

	fmt.Println(Art)
	opt := options()
	if opt.help || opt.host == "" || opt.username == "" || opt.password == "" {

		fmt.Println(usage)
		os.Exit(0)

	}

	if !auth(opt) {

		log.Fatal("[!] Error Signing in")

	}

	sess.username = opt.username
	sess.host = opt.host

	mmenu()

}
