package handlers

import (
    "../../common/constants"
    "../../common/logic"
    "../../common/models"
    "../../utils"
    "../logic"
    "../models"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "sort"
    "strconv"
    "sync"
)

var DatabaseConnector = commonModels.Connector{
    Router: commonModels.Router{Host: commonLogic.GetUrl(constants.Protocol, constants.DatabaseHost, constants.DatabasePort)},
    Mutex: sync.Mutex{},
}

func GetProductsList(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    replier := commonModels.Replier{Writer: &w}
    checker := commonModels.ErrorChecker{Replier: &replier}

    response, requestError := DatabaseConnector.Get("")
    if checker.CheckError(requestError) {
        return
    }

    defer response.Body.Close()
    body, bodyParseError := ioutil.ReadAll(response.Body)
    if checker.CheckCustomError(bodyParseError, http.StatusInternalServerError) {
        return
    }

    var parsedData []interface {}
    parseError := json.Unmarshal(body, &parsedData)
    if checker.CheckCustomError(parseError, http.StatusInternalServerError) {
        return
    }

    castError := func() bool { return checker.NewError("500 Internal server error", http.StatusInternalServerError) }
    array := make([]models.Product, len(parsedData))
    for index, item := range parsedData {
        switch item.(type) {
        case []interface{}:
            elem := item.([]interface{})
            if len(elem) > 2 {
                array[index] = models.Product{}

                if id, ok := elem[0].(float64); ok {
                    array[index].Id = int(id)
                } else {
                    castError()
                    return
                }

                if name, ok := elem[1].(string); ok {
                    array[index].Name = name
                } else {
                    castError()
                    return
                }

                if category, ok := elem[2].(string); ok {
                    array[index].Category = category
                } else {
                    castError()
                    return
                }
            } else {
                castError()
                return
            }
        default:
            castError()
            return
        }
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

    res := models.AllItems{}

    countStr := r.URL.Query().Get("count")
    pageStr := r.URL.Query().Get("page")
    if len(countStr) > 0 && len(pageStr) > 0 {
        count, offsetError := strconv.Atoi(countStr)
        if checker.CheckCustomError(offsetError, http.StatusBadRequest) {
            return
        }
        if count < 1 {
            checker.NewError("Invalid param count", http.StatusBadRequest)
        }

        page, pageError := strconv.Atoi(pageStr)
        if checker.CheckCustomError(pageError, http.StatusBadRequest) {
            return
        }
        if count < 1 {
            checker.NewError("Invalid param page", http.StatusBadRequest)
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

func AddProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    replier := commonModels.Replier{Writer: &w}
    checker := commonModels.ErrorChecker{Replier: &replier}

    authError := logic.IsAuthorised(r.Header.Get("AccessToken"))
    if authError != nil {
        checker.NewError("Invalid access token", http.StatusUnauthorized)
        return
    }

    body := r.Body
    defer body.Close()
    readBody, bodyParseError := ioutil.ReadAll(body)
    if checker.CheckCustomError(bodyParseError, http.StatusBadRequest) {
        return
    }

    response, requestError := DatabaseConnector.Post("product", &readBody)
    if checker.CheckError(requestError) {
        return
    }

    writeData, jsonError := logic.GetProductJSON(response.Body)
    if checker.CheckError(jsonError) {
        return
    }

    if checker.CheckCustomError(replier.ReplyWithData(*writeData), http.StatusInternalServerError) {
        return
    }
}

func ProductCard(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    replier := commonModels.Replier{Writer: &w}
    checker := commonModels.ErrorChecker{Replier: &replier}

    response, requestError := DatabaseConnector.Get("product/" + r.URL.Query().Get("id"))
    if checker.CheckError(requestError) {
        return
    }

    writeData, jsonError := logic.GetProductJSON(response.Body)
    if checker.CheckError(jsonError) {
        return
    }

    if checker.CheckCustomError(replier.ReplyWithData(*writeData), http.StatusInternalServerError) {
        return
    }
}

func EditProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    replier := commonModels.Replier{Writer: &w}
    checker := commonModels.ErrorChecker{Replier: &replier}

    authError := logic.IsAuthorised(r.Header.Get("AccessToken"))
    if authError != nil {
        checker.NewError("Invalid access token", http.StatusUnauthorized)
        return
    }

    jsonData, getJsonError := logic.GetProductJSON(r.Body)
    if checker.CheckError(getJsonError) {
        return
    }

    response, requestError := DatabaseConnector.Put("product", jsonData)
    if checker.CheckError(requestError) {
        return
    }

    writeData, jsonError := logic.GetProductJSON(response.Body)
    if checker.CheckError(jsonError) {
        return
    }

    if checker.CheckCustomError(replier.ReplyWithData(*writeData), http.StatusInternalServerError) {
        return
    }
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    replier := commonModels.Replier{Writer: &w}
    checker := commonModels.ErrorChecker{Replier: &replier}

    authError := logic.IsAuthorised(r.Header.Get("AccessToken"))
    if authError != nil {
        checker.NewError("Invalid access token", http.StatusUnauthorized)
        return
    }

    jsonData, jsonError := logic.GetProductJSON(r.Body)
    if checker.CheckError(jsonError) {
        return
    }

    response, requestError := DatabaseConnector.Delete("product", jsonData)
    if checker.CheckError(requestError) {
        return
    }

    writeData, getJsonError := logic.GetProductJSON(response.Body)
    if checker.CheckError(getJsonError) {
        return
    }

    if checker.CheckCustomError(replier.ReplyWithData(*writeData), http.StatusInternalServerError) {
        return
    }
}
