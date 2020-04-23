package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "shop/pkg/auth"
    "shop/pkg/constants"
    "shop/pkg/database"
    "shop/pkg/handlers"
    "shop/pkg/logic"
    "shop/pkg/models"
    "shop/pkg/product"
    "strconv"
    "sync"
    "time"
)

func main() {
    db, dbError := sql.Open("mysql", "root:guest@tcp(mysql:3306)/shop?charset=utf8&interpolateParams=true")
    if dbError != nil {
        log.Fatalf("Cannot open database: %s", dbError.Error())
    }

    for tries := 0; tries < 10; tries++ {
        dbError = db.Ping()
        if dbError == nil {
            break
        }
        log.Printf("Failed connect to database for %d times. Trying to reconnect...", tries + 1)
        time.Sleep(3 * time.Second)
    }
    if dbError != nil {
        log.Fatalf("Cannot connect to database: %s", dbError.Error())
    }

    authConnector := models.Connector{
        Router: models.Router{Host: logic.GetUrl(constants.Protocol, constants.AuthHost, constants.AuthPort)},
        Mutex:  sync.Mutex{},
    }

    handler := handlers.ProductHandler{
        Repo:   &product.Repo{
            Connector:  database.NewConnector(db),
        },
        Auth:   &auth.Repo{
            Connector:  &authConnector,
        },
        Mutex: &sync.RWMutex{},
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
