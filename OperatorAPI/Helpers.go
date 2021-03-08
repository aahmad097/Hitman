package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func hashpwd(pwd []byte) string {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(hash)

}

func comparePWD(hashedPwd string, pwd []byte) bool {

	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, pwd)
	if err != nil {

		return false

	}

	return true

}
