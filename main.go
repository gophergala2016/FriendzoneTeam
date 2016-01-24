package main

import ( 
    "fmt"
    "time"
    "net/http"
    "net/url"
    "encoding/json"
    "github.com/ChimeraCoder/anaconda"
    util "github.com/gophergala2016/FriendzoneTeam/util/dateformat"
)

func main()  {
    http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request){
        // Verifica si ha obtenido nuevos DMs
        anaconda.SetConsumerKey("qFwSZ5VTSw2E23u166iWPAZj1")
        anaconda.SetConsumerSecret("AgJPsgKrP1OW3fx9TpDwZ55beUMGdBLcybJLyUvgMJM0tQj5PD")
        api := anaconda.NewTwitterApi("151138456-QnyorQfe59zld78ATJdfwNpcpsywVPZqcqo5uXiZ", "TpG0upWDAAmXF81K0HscQOxXCqHfpWEYd66c7fVdLX1k8")
        // var dms []anaconda.DirectMessage
        uri := url.Values{}
        if api.Credentials == nil {
            fmt.Println("Twitter Api client has empty (nil) credentials")
        }
        searchResult, _ := api.GetDirectMessages(uri)
        output, err := json.Marshal(searchResult)
        
        if err != nil {
            fmt.Println("Error parseando JSON :c")
        }
        
        // Envio DMs
        /* message, _ := api.PostDMToScreenName("chemas sin json", "Chemasmas")
        if err != nil {
            fmt.Println("Error enviando el DM :c")
        }
        fmt.Println(message)*/
        currentDate := time.Now()
        fmt.Printf("%d-%s-%d\n", currentDate.Day(), currentDate.Month(), currentDate.Year())
        strdate, _ := util.DateFormat("Sat Jan 23 19:46:10 +0000 2016")
        fmt.Println(strdate)
        
        fmt.Fprintf(w, string(output))
    })
    
    http.ListenAndServe(":3000", nil)
}