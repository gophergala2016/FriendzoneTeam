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
    "upper.io/db.v2"
    "upper.io/db.v2/mongo"
)

type Scheduler struct {
    DmId int16 `json:"id_dm" bson:"id_dm"`
    Command string `json:"command" bson:"command"`
    Type string `json:"type" bson:"type"`
    UserId string `json:"user_id" bson:"user_id"`
    Created_At string `json:"created_at" bson:"created_at"`
}

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
    var settings = mongo.ConnectionURL{
        Address:  db.Host("ds049945.mongolab.com:49945/socialgopher"), // MongoDB hostname.
        Database: "socialgopher",            // Database name.
        User:     "friendzonedb",             // Optional user name.
        Password: "friendzonedb",             // Optional user password.
    }
    sess, err := db.Open(mongo.Adapter, settings)
    var regs []Scheduler
    if err != nil {
        log.Fatalf("db.Open(): %q\n", err)
    }
    defer sess.Close()
    // Scheduler
    schedulerCollection, err := sess.Collection("scheduler")
    if err != nil {
        log.Fatalf("Could not use collection: %q\n", err)
    }
    var res db.Result
    res = schedulerCollection.Find().Where("type = ?", "SHOW_FILE")
    err = res.All(&regs)
    if err != nil {
        log.Fatalf("res.All(): %q\n", err)
    }
    
    outputw, err := json.Marshal(regs)
    
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
    fmt.Printf("%s", string(output))
    fmt.Fprintf(w, string(outputw))
}

func createupdate() {
    
}

func Command(){
    
}