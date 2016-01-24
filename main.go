package main

import ( 
    "net/http"
    "github.com/gorilla/mux"
    service "github.com/gophergala2016/FriendzoneTeam/services"
)

func main()  {
    app := mux.NewRouter()
    app.HandleFunc("/messages", service.RevisarDM)
    http.Handle("/", app)
    http.ListenAndServe(":3000", nil)
}