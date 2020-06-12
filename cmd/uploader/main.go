package main

import (
	"encoding/json"
	"fmt"
	"log"
	"shop/pkg/constants"
	"shop/pkg/handlers"
	"shop/pkg/logic"
	"shop/pkg/models"
	"shop/pkg/product"
	"sync"
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

	connector := models.Connector{
		Router: models.Router{Host: logic.GetUrl(constants.Protocol, constants.MainHost, constants.MainPort)},
		Mutex:  sync.Mutex{},
	}

	for message := range messages {
		log.Printf("Received a message: %s\n", message.Body)

		var products []product.Product
		jsonError := json.Unmarshal(message.Body, &products)
		if jsonError != nil {
			log.Println("Failed to parse products")
			return
		}

		for _, prod := range products {
			jsonProd, jsonErr := prod.GetJson()
			if jsonErr != nil {
				log.Printf("Got broken JSON: %v, ERROR: %s", prod, jsonErr.ErrorString)
				continue
			}
			_, authError := r.Connector.Post("validate", &jsonToken)
			if authError != nil {
				fmt.Println(authError.ErrorString)
				return authError
			}
		}

		//ok := notifications.SendEmail(notification.Email, notification.Message)

		//if ok {
		//	//log.Println("Email sent")
		//	message.Ack(false)
		//}
	}
}
