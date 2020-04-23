package handlers

import (
	"encoding/json"
	"net/http"
	"shop/pkg/auth"
	"shop/pkg/constants"
	"shop/pkg/models"
	"shop/pkg/product"
	"shop/pkg/utils"
	"sort"
	"strconv"
	"sync"
)

type ProductHandler struct {
	Repo 	*product.Repo
	Auth 	*auth.Repo
	Mutex 	*sync.RWMutex
}

func (h *ProductHandler) GetProductsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	array, getError := h.Repo.GetAll()
	if checker.CheckError(getError) {
		return
	}

	sort.Slice(array, func(i int, j int) bool {
		if array[i].Id != array[j].Id {
			return array[i].Id < array[j].Id
		} else if array[i].Category != array[j].Category {
			return array[i].Category < array[j].Category
		} else {
			return array[i].Name < array[j].Name
		}
	})

	res := product.AllItems{}

	countStr := r.URL.Query().Get("count")
	pageStr := r.URL.Query().Get("page")
	if len(countStr) > 0 && len(pageStr) > 0 {
		count, offsetError := strconv.Atoi(countStr)
		if checker.CheckCustomError(offsetError, http.StatusBadRequest) {
			return
		}
		if count < 1 {
			checker.NewError(constants.InvalidParams, http.StatusBadRequest)
		}

		page, pageError := strconv.Atoi(pageStr)
		if checker.CheckCustomError(pageError, http.StatusBadRequest) {
			return
		}
		if count < 1 {
			checker.NewError(constants.InvalidParams, http.StatusBadRequest)
		}

		begin := utils.Min((page - 1) * count, 0)
		end := utils.Min(page * count, len(array))
		res.Items = array[begin : end]
		res.PagesCount = len(array) / count + utils.Int(len(array) % count != 0)
		res.CurrentPage = utils.Min(res.PagesCount, page)
		if len(array) == 0 {
			res.PagesCount = 1
			res.CurrentPage = 1
		}
	} else {
		res.Items = array
		res.PagesCount = 1
		res.CurrentPage = 1
	}

	jsonData, jsonError := json.Marshal(res)
	if checker.CheckCustomError(jsonError, http.StatusInternalServerError) {
		return
	}

	if checker.CheckCustomError(replier.ReplyWithData(jsonData), http.StatusInternalServerError) {
		return
	}
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	authError := h.Auth.Verify(r.Header.Get("AccessToken"))
	if authError != nil {
		checker.NewError(constants.Unauthorized, http.StatusUnauthorized)
		return
	}

	var prod product.Product
	parseError := prod.GetFromBody(r.Body)
	if checker.CheckError(parseError) {
		return
	}

	requestError := h.Repo.Add(&prod)
	if checker.CheckError(requestError) {
		return
	}

	jsonData, jsonError := prod.GetJson()
	if checker.CheckError(jsonError) {
		return
	}

	if checker.CheckCustomError(replier.ReplyWithData(jsonData), http.StatusInternalServerError) {
		return
	}
}

func (h *ProductHandler) ProductCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	id, intError := strconv.Atoi(r.URL.Query().Get("id"))
	if checker.CheckCustomError(intError, http.StatusBadRequest) {
		return
	}

	prod, getError := h.Repo.Get(id)
	if checker.CheckError(getError) {
		return
	}

	jsonData, jsonError := prod.GetJson()
	if checker.CheckError(jsonError) {
		return
	}

	if checker.CheckCustomError(replier.ReplyWithData(jsonData), http.StatusInternalServerError) {
		return
	}
}

func (h *ProductHandler) EditProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	authError := h.Auth.Verify(r.Header.Get("AccessToken"))
	if authError != nil {
		checker.NewError(constants.Unauthorized, http.StatusUnauthorized)
		return
	}

	var prod product.Product
	parseError := prod.GetFromBody(r.Body)
	if checker.CheckError(parseError) {
		return
	}

	h.Mutex.RLock()
	requestError := h.Repo.Edit(&prod)
	if checker.CheckError(requestError) {
		return
	}

	resp, getError := h.Repo.Get(prod.Id)
	if checker.CheckError(getError) {
		return
	}
	h.Mutex.RUnlock()

	jsonData, jsonError := resp.GetJson()
	if checker.CheckError(jsonError) {
		return
	}

	if checker.CheckCustomError(replier.ReplyWithData(jsonData), http.StatusInternalServerError) {
		return
	}
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	replier := models.Replier{Writer: &w}
	checker := models.ErrorChecker{Replier: &replier}

	authError := h.Auth.Verify(r.Header.Get("AccessToken"))
	if authError != nil {
		checker.NewError(constants.Unauthorized, http.StatusUnauthorized)
		return
	}

	var prod product.Product
	parseError := prod.GetFromBody(r.Body)
	if checker.CheckError(parseError) {
		return
	}

	h.Mutex.RLock()
	requestError := h.Repo.Delete(&prod)
	if checker.CheckError(requestError) {
		return
	}

	h.Mutex.RUnlock()

	jsonData, jsonError := prod.GetJson()
	if checker.CheckError(jsonError) {
		return
	}

	if checker.CheckCustomError(replier.ReplyWithData(jsonData), http.StatusInternalServerError) {
		return
	}
}
