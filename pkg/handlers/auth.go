package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"shop/pkg/auth"
	"shop/pkg/constants"
	"shop/pkg/logic"
	"shop/pkg/models"
	"shop/pkg/sessions"
	"shop/pkg/user"
	"strconv"
	"time"
)

type AuthHandler struct {
	Repo 			*user.Repo
	Sessions 		*sessions.Repo
	Notifications	*NotificationHandler
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	// parse user from body
	usr := &user.User{}
	if checker.CheckError(usr.GetUser(r.Body)) {
		return
	}
	log.Printf("Trying to register email %s\n", usr.Email)

	// add new user to db
	_, requestError := h.Repo.AddUser(usr)
	if requestError != nil {
		switch requestError.ErrorCode {
		case http.StatusConflict:
			checker.CheckCustomError(replier.ReplyWithMessage(constants.UserAlreadyExists), http.StatusInternalServerError)
			return
		default:
			checker.CheckError(requestError)
		}
	}

	// generate token to verify account
	confirmationToken, confirmationTokenError := logic.CreateConfirmationToken()
	if checker.CheckCustomError(confirmationTokenError, http.StatusInternalServerError) {
		return
	}

	// send confirmation request to notification queue
	confirmationError := h.Notifications.SendConfirmationRequest(usr, *confirmationToken)
	if checker.CheckCustomError(confirmationError, http.StatusInternalServerError) {
		log.Println("Failed to send confirmation request into queue")
		return
	}

	// make confirmation
	confirmation := auth.Confirmation{}
	confirmation.Init(usr.Email, *confirmationToken)
	log.Printf("Make confirmation request %s for user %s\n", *confirmationToken, usr.Email)

	// add confirmation to db
	_, addConfirmationError := h.Repo.AddConfirmation(&confirmation)
	if checker.CheckError(addConfirmationError) {
		return
	}

	// delete confirmation from db when token expired
	go func() {
		time.Sleep(constants.ConfirmationTokenExpireTime)
		resp, _ := h.Repo.DeleteConfirmation(&confirmation)
		if resp.StatusCode == http.StatusOK {
			log.Printf("Delete confirmation %s as it has been expired\n", confirmation.Token)
		}
	}()

	log.Printf("User %s has been signed up\n", usr.Email)
	if checker.CheckCustomError(
		replier.ReplyWithMessage(constants.SignUpOkMessage(usr.Email)),
		http.StatusInternalServerError,
	) {
		log.Println("Failed to reply")
		return
	}
}

func (h *AuthHandler) ConfirmRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	// get confirmation token from url
	token := mux.Vars(r)["token"]
	log.Printf("Trying to confirm user by token %s\n", token)

	// get confirmation from db by token
	confirmationData, confirmationError := h.Repo.GetConfirmation(token)
	if checker.CheckError(confirmationError) {
		return
	}

	// parse confirmation from db response
	var confirmation auth.Confirmation
	jsonError := confirmation.GetConfirmation(confirmationData.Body)
	if checker.CheckError(jsonError) {
		log.Printf("Bad response from database for confirmation token %s\n", token)
		return
	}

	// check if token has not been expired
	if expire, expireErr := strconv.Atoi(confirmation.Expire); expireErr != nil {
		if expire < int(time.Now().Unix()) {
			log.Printf("Token %s has been expired\n", token)
			checker.NewError(constants.ExpiredConfirmation, http.StatusBadRequest)
		}
	} else {
		log.Println("Failed to convert time into integer")
		checker.CheckCustomError(expireErr, http.StatusInternalServerError)
		return
	}

	// confirm user in db
	usr := user.User{Email: confirmation.Email}
	_, confirmError := h.Repo.ConfirmUser(&usr)
	if checker.CheckError(confirmError) {
		return
	}

	// delete confirmation in db as account has been verified
	go func() {
		h.Repo.DeleteConfirmation(&confirmation)
		log.Printf("Confirmation with token %s has been deleted\n", token)
	}()

	log.Printf("User %s has been verified\n", confirmation.Email)
	if checker.CheckCustomError(
		replier.ReplyWithMessage(constants.ConfirmOkMessage),
		http.StatusInternalServerError,
	) {
		return
	}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	// parse user from body
	usr := &user.User{}
	if checker.CheckError(usr.GetUser(r.Body)) {
		return
	}
	log.Printf("Trying to sign in for user %s\n", usr.Email)

	// get user by email from db to check confirmation
	jsonUserData, jsonUserDataError := h.Repo.GetUser(usr)
	if jsonUserData.StatusCode == http.StatusNotFound {
		log.Printf("User %s doesn't exist\n", usr.Email)
		checker.NewError(constants.InvalidUser, http.StatusBadRequest)
		return
	}
	if checker.CheckError(jsonUserDataError) {
		return
	}

	// parse db response
	if checker.CheckError(usr.GetUser(jsonUserData.Body)) {

		return
	}
	// check if password is valid and user is verified
	if usr.Password != usr.Password {
		log.Printf("User %s send invalid password\n", usr.Email)
		checker.NewError(constants.InvalidUser, http.StatusBadRequest)
		return
	}
	if !usr.Confirm {
		log.Printf("User %s is not confirmed\n", usr.Email)
		checker.NewError(constants.NotConfirmed, http.StatusBadRequest)
		return
	}

	// generate refresh token
	refreshToken, refreshTokenError := logic.CreateRefreshToken()
	if checker.CheckCustomError(refreshTokenError, http.StatusInternalServerError) {
		return
	}

	// create session
	session := sessions.Session{}
	session.Init(usr.Email, *refreshToken)
	log.Printf("Creating session for user %s\n", usr.Email)

	// add session to db
	_, requestError := h.Sessions.AddOrUpdate(&session)
	if checker.CheckError(requestError) {
		return
	}
	// delete session when it expired
	go func() {
		time.Sleep(constants.RefreshTokenExpireTime)
		resp, _ := h.Sessions.Delete(&session)
		if resp.StatusCode == http.StatusOK {
			log.Printf("Session for user %s has been deleted as expired\n", usr.Email)
		}
	}()

	// generate access token
	accessToken, tokenError := logic.CreateAccessToken(usr.Email)
	if checker.CheckCustomError(tokenError, http.StatusInternalServerError) {
		return
	}

	log.Printf("Sign in for user %s\n", usr.Email)
	checker.CheckCustomError(
		replier.ReplyWithJson(
			models.AuthToken{
				AccessToken: accessToken,
				RefreshToken: *refreshToken,
			},
		),
		http.StatusInternalServerError,
	)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	body := r.Body
	defer body.Close()
	readBody, bodyParseError := ioutil.ReadAll(body)
	if checker.CheckCustomError(bodyParseError, http.StatusBadRequest) {
		log.Println("Error while reading body")
		return
	}

	// parse access token from body
	var parsedBody models.AccessToken
	parsedBodyError := json.Unmarshal(readBody, &parsedBody)
	if checker.CheckCustomError(parsedBodyError, http.StatusBadRequest) {
		log.Println("Error while parsing body")
		return
	}

	// check if access token exists
	accessTokenStr := parsedBody.TokenStr
	if len(accessTokenStr) == 0 {
		log.Println("Authorization token is empty string")
		checker.NewError(constants.Unauthorized, http.StatusBadRequest)
		return
	}

	// parse access token into struct
	accessToken, accessTokenParseError := jwt.Parse(
		accessTokenStr,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(constants.SigningToken), nil
		},
	)
	if checker.CheckCustomError(accessTokenParseError, http.StatusInternalServerError) {
		log.Println("Failed to parse access token")
		return
	}

	// check if access token has not been expired
	if accessToken.Valid {
		log.Printf("Token %s is valid \n", accessTokenStr)
		checker.CheckCustomError(replier.ReplyWithMessage(constants.ValidAccessToken), http.StatusInternalServerError)
	} else {
		log.Printf("Token %s is invalid \n", accessTokenStr)
		checker.NewError(constants.Unauthorized, http.StatusForbidden)
	}
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	// parse session from body (with refresh token only)
	session := sessions.Session{}
	session.GetSession(r.Body)
	log.Printf("Trying to refresh session with token %s\n", session.RefreshToken)

	// generate new refresh token
	refreshToken, refreshTokenError := logic.CreateRefreshToken()
	if checker.CheckCustomError(refreshTokenError, http.StatusInternalServerError) {
		return
	}

	session.Update(*refreshToken)

	// update session in db
	response, requestError := h.Sessions.AddOrUpdate(&session)
	if checker.CheckError(requestError) {
		return
	}
	if response.StatusCode != http.StatusOK {
		switch response.StatusCode {
		case http.StatusNotFound:
			log.Printf("No session with refresh token %s\n", session.RefreshToken)
			checker.NewError(constants.InvalidRefreshToken, http.StatusNotFound)
			return
		case http.StatusForbidden:
			log.Printf("Refresh token %s has been expired\n", session.RefreshToken)
			checker.NewError(constants.InvalidRefreshToken, http.StatusForbidden)
			return
		}
	}

	// generate new access token
	accessToken, signedError := logic.CreateAccessToken(session.Email)
	if checker.CheckCustomError(signedError, http.StatusInternalServerError) {
		return
	}

	log.Println("New tokens have been generated")
	checker.CheckCustomError(
		replier.ReplyWithJson(
			models.AuthToken{
				AccessToken: accessToken,
				RefreshToken: *refreshToken,
			},
		),
		http.StatusInternalServerError,
	)
}
