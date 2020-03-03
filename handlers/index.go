package handlers

import (
    "../models"
    "../repository"
    "encoding/json"
    "fmt"
    "html/template"
    "io/ioutil"
    "net/http"
    "strconv"
)

var connector = repository.Connector{Host: "http://localhost:8888"}

func GetProductsList(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    resp, requestError := connector.Get("/get_products")
    if requestError != nil {
        //fmt.Println(requestError.Error())
        http.Error(w, requestError.Error(), http.StatusInternalServerError)
        return
    }
    if resp.StatusCode != http.StatusOK {
        //fmt.Println(resp.Status)
        http.Error(w, resp.Status, resp.StatusCode)
        return
    }
    defer resp.Body.Close()
    body, bodyParseError := ioutil.ReadAll(resp.Body)
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
    jsonData, parseError := json.Marshal(array)
    if parseError != nil {
       //fmt.Println(parseError.Error())
       http.Error(w, parseError.Error(), http.StatusInternalServerError)
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
        jsonData, parseError := json.Marshal(
            models.Product{
                Name: r.FormValue("name"),
                Id: 0,
                Category: r.FormValue("category")})
        if parseError != nil {
            //fmt.Println(parseError.Error())
            http.Error(w, parseError.Error(), http.StatusBadRequest)
            return
        }
        //fmt.Printf("%v\n", string(jsonData))
        resp, requestError := connector.Post("/add_product", jsonData)
        if requestError != nil {
            //fmt.Println(requestError.Error())
            http.Error(w, requestError.Error(), http.StatusInternalServerError)
            return
        }
        if resp.StatusCode != http.StatusOK {
            //fmt.Println(resp.Status)
            http.Error(w, resp.Status, resp.StatusCode)
            return
        }
        defer resp.Body.Close()
        body, bodyParseError := ioutil.ReadAll(resp.Body)
        if bodyParseError != nil {
           //fmt.Println(bodyParseError.Error())
            http.Error(w, bodyParseError.Error(), http.StatusInternalServerError)
           return
        }
        productId, intError := strconv.Atoi(string(body))
        if intError != nil {
            //fmt.Println(intError.Error())
            http.Error(w, intError.Error(), http.StatusInternalServerError)
            return
        }
        writeError := json.NewEncoder(w).Encode(struct{Id int `json:"id"`}{productId})
        if writeError != nil {
           //fmt.Println(writeError.Error())
           http.Error(w, writeError.Error(), http.StatusInternalServerError)
           return
        }
    }
}

func ProductCard(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    resp, requestError := connector.Get("/" + r.URL.Query().Get("id"))
    if requestError != nil {
        //fmt.Println(err.Error())
        http.Error(w, requestError.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    body, bodyParseError := ioutil.ReadAll(resp.Body)
    if bodyParseError != nil {
        //fmt.Println(bodyParseError.Error())
        http.Error(w, bodyParseError.Error(), http.StatusInternalServerError)
        return
    }
    if resp.StatusCode != http.StatusOK {
        //fmt.Println(resp.Status)
        http.Error(w, resp.Status, resp.StatusCode)
        return
    }
    var data models.Product
    unmarshalError := json.Unmarshal(body, &data)
    if unmarshalError != nil {
        //fmt.Println(unmarshalError.Error())
        http.Error(w, unmarshalError.Error(), http.StatusInternalServerError)
        return
    }
    writeError := json.NewEncoder(w).Encode(
        models.Product{
            Name: data.Name,
            Id: data.Id,
            Category: data.Category})
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
            fmt.Println(intError.Error())
            http.Error(w, intError.Error(), http.StatusBadRequest)
        }
        data.Id = intId
        if len(r.FormValue("name")) > 0 {
            data.Name = r.FormValue("name")
        }
        if len(r.FormValue("category")) > 0 {
            data.Category = r.FormValue("category")
        }
        jsonData, parseError := json.Marshal(data)
        if parseError != nil {
            //fmt.Println(parseError.Error())
            http.Error(w, parseError.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Printf("%v\n", string(jsonData))
        resp, requestError := connector.Post("/edit_product", jsonData)
        if requestError != nil {
            //fmt.Println(requestError.Error())
            http.Error(w, requestError.Error(), http.StatusInternalServerError)
            return
        }
        if resp.StatusCode != http.StatusOK {
            //fmt.Println(resp.Status)
            http.Error(w, resp.Status, resp.StatusCode)
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
    id, intError := strconv.Atoi(r.URL.Query().Get("id"))
    if intError != nil {
        //fmt.Println(intError)
        http.Error(w, intError.Error(), http.StatusBadRequest)
        return
    }
    jsonData, parseError := json.Marshal(struct{Id int `json:"id"`}{id})
    if parseError != nil {
        fmt.Println(parseError.Error())
        http.Error(w, parseError.Error(), http.StatusInternalServerError)
        return
    }
    resp, requestError := connector.Post("/delete_product", jsonData)
    if requestError != nil {
        fmt.Println(requestError.Error())
        http.Error(w, requestError.Error(), http.StatusInternalServerError)
        return
    }
    if resp.StatusCode != http.StatusOK {
        //fmt.Println(resp.Status)
        http.Error(w, resp.Status, resp.StatusCode)
        return
    }
}
