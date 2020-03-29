package main

import (
    "../common/constants"
    "./handlers"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "strconv"
)

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", handlers.GetProductsList).Methods(http.MethodGet)
    router.HandleFunc("/product", handlers.ProductCard).Methods(http.MethodGet)
    router.HandleFunc("/product", handlers.AddProduct).Methods(http.MethodPost)
    router.HandleFunc("/product", handlers.EditProduct).Methods(http.MethodPut)
    router.HandleFunc("/product", handlers.DeleteProduct).Methods(http.MethodDelete)

    http.Handle("/", router)

    fs := http.FileServer(http.Dir("./swaggerui"))
    http.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", fs))

    err := http.ListenAndServe(":" + strconv.Itoa(constants.MainServerPort), nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
