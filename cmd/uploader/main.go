package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"shop/pkg/auth"
	"shop/pkg/constants"
	"shop/pkg/database"
	"shop/pkg/handlers"
	"shop/pkg/product"
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

	messages, queueError := handler.InitQueue()
	if queueError != nil {
		return
	}

	for message := range messages {
		//log.Printf("Received a message: %s\n", message.Body)

		var request product.Uploading
		jsonError := json.Unmarshal(message.Body, &request)
		if jsonError != nil {
			log.Printf("Failed to parse request: %s\n", jsonError.Error())
			return
		}

		for _, prod := range request.Products {
			jsonProd, jsonErr := prod.GetJson()
			if jsonErr != nil {
				log.Printf("Got broken JSON: %v, ERROR: %s\n", prod, jsonErr.ErrorString)
				continue
			}
			req, reqError := http.NewRequest(
				"POST",
				constants.Protocol + "://" + constants.MainHost + ":" + strconv.Itoa(constants.MainPort) + "/product",
				bytes.NewBuffer(jsonProd),
			)
			if reqError != nil {
				log.Printf("Something went wrong with http client: %s\n", reqError.Error())
				return
			}
			req.Header.Set("accessToken", request.Token)
			req.Header.Set("Content-Type", "application/json")
			_, authError := http.DefaultClient.Do(req)
			if authError != nil {
				fmt.Printf("Can't add product %v: ERROR: %s\n", prod, authError.Error())
			}
		}
		go func() {
			//time.Sleep(time.Second * 20)
			usr, _ := auth.Verify(request.Token)
			//log.Println(usr.Email)
			//if checker.CheckError(authError) {
			//	return
			//}
			row := handler.Database.SelectOne("current, max", "uploads", "login = ?", usr.Email)
			var count int
			var all int
			err := row.Scan(&count, &all)
			if err != nil {
				log.Printf("ERROR: %s", err.Error())
			}
			//log.Println(all)
			handler.Database.Update("uploads", "current = ?", "login = ?", count+len(request.Products), usr.Email)
		}()
		//log.Println("Products array added")
		message.Ack(false)
	}
}
