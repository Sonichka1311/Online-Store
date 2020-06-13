package handlers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"log"
	"net/http"
	"shop/pkg/auth"
	"shop/pkg/constants"
	"shop/pkg/database"
	"shop/pkg/models"
	"shop/pkg/product"
	"time"
)

type ImporterHandler struct {
	Database    *database.Connector
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
	//log.Printf("Trying to import products")

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

	// check rights
	usr, authError := auth.Verify(r.Header.Get("AccessToken"))
	if checker.CheckError(authError) {
		return
	}

	// get file type
	fileType := r.URL.Query().Get("type")

	// get file
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

	// init array
	var products []product.Product
	flag := true

	// skip first row
	_, err = rr.ReadString('\n')
	if err != nil {
		flag = false
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// update upload info
	_, err = h.Database.Insert("uploads", "login", "?", usr.Email)
	if err != nil {
		log.Println("ERROR INSERT: " + err.Error())
		_, err = h.Database.Update("uploads", "current = ?, max = ?, enough = ?", "login = ?", 0, 0, false, usr.Email)
		if err != nil {
			log.Println("ERROR INSERT OR UPDATE: " + err.Error())
		}
	}

	// find end pattern for xml files
	var end string
	if fileType == "xml" {
		begin, err := rr.ReadString('\n')
		if err == io.EOF {
			flag = false
		}
		end = begin[:1] + "/" + begin[1:]
	}
	//log.Printf("END: %s\n", end)

	// parse file
	counter := 0
	for flag {
		// read row
		line, err := rr.ReadString('\n')
		//log.Println("LINE: " + line)
		if err == io.EOF {
			flag = false
		} else if err != nil {
			log.Printf("Bad line: %v\n", err)
			continue
		}

		// declare product
		var prod product.Product
		var parseErr *models.Error

		// switch for file type
		if fileType == "csv" {
			parseErr = prod.GetFromCsv(line)
		} else if fileType == "xml" {
			// check if now is end pattern
			if line == end {
				flag = false
			}

			// get request if it isn't end
			if flag {
				// read next 3 (fields number) rows
				var lines string
				for i := 0; i < 3; i++ {
					line, _ := rr.ReadString('\n')
					//log.Println("New line: " + line)
					if err != nil {
						break
					}
					lines += line
				}
				//log.Println("Lines: " + lines)
				// skip row
				rr.ReadString('\n')

				// parse into product
				parseErr = prod.GetFromXML(lines)
			}
		} else {
			log.Printf("Unknown format: %s\n", fileType)
			newError, _ := models.NewError(errors.New(fmt.Sprintf("Unknown format: %s\n", fileType)), http.StatusBadRequest)
			replier.ReplyWithError(newError)
			return
		}

		if parseErr != nil {
			log.Printf("Bad line: %v\n", err)
			continue
		}
		if fileType != "xml" || flag {
			products = append(products, prod)
		}
		//log.Printf("Append product with id %d\n", prod.Id)

		if len(products) == 20 || (!flag && len(products) > 0) {
			//log.Println("Send upload request")
			counter += len(products)
			h.Database.Update("uploads", "max = ?", "login = ?", counter, usr.Email)
			h.SendUploadRequest(products, r.Header.Get("AccessToken"))
			products = products[:0]
		}
	}

	_, err = h.Database.Update("uploads", "max = ?, enough = ?", "login = ?", counter, true, usr.Email)
	if err != nil {
		log.Println("ERROR UPDATE: " + err.Error())
	}
	log.Printf("File uploaded, all count: %d\n", counter)

	w.WriteHeader(http.StatusOK)
}

func (h *ImporterHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	usr, authError := auth.Verify(r.Header.Get("AccessToken"))
	if checker.CheckError(authError) {
		return
	}

	row := h.Database.SelectOne("current, max, enough", "uploads", "login = ?", usr.Email)
	var count int
	var all int
	var enough bool
	err := row.Scan(&count, &all, &enough)
	if err != nil {
		log.Printf("ERROR GET STATUS: %s\n", err.Error())
		return
	}

	var message string

	if enough && count >= all{
		message = "Import and download is completed"
	} else if enough {
		message = fmt.Sprintf("Import completed. Download is not completed: downloaded %d, imported %d\n", count, all)
	} else {
		message = fmt.Sprintf("Import and download are not completed: downloaded %d, imported %d\n", count, all)
	}

	replier.ReplyWithMessage(message)
}
