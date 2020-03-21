package handlers

import (
    "../logic"
    "../models"
    "encoding/json"
    "html/template"
    "io/ioutil"
    "net/http"
    "strconv"
)

func GetProductsList(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    response, requestError := logic.Get("get_products")
    if requestError != nil {
        // ToDO: log
        http.Error(w, requestError.ErrorString, requestError.ErrorCode)
        return
    }
    defer response.Body.Close()
    body, bodyParseError := ioutil.ReadAll(response.Body)
    if bodyParseError != nil {
        //fmt.Println(bodyParseError.Error())
        http.Error(w, bodyParseError.Error(), http.StatusInternalServerError)
        return
    }
    var parsedData []interface {}
    parseError := json.Unmarshal(body, &parsedData)
    if parseError != nil {
        //fmt.Println(parseError.Error())
        http.Error(w, parseError.Error(), http.StatusInternalServerError)
        return
    }
    var array = make([]models.Product, len(parsedData))
    for index, item := range parsedData {
        switch item.(type) {
        case []interface{}:
            elem := item.([]interface{})
            if len(elem) == 3 {
                array[index] = models.Product{}
                switch elem[0].(type) {
                case float64:
                    array[index].Id = int(elem[0].(float64))
                default:
                    w.WriteHeader(http.StatusInternalServerError)
                    return
                }
                switch elem[1].(type) {
                case string:
                    array[index].Name = elem[1].(string)
                default:
                    w.WriteHeader(http.StatusInternalServerError)
                    return
                }
                switch elem[2].(type) {
                case string:
                    array[index].Category = elem[2].(string)
                default:
                    w.WriteHeader(http.StatusInternalServerError)
                    return
                }
            } else {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
        default:
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }
    jsonData, jsonError := json.Marshal(array)
    if jsonError != nil {
       //fmt.Println(parseError.Error())
       http.Error(w, jsonError.Error(), http.StatusInternalServerError)
       return
    }
    _, writeError := w.Write(jsonData)
    if writeError != nil {
      //fmt.Println(writeError.Error())
      http.Error(w, writeError.Error(), http.StatusInternalServerError)
      return
    }
    //t, err := template.ParseFiles("handlers/main.html")
    //if err != nil {
    //    http.Error(w, err.Error(), 500)
    //    return
    //}
    //if err := t.Execute(w, nil); err != nil {
    //    http.Error(w, err.Error(), 500)
    //    return
    //}
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        http.ServeFile(w, r, "frontend/addProduct.html")
    } else if r.Method == http.MethodPost {
        w.Header().Set("Content-Type", "application/json")
        body := r.Body
        defer body.Close()
        readBody, bodyParseError := ioutil.ReadAll(body)
        if bodyParseError != nil {
            http.Error(w, bodyParseError.Error(), http.StatusBadRequest)
            return
        }
        response, requestError := logic.Post("add_product", &readBody)
        if requestError != nil {
            // ToDO: log
            http.Error(w, requestError.ErrorString, requestError.ErrorCode)
            return
        }
        writeData, jsonError := logic.GetProductJSON(response.Body, nil)
        if jsonError != nil {
            // ToDO: log
            http.Error(w, jsonError.ErrorString, jsonError.ErrorCode)
            return
        }
        _, writeError := w.Write(*writeData)
        if writeError != nil {
            //fmt.Println(writeError.Error())
            http.Error(w, writeError.Error(), http.StatusInternalServerError)
            return
        }
    }
}

func ProductCard(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    response, requestError := logic.Get(r.URL.Query().Get("id"))
    if requestError != nil {
        // ToDO: log
        http.Error(w, requestError.ErrorString, requestError.ErrorCode)
        return
    }
    writeData, jsonError := logic.GetProductJSON(response.Body, nil)
    if jsonError != nil {
        // ToDO: log
        http.Error(w, jsonError.ErrorString, jsonError.ErrorCode)
        return
    }
    _, writeError := w.Write(*writeData)
    if writeError != nil {
        //fmt.Println(writeError.Error())
        http.Error(w, writeError.Error(), http.StatusInternalServerError)
        return
    }
    //data := map[string]interface{}{"name": string(body), "id": r.URL.Query().Get("id")}
    //t, err := template.ParseFiles("handlers/bookCardPage.html")
    //if err != nil {
    //    http.Error(w, err.Error(), 500)
    //    return
    //}
    //if err := t.Execute(w, data); err != nil {
    //    http.Error(w, err.Error(), 500)
    //    return
    //}
}

func EditProduct(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        data := map[string]interface{}{"id": r.URL.Query().Get("id")}
        // ToDo: edit get case
        t, err := template.ParseFiles("handlers/editBookPage.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        if err := t.Execute(w, data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    } else if r.Method == http.MethodPost {
        w.Header().Set("Content-Type", "application/json")
        var data models.Product
        intId, intError := strconv.Atoi(r.URL.Query().Get("id"))
        if intError != nil {
            //fmt.Println(intError.Error())
            http.Error(w, intError.Error(), http.StatusBadRequest)
            return
        }
        data.Id = intId
        jsonData, getJsonError := logic.GetProductJSON(r.Body, &data)
        if getJsonError != nil {
            // ToDO: log
            http.Error(w, getJsonError.ErrorString, getJsonError.ErrorCode)
            return
        }
        response, requestError := logic.Post("edit_product", jsonData)
        if requestError != nil {
            // ToDO: log
            http.Error(w, requestError.ErrorString, requestError.ErrorCode)
            return
        }
        writeData, jsonError := logic.GetProductJSON(response.Body, nil)
        if jsonError != nil {
            // ToDO: log
            http.Error(w, jsonError.ErrorString, jsonError.ErrorCode)
            return
        }
        _, writeError := w.Write(*writeData)
        if writeError != nil {
            //fmt.Println(writeError.Error())
            http.Error(w, writeError.Error(), http.StatusInternalServerError)
            return
        }
        //defer resp.Body.Close()
        //body, bodyParseError := ioutil.ReadAll(resp.Body)
        //if bodyParseError != nil {
        //    //fmt.Println(bodyParseError.Error())
        //    http.Error(w, bodyParseError.Error(), http.StatusInternalServerError)
        //    return
        //}
        //fmt.Println(w, body)
        //fmt.Fprint(w, "Название книги успешно изменено.\n")
    }
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    id, intError := strconv.Atoi(r.URL.Query().Get("id"))
    if intError != nil {
        //fmt.Println(intError)
        http.Error(w, intError.Error(), http.StatusBadRequest)
        return
    }
    jsonData, parseError := json.Marshal(struct{Id int `json:"id"`}{id})
    if parseError != nil {
        //fmt.Println(parseError.Error())
        http.Error(w, parseError.Error(), http.StatusInternalServerError)
        return
    }
    response, requestError := logic.Post("delete_product", &jsonData)
    if requestError != nil {
        // ToDO: log
        http.Error(w, requestError.ErrorString, requestError.ErrorCode)
        return
    }
    writeData, jsonError := logic.GetProductJSON(response.Body, nil)
    if jsonError != nil {
        // ToDO: log
        http.Error(w, jsonError.ErrorString, jsonError.ErrorCode)
        return
    }
    _, writeError := w.Write(*writeData)
    if writeError != nil {
        //fmt.Println(writeError.Error())
        http.Error(w, writeError.Error(), http.StatusInternalServerError)
        return
    }
}
