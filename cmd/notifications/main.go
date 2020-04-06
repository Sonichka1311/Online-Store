package main

import (
	"encoding/json"
	"log"
	"shop/pkg/handlers"
	"shop/pkg/models"
	"shop/pkg/notifications"
)

func main() {
	handler := &handlers.NotificationHandler{}
	handlerError := handler.Init()
	if handlerError != nil {
		log.Println("Failed to connect to notification queue")
		return
	}

	messages, queueError := handler.InitQueue()
	if queueError != nil {
		return
	}

	for message := range messages {
		log.Printf("Received a message: %s\n", message.Body)

		notification := models.EmailNotification{}
		jsonError := json.Unmarshal(message.Body, &notification)
		if jsonError != nil {
			log.Println("Failed to parse message")
			return
		}

		ok := notifications.SendEmail(notification.Email, notification.Message)

		if ok {
			log.Println("Email sent")
			message.Ack(false)
		}
	}
}
