package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"shop/pkg/constants"
	"shop/pkg/logic"
	"shop/pkg/models"
	"shop/pkg/user"
	"strconv"
	"time"
)

type AuthHandler struct {
	Repo *user.Repo
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	body := r.Body
	defer body.Close()
	readBody, bodyParseError := ioutil.ReadAll(body)
	if checker.CheckCustomError(bodyParseError, http.StatusBadRequest) {
		return
	}

	var parsedBody *models.User
	parsedBodyError := json.Unmarshal(readBody, &parsedBody)
	if checker.CheckCustomError(parsedBodyError, http.StatusBadRequest) {
		return
	}

	_, requestError := h.Repo.Add(&readBody)
	if requestError != nil {
		switch requestError.ErrorCode {
		case http.StatusConflict:
			checker.CheckCustomError(replier.ReplyWithMessage("User with this email already exists"), http.StatusInternalServerError)
			return
		default:
			checker.NewError(requestError.ErrorString, requestError.ErrorCode)
		}
	}

	if checker.CheckCustomError(replier.ReplyWithMessage(fmt.Sprintf("Welcome, %s!\n", parsedBody.Email)), http.StatusInternalServerError) {
		return
	}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	body := r.Body
	defer body.Close()
	readBody, bodyParseError := ioutil.ReadAll(body)
	if checker.CheckCustomError(bodyParseError, http.StatusBadRequest) {
		return
	}

	var parsedBody *models.User
	parsedBodyError := json.Unmarshal(readBody, &parsedBody)
	if checker.CheckCustomError(parsedBodyError, http.StatusBadRequest) {
		return
	}

	refreshToken := logic.CreateRefreshToken()

	session := models.Session{
		Email: 			parsedBody.Email,
		Password: 		parsedBody.Password,
		RefreshToken: 	refreshToken,
		Expire: 		strconv.Itoa(int(time.Now().Add(constants.RefreshTokenExpireTime).Unix())),
	}

	jsonData, jsonError := json.Marshal(session)
	if checker.CheckCustomError(jsonError, http.StatusInternalServerError) {
		return
	}

	_, requestError := h.Repo.Get(&jsonData)
	if requestError != nil {
		switch requestError.ErrorCode {
		case http.StatusBadRequest:
			checker.NewError("Incorrect password for user " + parsedBody.Email, http.StatusBadRequest)
			return
		case http.StatusNotFound:
			checker.NewError("There is no user " + parsedBody.Email, http.StatusBadRequest)
			return
		default:
			checker.NewError(requestError.ErrorString, requestError.ErrorCode)
			return
		}
	}

	accessToken, tokenError := logic.CreateAccessToken(parsedBody.Email)
	if checker.CheckCustomError(tokenError, http.StatusInternalServerError) {
		return
	}

	checker.CheckCustomError(
		replier.ReplyWithJson(models.AuthToken{AccessToken: accessToken, RefreshToken: refreshToken}),
		http.StatusInternalServerError)
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

	var parsedBody models.AccessToken
	parsedBodyError := json.Unmarshal(readBody, &parsedBody)
	if checker.CheckCustomError(parsedBodyError, http.StatusBadRequest) {
		return
	}

	accessTokenStr := parsedBody.TokenStr
	if accessTokenStr == "" {
		checker.NewError("No access token", http.StatusBadRequest)
		return
	}

	accessToken, accessTokenParseError := jwt.Parse(
		accessTokenStr,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(constants.SigningToken), nil
		},
	)
	if checker.CheckCustomError(accessTokenParseError, http.StatusInternalServerError) {
		return
	}

	if accessToken.Valid {
		checker.CheckCustomError(replier.ReplyWithMessage("Access token is valid"), http.StatusInternalServerError)
	} else {
		checker.NewError("Invalid token", http.StatusForbidden)
	}
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	body := r.Body
	defer body.Close()
	readBody, bodyParseError := ioutil.ReadAll(body)
	if checker.CheckCustomError(bodyParseError, http.StatusBadRequest) {
		return
	}

	var parsedBody *models.RefreshToken
	parsedBodyError := json.Unmarshal(readBody, &parsedBody)
	if checker.CheckCustomError(parsedBodyError, http.StatusBadRequest) {
		return
	}

	refreshToken := logic.CreateRefreshToken()

	data := models.RefreshTokenUpdate{
		OldToken: parsedBody.TokenStr,
		Token: refreshToken,
		Expire: strconv.Itoa(int(time.Now().Add(constants.RefreshTokenExpireTime).Unix())),
		Time: strconv.Itoa(int(time.Now().Unix())),
	}

	jsonData, jsonError := json.Marshal(data)
	if checker.CheckCustomError(jsonError, http.StatusInternalServerError) {
		return
	}

	response, requestError := h.Repo.Update(&jsonData)
	if checker.CheckError(requestError) {
		return
	}
	if response.StatusCode != http.StatusOK {
		switch response.StatusCode {
		case http.StatusNotFound:
			checker.NewError("Invalid refresh_token", http.StatusNotFound)
			return
		case http.StatusForbidden:
			checker.NewError("Expired refresh_token", http.StatusForbidden)
			return
		}
	}

	var parsedResp models.User
	parsedRespError := json.Unmarshal(readBody, parsedBody)
	if checker.CheckCustomError(parsedRespError, http.StatusInternalServerError) {
		return
	}

	accessToken, signedError := logic.CreateAccessToken(parsedResp.Email)
	if checker.CheckCustomError(signedError, http.StatusInternalServerError) {
		return
	}

	checker.CheckCustomError(
		replier.ReplyWithJson(models.AuthToken{AccessToken: accessToken, RefreshToken: refreshToken}),
		http.StatusInternalServerError)
}
