package main

import ( 
    "net/http"
    service "github.com/gophergala2016/FriendzoneTeam/services"
)

func main()  {
    http.HandleFunc("/test", service.RevisarDM)
    http.ListenAndServe(":3000", nil)
}