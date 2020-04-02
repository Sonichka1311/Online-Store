package main

import (
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "shop/pkg/auth"
    "shop/pkg/constants"
    "shop/pkg/handlers"
    "shop/pkg/logic"
    "shop/pkg/models"
    "shop/pkg/product"
    "strconv"
    "sync"
)

func main() {
    databaseConnector := models.Connector{
        Router: models.Router{Host: logic.GetUrl(constants.Protocol, constants.DatabaseHost, constants.DatabasePort)},
        Mutex:  sync.Mutex{},
    }

    authConnector := models.Connector{
        Router: models.Router{Host: logic.GetUrl(constants.Protocol, constants.AuthHost, constants.AuthPort)},
        Mutex:  sync.Mutex{},
    }

    handler := handlers.ProductHandler{
        Repo:   &product.Repo{
            Connector:  &databaseConnector,
        },
        Auth:   &auth.Repo{
            Connector:  &authConnector,
        },
    }

    router := mux.NewRouter()
    router.HandleFunc("/", handler.GetProductsList).Methods(http.MethodGet)
    router.HandleFunc("/product", handler.ProductCard).Methods(http.MethodGet)
    router.HandleFunc("/product", handler.AddProduct).Methods(http.MethodPost)
    router.HandleFunc("/product", handler.EditProduct).Methods(http.MethodPut)
    router.HandleFunc("/product", handler.DeleteProduct).Methods(http.MethodDelete)

    http.Handle("/", router)

    fs := http.FileServer(http.Dir("./swaggerui"))
    http.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", fs))

    err := http.ListenAndServe(":" + strconv.Itoa(constants.MainServerPort), nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
