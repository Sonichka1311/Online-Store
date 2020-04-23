package product

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

var productsTable = "products"

func (r *Repo) Get(id int) (*Product, *models.Error) {
	row := r.Connector.SelectOne("name, category", productsTable, "id = ?", id)
	var product Product
	product.Id = id
	dbError := row.Scan(&product.Name, &product.Category)
	if dbError == sql.ErrNoRows {
		if err, isErr := models.NewError(errors.New(constants.NoProduct), http.StatusNotFound); isErr {
			return nil, err
		}
	}
	if err, isErr := models.NewError(dbError, http.StatusInternalServerError); isErr {
		return nil, err
	}
	return &product, nil
}

func (r *Repo) GetAll() ([]*Product, *models.Error) {
	rows, queryError := r.Connector.SelectAll("id, name, category", productsTable)
	if err, isErr := models.NewError(queryError, http.StatusInternalServerError); isErr {
		return nil, err
	}
	products := make([]*Product, 0)
	for rows.Next() {
		var product Product
		dbError := rows.Scan(&product.Id, &product.Name, &product.Category)
		if err, isErr := models.NewError(dbError, http.StatusInternalServerError); isErr {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func (r *Repo) Add(product *Product) *models.Error {
	row := r.Connector.SelectOne("id", productsTable, "name = ?", product.Name)
	dbError := row.Scan(&product.Id)
	if dbError != sql.ErrNoRows {
		if err, isErr := models.NewError(dbError, http.StatusInternalServerError); isErr {
			return err
		}
		if err, isErr := models.NewError(errors.New(constants.ProductAlreadyExists), http.StatusBadRequest); isErr {
			return err
		}
	}

	res, dbError := r.Connector.Insert(productsTable, "name, category", "?, ?", product.Name, product.Category)
	if err, isErr := models.NewError(dbError, http.StatusInternalServerError); isErr {
		return err
	}
	id, resErr := res.LastInsertId()
	if err, isErr := models.NewError(resErr, http.StatusInternalServerError); isErr {
		return err
	}
	product.Id = int(id)
	return nil
}

func (r *Repo) Edit(product *Product) *models.Error {
	row := r.Connector.SelectOne("id", productsTable, "id = ?", product.Id)
	dbError := row.Scan(&product.Id)
	if dbError == sql.ErrNoRows {
		if err, isErr := models.NewError(errors.New(constants.NoProduct), http.StatusBadRequest); isErr {
			return err
		}
	} else if err, isErr := models.NewError(dbError, http.StatusInternalServerError); isErr {
		return err
	}

	if len(product.Name) != 0 {
		_, dbError = r.Connector.Update(productsTable, "name = ?", "id = ?", product.Name, product.Id)
	}
	if len(product.Category) != 0 {
		_, dbError = r.Connector.Update(productsTable, "category = ?", "id = ?", product.Category, product.Id)
	}
	if err, isErr := models.NewError(dbError, http.StatusInternalServerError); isErr {
		return err
	}
	return nil
}

func (r *Repo) Delete(product *Product) *models.Error {
	log.Println(product.Id)
	row := r.Connector.SelectOne("name, category", productsTable, "id = ?", product.Id)
	dbError := row.Scan(&product.Name, &product.Category)
	if dbError == sql.ErrNoRows {
		if err, isErr := models.NewError(errors.New(constants.NoProduct), http.StatusBadRequest); isErr {
			return err
		}
	}
	_, dbError = r.Connector.Delete(productsTable, "id = ?", product.Id)
	if err, isErr := models.NewError(dbError, http.StatusInternalServerError); isErr {
		return err
	}
	return nil
}
