package main

import (
    "./handlers"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", handlers.GetProductsList)
    http.HandleFunc("/product", handlers.ProductCard)
    http.HandleFunc("/add_product", handlers.AddProduct)
    http.HandleFunc("/edit_product", handlers.EditProduct)
    http.HandleFunc("/delete_product", handlers.DeleteProduct)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
