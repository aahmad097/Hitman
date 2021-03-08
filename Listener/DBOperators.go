package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func conn(host string, port int, database string, username string, password string) *sql.DB {

	DB_DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, host, port, database)
	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		fmt.Println("Unable to connect to DB: ", err)
		os.Exit(1)
	}

	return db
}

func qRows(db *sql.DB) int {

	count := 0
	_ = db.QueryRow("SELECT MAX(id) FROM sessions.sessions;").Scan(&count)
	count++
	return count

}

func registerer(db *sql.DB, session Session) bool {

	sqlStatement := `INSERT INTO sessions.sessions(ID,SESSIONID, IMPLANTTYPE, IP, COMPUTERNAME, USERNAME, DOMAIN, ENCRYPTIONKEY) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.Exec(sqlStatement, session.Sessionid, session.Sessionhash, session.Implanttype, session.Ip, session.Compname, session.Username, session.Domain, session.Cryptkey)
	if err != nil {
		fmt.Println("Database error: ", err)
		return false
	}

	return true

}

func taskfetcher(db *sql.DB, sToken string) string {

	var response string
	sqlStatement := "SELECT task FROM sessions.tasking WHERE COMPLETE='f' AND SESSIONID=$1 ORDER BY ID LIMIT 1"
	err := db.QueryRow(sqlStatement, sToken).Scan(&response)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			fmt.Println("[!] Session ", sToken, " has no tasks queued")
			return response
		}

		fmt.Println("Failed to query task for session: ", err)
		return response
	}

	return response
}

func taskresponse(db *sql.DB, session string, taskid int, response string) {

	sqlStatement := "UPDATE sessions.tasking SET COMPLETE='true', RESPONSE = $3 WHERE SESSIONID = $1 AND TASKID = $2;"
	_, err := db.Exec(sqlStatement, session, taskid, response)
	if err != nil {

		fmt.Println("[!] Unable to update session tasking response: ", err)
		return

	}

	return

}
