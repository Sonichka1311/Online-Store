package handlers

import (
	"bufio"
	"encoding/json"
	"github.com/streadway/amqp"
	"io"
	"log"
	"net/http"
	"shop/pkg/auth"
	"strings"

	//"shop/pkg/auth"
	"shop/pkg/constants"
	"shop/pkg/models"
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

func (h *ImporterHandler) SendUploadRequest(products []product.Product, token string) error {
	log.Printf("Trying to import products")

	// init request
	req := product.Uploading{
		Token:    token,
		Products: products,
	}

	// request to json
	jsonReq, jsonError := json.Marshal(req)
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
			Body:        jsonReq,
		},
	)
	if publishError != nil {
		log.Printf("Notification: SendRequest error: %s\n", publishError.Error())
		return publishError
	}

	return nil
}

func (h *ImporterHandler) ImportFile(w http.ResponseWriter, r *http.Request) {
	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	_, authError := auth.Verify(r.Header.Get("AccessToken"))
	if checker.CheckError(authError) {
		return
	}

	_ = r.ParseMultipartForm(10 << 30)
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("bad file: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Upload file: name=%s, size=%d", header.Filename, header.Size)

	defer file.Close()
	rr := bufio.NewReader(file)

	var products []product.Product
	flag := true
	_, err = rr.ReadString('\n')
	if err != nil {
		flag = false
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for flag {
		line, err := rr.ReadString('\n')
		line = strings.TrimSuffix(line, "\r")
		log.Printf("Read line: %s\n", line)

		if err == io.EOF {
			flag = false
		} else if err != nil {
			log.Printf("Bad line: %v\n", err)
			continue
		}

		var prod product.Product
		parseErr := prod.GetFromCsv(line)
		if parseErr != nil {
			log.Printf("Bad line: %v\n", err)
			continue
		}
		products = append(products, prod)
		log.Printf("Append product with id %d\n", prod.Id)

		if len(products) == 20 || (!flag && len(products) > 0) {
			log.Println("Send upload request")
			h.SendUploadRequest(products, r.Header.Get("AccessToken"))
			products = products[:0]
		}
	}
	log.Println("File uploaded")

	w.WriteHeader(http.StatusOK)
}
