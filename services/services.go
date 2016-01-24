package service

import (
    "fmt"
    "log"
    "time"
    "net/http"
    "net/url"
    "encoding/json"
    "github.com/ChimeraCoder/anaconda"
    util "github.com/gophergala2016/FriendzoneTeam/util/dateformat"
)

// Regresa los DMs que se han recibido en base a la fecha actual
func RevisarDM(w http.ResponseWriter, r *http.Request){
    // Verifica si ha obtenido nuevos DMs
    anaconda.SetConsumerKey("lCqA4GsOhivuJumCMynVuOI2B")
    anaconda.SetConsumerSecret("B3XOc0n1FLw0faxl2SSCC7xNFxAAWdPnL7shzLj0Sq3l3OUqvE")
    api := anaconda.NewTwitterApi("4804703832-9fNj7vcJBobyDdbYhDPKOfFbTJzXkq64VFr99Qr", "aK3dWB3HoM01p79UFRTD8eQh9SYCwElr1RCicqe3imHBf")// Revisamos que ejecute bien la conexion con las credenciales
    if api.Credentials == nil {
        log.Println("Twitter Api client has empty (nil) credentials")
    }
    uri := url.Values{}
    dmResults, err := api.GetDirectMessages(uri)
    if err != nil {
        log.Printf("Error: %s", err.Error())
    }
    
    output, err := json.Marshal(dmResults)
    if err != nil {
        log.Printf("Error: %s", err.Error())
    }
    
    cd := time.Now()
    currentDate := fmt.Sprintf("%d-%s-%d\n", cd.Day(), cd.Month(), cd.Year())
    i := 0
    for _, message := range dmResults {
        strFormat, err := util.DateFormat(message.CreatedAt)
        if err != nil {
            log.Fatal(err.Error())
        }
        fmt.Printf("%s\n", message.CreatedAt)
        fmt.Printf("%s - %s\n", currentDate, strFormat)
        if currentDate == strFormat {
            i++
        }
        // fmt.Printf("%d", i)
    }
    
    fmt.Fprintf(w, string(output))
}

func createupdate() {
    
}

func Command(){
    
}