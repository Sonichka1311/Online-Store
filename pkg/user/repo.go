package user

import (
	"database/sql"
	"errors"
	"net/http"
	"shop/pkg/auth"
	"shop/pkg/constants"
	"shop/pkg/database"
	"shop/pkg/models"
	"time"
)

type Repo struct {
	Connector *database.Connector
}

var (
	usersTable = "users"
	confirmationsTable = "confirmations"
)

func (r *Repo) AddUser(usr *User) *models.Error {
	row := r.Connector.SelectOne("login", usersTable, "login = ?", usr.Email)
	var login string
	dbErr := row.Scan(&login)
	if dbErr != sql.ErrNoRows {
		if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
			return err
		}
		if err, isErr := models.NewError(errors.New(constants.UserAlreadyExists), http.StatusConflict); isErr {
			return err
		}
	}
	_, dbErr = r.Connector.Insert(usersTable, "login, password, confirm", "?, ?, ?", usr.Email, usr.Password, false)
	if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
		return err
	}
	return nil
}

func (r *Repo) GetUser(usr *User) *models.Error {
	row := r.Connector.SelectOne("password, confirm, role", usersTable, "login = ?", usr.Email)
	dbErr := row.Scan(&usr.Password, &usr.Confirm, &usr.Role)
	if dbErr == sql.ErrNoRows {
		if err, isErr := models.NewError(errors.New(constants.NoUser), http.StatusBadRequest); isErr {
			return err
		}
	}
	if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
		return err
	}
	return nil
}

func (r *Repo) AddConfirmation(confirmation *auth.Confirmation) *models.Error {
	row := r.Connector.SelectOne("login", confirmationsTable, "login = ?", confirmation.Email)
	var log string
	dbErr := row.Scan(&log)
	if dbErr != sql.ErrNoRows {
		if err, isErr := models.NewError(errors.New(constants.UserAlreadyExists), http.StatusConflict); isErr {
			return err
		}
	}

	_, dbErr = r.Connector.Insert(
		confirmationsTable,
		"login, token, expire", "?, ?, ?",
		confirmation.Email, confirmation.Token, confirmation.Expire,
	)
	if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
		return err
	}

	go func() {
		time.Sleep(2 * constants.ConfirmationTokenExpireTime)
		r.Connector.Delete(confirmationsTable, "token = ?", confirmation.Token)
	}()

	return nil
}

func (r *Repo) GetConfirmation(confirmation *auth.Confirmation) *models.Error {
	row := r.Connector.SelectOne("login, expire", confirmationsTable, "token = ?", confirmation.Token)
	dbErr := row.Scan(&confirmation.Email, &confirmation.Expire)
	if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
		return err
	}
	return nil
}

func (r *Repo) DeleteConfirmation(confirmation *auth.Confirmation) bool {
	_, err := r.Connector.Delete(confirmationsTable, "token = ?", confirmation.Token)
	return err == nil
}

func (r *Repo) ConfirmUser(confirmation *auth.Confirmation) *models.Error {
	_, dbErr := r.Connector.Update(usersTable, "confirm = ?", "login = ?", true, confirmation.Email)
	if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
		return err
	}
	return nil
}

func (r *Repo) UpgradeUser(usr *User) *models.Error {
	_, dbErr := r.Connector.Update(usersTable, "role = ?", "login = ?", "admin", usr.Email)
	if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
		return err
	}
	return nil
}

func (r *Repo) DowngradeUser(usr *User) *models.Error {
	_, dbErr := r.Connector.Update(usersTable, "role = ?", "login = ?", "user", usr.Email)
	if err, isErr := models.NewError(dbErr, http.StatusInternalServerError); isErr {
		return err
	}
	return nil
}
