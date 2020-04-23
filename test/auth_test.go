package test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"shop/pkg/constants"
	"shop/pkg/logic"
	"strconv"
	"strings"
	"testing"
)

type AuthTestCase struct {
	Handler 		string
	UserJson 		string
	Message			string
	GetTokens		bool
	GetMessage		bool
}

type TestUser struct {
	Email 		string 	`json:"email,omitempty"`
	Password 	string 	`json:"password,omitempty"`
	Confirm 	bool	`json:"confirm,omitempty"`
}

type TestConfirmation struct {
	Email 	string `json:"email,omitempty"`
	Token 	string `json:"token,omitempty"`
	Expire  string `json:"expire,omitempty"`
}

type TestSession struct {
	Email 	string `json:"email,omitempty"`
	Token 	string `json:"refresh_token,omitempty"`
	Expire  string `json:"expire,omitempty"`
}

type TestDB struct {
	Users				map[string]TestUser 		// email to user
	UserToConfirmation 	map[string]TestConfirmation // email to confirmation
	TokenToConfirmation map[string]TestConfirmation // confirm token to confirmation
	UserToSession	 	map[string]TestSession		// email to session
	TokenToSession		map[string]TestSession		// access token to session
}

func (db *TestDB) GetUser(email string) ([]byte, int) {
	user, exists := database.Users[email]
	if !exists {
		return nil, http.StatusNotFound
	}
	user.Email = ""
	jsonUser, _ := json.Marshal(user)
	return jsonUser, http.StatusOK
}

func (db *TestDB) AddUser(u TestUser) int {
	if u.Email == "" || u.Password == "" {
		return http.StatusBadRequest
	}
	if _, exists := db.Users[u.Email]; exists {
		if db.Users[u.Email].Confirm {
			return http.StatusOK
		}
		return http.StatusConflict
	}
	db.Users[u.Email] = u
	return http.StatusOK
}

func (db *TestDB) UpdateUser(u TestUser) int {
	if u.Email == "" {
		return http.StatusBadRequest
	}
	if _, exists := db.Users[u.Email]; !exists {
		return http.StatusNotFound
	}
	db.Users[u.Email] = u
	return http.StatusOK
}

func (db *TestDB) GetSession(token string) ([]byte, int) {
	session, exists := database.TokenToSession[token]
	if !exists {
		return nil, http.StatusNotFound
	}
	session.Token = ""
	jsonUser, _ := json.Marshal(session)
	return jsonUser, http.StatusOK
}

func (db *TestDB) AddSession(s TestSession) int {
	if _, exists := db.UserToSession[s.Email]; !exists {
		return http.StatusNotFound
	}
	db.UserToSession[s.Email] = s
	db.TokenToSession[s.Token] = s
	return http.StatusOK
}

func (db *TestDB) DeleteSession(s TestSession) int {
	if session, exists := db.UserToSession[s.Email]; exists && session.Token == s.Token{
		delete(db.UserToSession, s.Email)
		delete(db.TokenToSession, s.Token)
	}
	return http.StatusOK
}

func (db *TestDB) GetConfirmation(token string) ([]byte, int) {
	session, exists := database.TokenToConfirmation[token]
	if !exists {
		return nil, http.StatusNotFound
	}
	session.Token = ""
	jsonUser, _ := json.Marshal(session)
	return jsonUser, http.StatusOK
}

func (db *TestDB) AddConfirmation(c TestConfirmation) int {
	if _, exists := db.UserToSession[c.Email]; exists {
		return http.StatusConflict
	}
	db.UserToConfirmation[c.Email] = c
	db.TokenToConfirmation[c.Token] = c
	return http.StatusOK
}

func (db *TestDB) DeleteConfirmation(c TestConfirmation) int {
	if _, exists := db.UserToSession[c.Email]; exists {
		delete(db.UserToSession, c.Email)
		delete(db.TokenToSession, c.Token)
	}
	return http.StatusOK
}


var database = TestDB{}

func DatabaseServer(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	query := strings.Split(path, "/")
	switch query[0] {
	case "user":
		switch r.Method {
		case http.MethodGet:
			user, st := database.GetUser(query[1])
			if st != http.StatusOK {
				http.Error(w, "", st)
			} else {
				w.Write(user)
			}
		case http.MethodPost:
			var user TestUser
			defer r.Body.Close()
			readBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(readBody, user)
			user.Confirm = false
			status := database.AddUser(user)
			if status != http.StatusOK {
				http.Error(w, "", status)
			}
		case http.MethodPut:
			var user TestUser
			defer r.Body.Close()
			readBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(readBody, user)
			user.Confirm = true
			status := database.UpdateUser(user)
			if status != http.StatusOK {
				http.Error(w, "", status)
			}
		}
	case "session":
		switch r.Method {
		case http.MethodGet:
			user, st := database.GetSession(query[1])
			if st != http.StatusOK {
				http.Error(w, "", st)
			} else {
				w.Write(user)
			}
		case http.MethodPost:
			var session TestSession
			defer r.Body.Close()
			readBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(readBody, session)
			status := database.AddSession(session)
			if status != http.StatusOK {
				http.Error(w, "", status)
			}
		case http.MethodDelete:
			var session TestSession
			defer r.Body.Close()
			readBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(readBody, session)
			status := database.DeleteSession(session)
			if status != http.StatusOK {
				http.Error(w, "", status)
			}
		}
	case "confirmation":
		switch r.Method {
		case http.MethodGet:
			user, st := database.GetConfirmation(query[1])
			if st != http.StatusOK {
				http.Error(w, "", st)
			} else {
				w.Write(user)
			}
		case http.MethodPost:
			var conf TestConfirmation
			defer r.Body.Close()
			readBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(readBody, conf)
			status := database.AddConfirmation(conf)
			if status != http.StatusOK {
				http.Error(w, "", status)
			}
		case http.MethodDelete:
			var conf TestConfirmation
			defer r.Body.Close()
			readBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(readBody, conf)
			status := database.DeleteConfirmation(conf)
			if status != http.StatusOK {
				http.Error(w, "", status)
			}
		}
	}
}

var message string

func SmtpMock(w http.ResponseWriter, r *http.Request) {
	listen, _ := net.Listen("tcp", ":"+strconv.Itoa(constants.MockPort))

	for {
		conn, _ := listen.Accept()

		reader := bufio.NewReader(conn)
		_, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		_, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		_, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		_, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		_, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		ans := ""
		for {
			buf, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}
			msg := strings.TrimRight(buf, "\r\n")
			ans += msg
			if msg == "." {
				log.Println("mock: exit circle")
				break
			}
		}

		message = ans
		conn.Close()
	}
}

func TestOK(t *testing.T) {
	url := func(handler string) string { return fmt.Sprintf("http://localhost:3687/%s", handler) }
	cases := []AuthTestCase{
		{
			Handler: 			`signup`,
			UserJson:   		`{"email": "test@mail.ru","password": "testpassword"}`,
			Message:    		`{"message": "Confirmation request was sent to test@mail.ru\n"}`,
			GetTokens:      	false,
			GetMessage:     	true,
		},
		{
			Handler: 			`verify`,
			UserJson:   		``,
			Message:    		`{"message": "User test@mail.ru has been verified\n"}`,
			GetTokens:      	false,
			GetMessage:     	true,
		},
		{
			Handler: 			`signin`,
			UserJson:   		`{"email": "test@mail.ru","password": "testpassword"}`,
			Message:    		``,
			GetTokens:      	true,
			GetMessage:     	false,
		},
		{
			Handler: 			`validate`,
			UserJson:   		``,
			Message:    		`{"message": "Access token is valid\n"}`,
			GetTokens:      	false,
			GetMessage:     	true,
		},
		{
			Handler: 			`refresh`,
			UserJson:   		``,
			Message:    		``,
			GetTokens:      	true,
			GetMessage:     	false,
		},
	}

	db := httptest.NewServer(http.HandlerFunc(DatabaseServer))
	dbParams := strings.Split(strings.TrimPrefix(db.URL, "http://"), ":")
	constants.DatabaseHost = dbParams[0]
	constants.DatabasePort, _ = strconv.Atoi(dbParams[1])

	mock := httptest.NewServer(http.HandlerFunc(SmtpMock))
	mockParams := strings.Split(strings.TrimPrefix(mock.URL, "http://"), ":")
	constants.MockServer = mockParams[0]
	constants.MockPort, _ = strconv.Atoi(mockParams[1])

	for idx, cs := range cases {
		var resp []byte
		if len(cs.UserJson) > 0 {
			r, err := http.Post(url(cs.Handler), "application/json", bytes.NewBuffer([]byte(cs.UserJson)))
			if err != nil {
				log.Printf("Unexpected error: %s", err.Error())
			}
			resp, err = ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("Fail to read body")
			}
		} else if cs.Handler == "verify" || cs.Handler == "validate" {
			token := logic.ReturnedAccessToken
			if cs.Handler == "verify" {
				token = logic.ReturnedConfirmationToken
			}
			r, err := http.Get(url(cs.Handler + token))
			if err != nil {
				log.Printf("Unexpected error: %s", err.Error())
			}
			resp, err = ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("Fail to read body")
			}
		} else if cs.Handler == "refresh" {
			jsonData, _ := json.Marshal(struct{
				refresh_token string
			}{logic.ReturnedRefreshToken})
			r, err := http.Post(url(cs.Handler), "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("Unexpected error: %s", err.Error())
			}
			resp, err = ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("Fail to read body")
			}
		}
		if cs.GetMessage {
			if string(resp) != cs.Message {
				log.Printf("FAIL: Test %d: expected %s, got %s", idx, cs.Message, string(resp))
			}
		}
		if cs.GetTokens {
			var tokens struct {
				access_token string
				refresh_token string
			}
			err := json.Unmarshal(resp, tokens)
			if err != nil {
				log.Printf("Unexpected error: %s", err.Error())
			}
			if tokens.access_token != logic.ReturnedAccessToken {
				log.Printf("FAIL: Test %d: expected %s, got %s", idx, logic.ReturnedAccessToken, tokens.access_token)
			}
			if tokens.refresh_token != logic.ReturnedRefreshToken {
				log.Printf("FAIL: Test %d: expected %s, got %s", idx, logic.ReturnedAccessToken, tokens.access_token)
			}
		}
	}
}
