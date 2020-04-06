package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io/ioutil"
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
	confirmationError := h.Notifications.SendRequest(usr, *confirmationToken)
	if checker.CheckCustomError(confirmationError, http.StatusInternalServerError) {
		return
	}

	// make confirmation
	confirmation := auth.Confirmation{}
	confirmation.Init(usr.Email, *confirmationToken)

	// add confirmation to db
	_, addConfirmationError := h.Repo.AddConfirmation(&confirmation)
	if checker.CheckError(addConfirmationError) {
		return
	}

	// delete confirmation from db when token expired
	go func() {
		time.Sleep(constants.ConfirmationTokenExpireTime)
		h.Repo.DeleteConfirmation(&confirmation)
	}()

	if checker.CheckCustomError(
		replier.ReplyWithMessage(constants.SignUpOkMessage(usr.Email)),
		http.StatusInternalServerError,
	) {
		return
	}
}

func (h *AuthHandler) ConfirmRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	// get confirmation token from url
	token := mux.Vars(r)["token"]

	// get confirmation from db by token
	confirmationData, confirmationError := h.Repo.GetConfirmation(token)
	if checker.CheckError(confirmationError) {
		return
	}

	// parse confirmation from db response
	var confirmation auth.Confirmation
	jsonError := confirmation.GetConfirmation(confirmationData.Body)
	if checker.CheckError(jsonError) {
		return
	}

	// check if token has not been expired
	if expire, expireErr := strconv.Atoi(confirmation.Expire); expireErr != nil {
		if expire < int(time.Now().Unix()) {
			checker.NewError(constants.ExpiredConfirmation, http.StatusBadRequest)
		}
	} else {
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
	}()

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

	// get user by email from db to check confirmation
	jsonUserData, jsonUserDataError := h.Repo.GetUser(usr)
	if jsonUserData.StatusCode == http.StatusNotFound {
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
		checker.NewError(constants.InvalidUser, http.StatusBadRequest)
		return
	}
	if !usr.Confirm {
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

	// add session to db
	_, requestError := h.Sessions.AddOrUpdate(&session)
	if checker.CheckError(requestError) {
		return
	}
	// delete session when it expired
	go func() {
		time.Sleep(constants.RefreshTokenExpireTime)
		h.Sessions.Delete(&session)
	}()

	// generate access token
	accessToken, tokenError := logic.CreateAccessToken(usr.Email)
	if checker.CheckCustomError(tokenError, http.StatusInternalServerError) {
		return
	}

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
		return
	}

	// parse access token from body
	var parsedBody models.AccessToken
	parsedBodyError := json.Unmarshal(readBody, &parsedBody)
	if checker.CheckCustomError(parsedBodyError, http.StatusBadRequest) {
		return
	}

	// check if access token exists
	accessTokenStr := parsedBody.TokenStr
	if len(accessTokenStr) == 0 {
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
		return
	}

	// check if access token has not been expired
	if accessToken.Valid {
		checker.CheckCustomError(replier.ReplyWithMessage(constants.ValidAccessToken), http.StatusInternalServerError)
	} else {
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
			checker.NewError(constants.InvalidRefreshToken, http.StatusNotFound)
			return
		case http.StatusForbidden:
			checker.NewError(constants.InvalidRefreshToken, http.StatusForbidden)
			return
		}
	}

	// generate new access token
	accessToken, signedError := logic.CreateAccessToken(session.Email)
	if checker.CheckCustomError(signedError, http.StatusInternalServerError) {
		return
	}

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
