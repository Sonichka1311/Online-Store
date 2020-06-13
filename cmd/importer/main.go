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
	"strconv"
	"time"
)

func main() {
	db, dbError := sql.Open("mysql", "root:guest@tcp(mysql:3306)/shop?charset=utf8&interpolateParams=true")
	if dbError != nil {
		log.Fatalf("Cannot open database: %s", dbError.Error())
	}

	for tries := 0; tries < constants.DatabaseConnectionRetries; tries++ {
		dbError = db.Ping()
		if dbError == nil {
			break
		}
		log.Printf("Failed connect to database for %d times. Trying to reconnect...", tries + 1)
		time.Sleep(constants.DatabaseConnectionSleepTime)
	}
	if dbError != nil {
		log.Fatalf("Cannot connect to database: %s", dbError.Error())
	}

	handler := &handlers.ImporterHandler{
		Database: database.NewConnector(db),
	}
	handlerError := handler.Init()
	if handlerError != nil {
		log.Println("Failed to connect to upload queue")
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/upload", handler.ImportFile)
	router.HandleFunc("/info", handler.GetStatus)

	err := http.ListenAndServe(":" + strconv.Itoa(constants.UploadPort), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
