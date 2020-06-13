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

		notification := models.Notification{}
		jsonError := json.Unmarshal(message.Body, &notification)
		if jsonError != nil {
			log.Println("Failed to parse message")
			return
		}

		emailOk := notifications.SendEmail(notification.Email, notification.Message)
		if emailOk {
			log.Println("Email sent")
			message.Ack(false)
		}

		if len(notification.Phone) > 0 {
			smsOk := notifications.SendSms(notification.Phone, notification.Message)
			if smsOk {
				log.Println("Sms sent")
				//message.Ack(false)
			}
		}
	}
}
