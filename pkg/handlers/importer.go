package handlers

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"shop/pkg/constants"
	"shop/pkg/product"
	"time"
)

type ImporterHandler struct {
	Connector	*amqp.Connection
	Channel 	*amqp.Channel
}

func (h *ImporterHandler) Init() error {

	// connect to rabbit.mq
	for tries := 0; tries < constants.QueueConnectionRetries; tries++ {
		var connectionError error
		h.Connector, connectionError = amqp.Dial(constants.QueueUploadServer)
		if connectionError == nil {
			break
		}
		log.Printf("Failed to connect to queue: %s. Retrying...", connectionError)
		time.Sleep(constants.QueueConnectionSleepTime)
	}
	log.Println("Connected to queue.")

	// create channel
	var channelError error
	h.Channel, channelError = h.Connector.Channel()
	if channelError != nil {
		log.Printf("Notifications: Init error: %s\n", channelError.Error())
		return channelError
	}

	// create exchange
	exchangeError := h.Channel.ExchangeDeclare("upload", "fanout", true, false, false, false, nil)
	if exchangeError != nil {
		log.Printf("Notifications: Init error: %s\n", exchangeError.Error())
		return exchangeError
	}
	return nil
}

func (h *ImporterHandler) InitQueue() (<-chan amqp.Delivery, error) {
	// create queue
	queue, declareError := h.Channel.QueueDeclare("products", false, false, false, false, nil)
	if declareError != nil {
		log.Printf("Failed to create queue: %s", declareError)
		return nil, declareError
	}

	// .. dont know what is it actually... ..
	bindError := h.Channel.QueueBind(queue.Name, "#", "upload", false, nil)
	if bindError != nil {
		log.Printf("Failed to bind queue: %s", bindError)
		return nil, bindError
	}

	return h.Channel.Consume(queue.Name, "", false, false, false, false, nil)
}

func (h *ImporterHandler) Close() {
	_ = h.Channel.Close()
	_ = h.Connector.Close()
}

func (h *ImporterHandler) SendUploadRequest(products []product.Product) error {
	log.Printf("Trying to import products")

	// array to json
	jsonArray, jsonError := json.Marshal(products)
	if jsonError != nil {
		log.Printf("Notification: SendRequest error: %s\n", jsonError.Error())
		return jsonError
	}

	// add to queue
	publishError := h.Channel.Publish(
		"upload",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonArray,
		},
	)
	if publishError != nil {
		log.Printf("Notification: SendRequest error: %s\n", publishError.Error())
		return publishError
	}

	return nil
}
