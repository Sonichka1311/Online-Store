package main

import (
	"../common/constants"
	"./handlers"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/signup", handlers.SignUp)
	http.HandleFunc("/signin", handlers.SignIn)
	http.HandleFunc("/validate", handlers.ValidateToken)
	http.HandleFunc("/refresh", handlers.RefreshToken)

	err := http.ListenAndServe(":" + strconv.Itoa(constants.AuthPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
