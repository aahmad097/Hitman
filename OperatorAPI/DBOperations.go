package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type DBTask struct {
	TASKID    string
	SESSIONID string
	COMPLETE  bool
	TASK      string
	METHOD    string
	TARGET    string
	DATA      string
	RESPONSE  string
}

type Task struct {
	TASKID string
	TASK   string
	METHOD string
	TARGET string
	DATA   string
}

func conn(host string, port int, database string, username string, password string) *sql.DB {

	DB_DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, host, port, database)
	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		fmt.Println("Unable to connect to DB: ", err)
		os.Exit(1)
	}

	return db
}

func taskInserter(db *sql.DB, sessionid string, task *Task) bool {

	task.TASKID = fmt.Sprintf("%d", qRows(db, sessionid, true))
	serTask, err := json.Marshal(task)
	if err != nil {

		fmt.Println("Unable to serialize task")
		return false

	}
	dbTaskData := base64.StdEncoding.EncodeToString([]byte(serTask))

	sqlStatement := `INSERT INTO sessions.tasking(ID, SESSIONID, TASKID, COMPLETE, TASK) VALUES ($1, $2, $3, $4, $5)`

	_, dberr := db.Exec(sqlStatement, qRows(db, "", false), sessionid, task.TASKID, false, dbTaskData)
	if dberr != nil {
		fmt.Println("Database error: ", dberr)
		return false
	}

	return true
}

func qRows(db *sql.DB, sessionid string, sess bool) int {

	count := 0
	if sess {
		sqlStatement := "SELECT MAX(taskid) FROM sessions.tasking WHERE SESSIONID=$1"
		_ = db.QueryRow(sqlStatement, sessionid).Scan(&count)
		count++
	} else {
		_ = db.QueryRow("SELECT MAX(id) FROM sessions.tasking;").Scan(&count)
		count++
	}
	return count

}

func taskfetcher(db *sql.DB, sToken string) string {

	var response string
	sqlStatement := "SELECT task FROM sessions.tasking WHERE COMPLETE='f' AND SESSIONID=$1 ORDER BY ID LIMIT 1"
	err := db.QueryRow(sqlStatement, sToken).Scan(&response)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			fmt.Println("Session has no tasks queued")
			return response
		}

		fmt.Println("Failed to query task for session: ", err)
		return response
	}

	return response
}

func qSessions(db *sql.DB) []session {

	var sess session
	var sessions []session
	rows, err := db.Query("select SESSIONID, IMPLANTTYPE, COMPUTERNAME from sessions.sessions;")
	if err != nil {
		fmt.Println("[!] Unable to fetch sessions!")

	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&sess.Sessionid, &sess.Implanttype, &sess.Computername); err != nil {
			fmt.Println("[!] Cannot grab session info")
		}

		sessions = append(sessions, sess)
	}

	return sessions

}

func qTasks(db *sql.DB, sessionid string) []DBTask {

	var dbtask DBTask
	var tasks []DBTask

	rows, err := db.Query("select TASKID, COMPLETE from sessions.tasking where sessionid = $1", sessionid)
	if err != nil {

		fmt.Println("[!] Unable to fetch tasks! Error: ", err)
		return tasks

	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&dbtask.TASKID, &dbtask.COMPLETE); err != nil {
			fmt.Println("[!] Cannot grab task info")
			return tasks
		}

		tasks = append(tasks, dbtask)
	}

	return tasks

}

func qTask(db *sql.DB, sessionid string, task int) string {

	var response string

	rows, err := db.Query("select RESPONSE from sessions.tasking where sessionid=$1 and taskid = $2", sessionid, task)
	if err != nil {

		fmt.Println("[!] Unable to fetch sessions!")
		return ""
	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&response); err != nil {

			fmt.Println("[!] Cannot grab task info")
			return ""

		}

	}

	return response

}

func getHash(db *sql.DB, username string) string {

	tarhash := ""
	sqlStatement := "select PWDHASH from OPERATORS.OPERATORS where USERNAME = $1;"

	rows, err := db.Query(sqlStatement, username)
	if err != nil {

		fmt.Println("[!] Error fetching user")
		return ""

	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&tarhash); err != nil {

			fmt.Println("[!] Cannot fetch hash")
			return ""

		}

	}

	return tarhash

}

func setCookie(db *sql.DB, username string, cookie string, exp int64) {

	sqlStatement := "UPDATE OPERATORS.OPERATORS SET SESSIONHASH = $1, SESSIONEXP = $2 WHERE USERNAME = $3;"
	_, err := db.Exec(sqlStatement, cookie, exp, username)
	if err != nil {

		fmt.Println("[!] Unable to update session tasking response: ", err)
		return

	}

	return

}

func checksession(db *sql.DB, session string, time int64) bool {

	var t int64 = 0

	sqlStatement := "SELECT SESSIONEXP FROM OPERATORS.OPERATORS WHERE SESSIONHASH = $1;"
	rows, err := db.Query(sqlStatement, session)
	if err != nil {

		fmt.Println("[!] Error fetching session info")
		return false
	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&t); err != nil {

			fmt.Println("[!] Cannot fetch hash")
			return false
		}

	}

	if t > time {

		return true

	} else {

		return false

	}

}

func getHashbySession(db *sql.DB, session string) string {

	tarhash := ""
	sqlStatement := "select PWDHASH from OPERATORS.OPERATORS where SESSIONHASH = $1;"

	rows, err := db.Query(sqlStatement, session)
	if err != nil {

		fmt.Println("[!] Error fetching user")
		return ""
	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&tarhash); err != nil {

			fmt.Println("[!] Cannot fetch hash")
			return ""

		}

	}

	return tarhash

}

func getRoleBySession(db *sql.DB, session string) string {

	role := ""
	sqlStatement := "SELECT role FROM OPERATORS.OPERATORS WHERE sessionhash = $1;"

	rows, err := db.Query(sqlStatement, session)
	if err != nil {

		fmt.Println("[!] Error fetching user role")
		return ""
	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&role); err != nil {

			fmt.Println("[!] Cannot set role")
			return ""

		}

	}

	return role

}

func updateUserHash(db *sql.DB, session string, newhash string) bool {

	sqlStatement := "UPDATE OPERATORS.OPERATORS SET PWDHASH = $1 WHERE SESSIONHASH = $2;"
	_, err := db.Exec(sqlStatement, newhash, session)
	if err != nil {

		fmt.Println("[!] Unable to update password hash ", err)
		return false

	}

	return true

}

func addDBUser(db *sql.DB, username string, role string, pwdhash string) bool {

	count := 0

	sqlStatement := "SELECT MAX(id) FROM OPERATORS.OPERATORS;"
	_ = db.QueryRow(sqlStatement).Scan(&count)
	count++

	sqlStatement2 := "INSERT INTO OPERATORS.OPERATORS(ID, USERNAME, ROLE, PWDHASH) VALUES ($1, $2, $3, $4);"
	_, err := db.Exec(sqlStatement2, count, username, role, pwdhash)
	if err != nil {

		fmt.Println("[!] Unable to add user to databse ", err)
		return false

	}

	return true

}

func userExist(db *sql.DB, username string) bool {

	exist := false
	count := 0

	sqlStatement := "SELECT ID FROM OPERATORS.OPERATORS WHERE USERNAME=$1"
	_ = db.QueryRow(sqlStatement, username).Scan(&count)

	if count > 0 {
		exist = true
	}

	return exist

}
