package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"shop/pkg/constants"
	"shop/pkg/database"
	"shop/pkg/handlers"
	"shop/pkg/sessions"
	"shop/pkg/user"
	"strconv"
	"time"
)

func main() {
	db, dbError := sql.Open("mysql", "root:guest@tcp(mysql:3306)/shop?charset=utf8&interpolateParams=true")
	if dbError != nil {
		log.Fatalf("Cannot open database: %s", dbError.Error())
	}

	for tries := 0; tries < 10; tries++ {
		dbError = db.Ping()
		if dbError == nil {
			break
		}
		log.Printf("Failed connect to database for %d times. Trying to reconnect...", tries + 1)
		time.Sleep(3 * time.Second)
	}
	if dbError != nil {
		log.Fatalf("Cannot connect to database: %s", dbError.Error())
	}

	notificationHandler := &handlers.NotificationHandler{}
	notificationError := notificationHandler.Init()
	if notificationError != nil {
		log.Println("Failed to connect to notification queue")
		return
	}

	handler := handlers.AuthHandler{
		Repo:   &user.Repo{
			Connector:  database.NewConnector(db),
		},
		Sessions: &sessions.Repo{
			Connector:	database.NewConnector(db),
		},
		Notifications: notificationHandler,
	}

	router := mux.NewRouter()

	router.HandleFunc("/signup", handler.SignUp)
	router.HandleFunc("/verify/{token}", handler.ConfirmRegister).Methods(http.MethodGet)
	router.HandleFunc("/signin", handler.SignIn)
	router.HandleFunc("/validate", handler.ValidateToken)
	router.HandleFunc("/refresh", handler.RefreshToken)

	err := http.ListenAndServe(":" + strconv.Itoa(constants.AuthPort), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
