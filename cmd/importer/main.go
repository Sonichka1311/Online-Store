package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"shop/pkg/constants"
	"shop/pkg/handlers"
	"strconv"
)

func main() {
	handler := &handlers.ImporterHandler{}
	handlerError := handler.Init()
	if handlerError != nil {
		log.Println("Failed to connect to upload queue")
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/upload", handler.ImportFile)

	err := http.ListenAndServe(":" + strconv.Itoa(constants.UploadPort), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
