package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"shop/pkg/constants"
	"shop/pkg/handlers"
	"shop/pkg/logic"
	"shop/pkg/models"
	"shop/pkg/sessions"
	"shop/pkg/user"
	"strconv"
	"sync"
)

func main() {
	databaseConnector := models.Connector{
		Router: models.Router{Host: logic.GetUrl(constants.Protocol, constants.DatabaseHost, constants.DatabasePort)},
		Mutex:  sync.Mutex{},
	}

	notificationHandler := &handlers.NotificationHandler{}
	notificationError := notificationHandler.Init()
	if notificationError != nil {
		log.Println("Failed to connect to notification queue")
		return
	}

	handler := handlers.AuthHandler{
		Repo:   &user.Repo{
			Connector:  &databaseConnector,
		},
		Sessions: &sessions.Repo{
			Connector:	&databaseConnector,
		},
		Notifications: notificationHandler,
	}

	router := mux.NewRouter()

	router.HandleFunc("/signup", handler.SignUp)
	router.HandleFunc("/{token}", handler.ConfirmRegister).Methods(http.MethodGet)
	router.HandleFunc("/signin", handler.SignIn)
	router.HandleFunc("/validate", handler.ValidateToken)
	router.HandleFunc("/refresh", handler.RefreshToken)

	err := http.ListenAndServe(":" + strconv.Itoa(constants.AuthPort), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
