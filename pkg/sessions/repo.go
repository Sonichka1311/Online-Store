package sessions

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"shop/pkg/constants"
	"shop/pkg/database"
	"shop/pkg/models"
)

type Repo struct {
	Connector *database.Connector
}

var (
	sessionsTable = "sessions"
)

func (r *Repo) Add(session *Session) *models.Error {
	log.Printf("Add session with token %s", session.RefreshToken)
	row := r.Connector.SelectOne("login", sessionsTable, "login = ?", session.Email)
	dbErr := row.Scan(&session.Email)
	switch dbErr {
	case sql.ErrNoRows:
		_, dbErr = r.Connector.Insert(sessionsTable, "login, token, expire", "?, ?, ?", session.Email, session.RefreshToken, session.Expire)
		if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
			return err
		}
	case nil:
		_, dbErr = r.Connector.Update(sessionsTable, "token = ?, expire = ?", "login = ?", session.RefreshToken, session.Expire, session.Email)
		if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
			return err
		}
	default:
		if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
			return err
		}
	}
	return nil
}

func (r *Repo) Get(session *Session) *models.Error {
	row := r.Connector.SelectOne("login, expire", sessionsTable, "token = ?", session.RefreshToken)
	dbErr := row.Scan(&session.Email, &session.Expire)
	if dbErr == sql.ErrNoRows {
		if err, isErr := models.NewError(errors.New(constants.TokenIsExpired), http.StatusBadRequest); isErr {
			return err
		}
	}
	if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
		return err
	}
	return nil
}

func (r *Repo) Delete(session *Session) bool {
	_, dbErr := r.Connector.Delete(sessionsTable, "login = ?", session.Email)
	return dbErr == nil
}
