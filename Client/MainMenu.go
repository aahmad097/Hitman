package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/inancgumus/screen"
)

const (
	mainmenu = `
	
Options:
	
[+] session [Session Interaction Menu]
[+] account [Account Menu]
[+] clear   [Clear Screen]
[+] exit    [Exit]
[+] help    [Print This Help Menu]`
)

func mmenu() {

	var opt string
	fmt.Println(mainmenu + "\n")

	for true {

		fmt.Print(sess.username + "@" + sess.host + " > ")
		var r = bufio.NewReader(os.Stdin)
		fmt.Fscanf(r, "%s", &opt)

		if opt == "session" {

			opt = ""
			sessionMenu()

		} else if opt == "account" {

			opt = ""
			acctMenu()

		} else if opt == "help" {

			opt = ""
			fmt.Println(mainmenu + "\n")

		} else if opt == "clear" {

			opt = ""
			screen.Clear()

		} else if opt == "exit" {

			os.Exit(0)

		} else {

			fmt.Println("[!] Do you ever read the help menu?")
			opt = ""

		}

	}

}
