package handlers

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"shop/pkg/constants"
	"shop/pkg/models"
	"shop/pkg/user"
	"time"
)

type NotificationHandler struct {
	Connector	*amqp.Connection
	Channel 	*amqp.Channel
}

func (h *NotificationHandler) Init() error {
	counter := 0
	for {
		counter++
		var connectionError error
		h.Connector, connectionError = amqp.Dial(constants.QueueNotificationServer)
		if connectionError == nil {
			break
		}
		if counter == constants.QueueConnectionRetries {
			return connectionError
		}
		log.Printf("Failed to connect to queue: %s. Retrying...", connectionError)
		time.Sleep(constants.QueueConnectionSleepTime)
	}
	log.Println("Connected to queue.")

	var channelError error
	h.Channel, channelError = h.Connector.Channel()
	if channelError != nil {
		log.Printf("Notifications: Init error: %s\n", channelError.Error())
		return channelError
	}

	exchangeError := h.Channel.ExchangeDeclare("notifications", "fanout", true, false, false, false, nil)
	if exchangeError != nil {
		log.Printf("Notifications: Init error: %s\n", exchangeError.Error())
		return exchangeError
	}
	return nil
}

func (h *NotificationHandler) InitQueue() (<-chan amqp.Delivery, error) {
	queue, declareError := h.Channel.QueueDeclare("emails", false, false, false, false, nil)
	if declareError != nil {
		log.Printf("Failed to create queue: %s", declareError)
		return nil, declareError
	}

	bindError := h.Channel.QueueBind(queue.Name, "#", "notifications", false, nil)
	if bindError != nil {
		log.Printf("Failed to bind queue: %s", bindError)
		return nil, bindError
	}

	return h.Channel.Consume(queue.Name, "", false, false, false, false, nil)
}

func (h *NotificationHandler) Close() {
	_ = h.Channel.Close()
	_ = h.Connector.Close()
}

func (h *NotificationHandler) SendConfirmationRequest(userData *user.User, token string) error {
	log.Printf("Trying to send confirmation request for user %s", userData.Email)
	notification := models.Notification{Email: userData.Email, Message: constants.ConfirmationMessage(token)}
	if len(userData.Phone) > 0 {
		notification.Phone = userData.Phone
	}
	jsonNotification, jsonError := json.Marshal(notification)
	if jsonError != nil {
		log.Printf("Notification: SendRequest error: %s\n", jsonError.Error())
		return jsonError
	}

	publishError := h.Channel.Publish(
		"notifications",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonNotification,
		},
	)
	if publishError != nil {
		log.Printf("Notification: SendRequest error: %s\n", publishError.Error())
		return publishError
	}
	return nil
}
