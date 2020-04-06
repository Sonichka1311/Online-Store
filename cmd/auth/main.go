package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"shop/pkg/constants"
	"shop/pkg/handlers"
	"shop/pkg/logic"
	"shop/pkg/models"
	"shop/pkg/user"
	"strconv"
	"sync"
)

func main() {
	databaseConnector := models.Connector{
		Router: models.Router{Host: logic.GetUrl(constants.Protocol, constants.DatabaseHost, constants.DatabasePort)},
		Mutex:  sync.Mutex{},
	}

	handler := handlers.AuthHandler{
		Repo:   &user.Repo{
			Connector:  &databaseConnector,
		},
	}

	router := mux.NewRouter()

	router.HandleFunc("/signup", handler.SignUp)
	router.HandleFunc("/{token}", handler.SignUp).Methods(http.MethodGet)
	router.HandleFunc("/signin", handler.SignIn)
	router.HandleFunc("/validate", handler.ValidateToken)
	router.HandleFunc("/refresh", handler.RefreshToken)

	err := http.ListenAndServe(":" + strconv.Itoa(constants.AuthPort), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
