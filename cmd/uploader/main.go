package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"shop/pkg/constants"
	"shop/pkg/handlers"
	"shop/pkg/product"
	"strconv"
)

func main() {
	handler := &handlers.ImporterHandler{}
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
		log.Printf("Received a message: %s\n", message.Body)

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

		log.Println("Products array added")
		message.Ack(false)
	}
}
