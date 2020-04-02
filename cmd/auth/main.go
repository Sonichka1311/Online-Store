package main

import (
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

	http.HandleFunc("/signup", handler.SignUp)
	http.HandleFunc("/signin", handler.SignIn)
	http.HandleFunc("/validate", handler.ValidateToken)
	http.HandleFunc("/refresh", handler.RefreshToken)

	err := http.ListenAndServe(":" + strconv.Itoa(constants.AuthPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
