package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/inancgumus/screen"
)

const (
	accountHelp = `
	
Session Menu:

[+] adduser <username> <pwd>              [List Available Sessions]
[+] chpwd <currentpwd> <pwd> <conf. pwd>  [Interact With Target Session]
[+] clear                                 [Clear Screen]
[+] back                                  [Back to Main Menu]
[+] help                                  [Print This Help Menu]`
)

func acctMenu() {

	var opt string
	fmt.Println(accountHelp + "\n")

	for true {

		fmt.Print(sess.username + "@" + sess.host + " > ")
		var r = bufio.NewReader(os.Stdin)
		fmt.Fscanf(r, "%s", &opt)

		if opt == "adduser" {

			opt = ""
			var username string
			var password string
			var c string
			var admin bool

			fmt.Fscanf(r, "%s", &username)
			fmt.Fscanf(r, "%s", &password)

			for true {

				fmt.Print("[?] Make user admin (y/n): ")
				var b = bufio.NewReader(os.Stdin)
				fmt.Fscanf(b, "%s", &c)

				if c == "y" {

					admin = true
					break

				} else if c == "n" {

					admin = false
					break

				} else {

					fmt.Println("[?] I didn't understand that!")

				}

			}

			addUser(username, password, admin)

		} else if opt == "chpwd" {

			opt = ""

			var cpwrd string
			var pwd string
			var confpwd string

			fmt.Fscanf(r, "%s", &cpwrd)
			fmt.Fscanf(r, "%s", &pwd)
			fmt.Fscanf(r, "%s", &confpwd)

			chpasswd(cpwrd, pwd, confpwd)

		} else if opt == "help" {

			opt = ""
			fmt.Println(accountHelp + "\n")

		} else if opt == "clear" {

			opt = ""
			screen.Clear()

		} else if opt == "back" {

			break

		} else {

			opt = ""
			fmt.Println("[!] Do you ever read the help menu?")

		}

	}

}

func addUser(username string, password string, admin bool) {

	var role string

	if admin {
		role = "admin"
	} else {

		role = "operator"

	}

	uri := sess.url + "/adduser"

	resp, _ := sess.client.PostForm(uri, url.Values{"username": {username}, "password": {password}, "role": {role}})
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(string(b))

}

func chpasswd(currentpwd string, password string, confpassword string) {

	uri := sess.url + "/changepassword"

	resp, _ := sess.client.PostForm(uri, url.Values{"oldpassword": {currentpwd}, "password": {password}, "confirmpassword": {confpassword}})
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(string(b))

}
