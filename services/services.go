package service

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "encoding/json"
    "github.com/ChimeraCoder/anaconda"
    "upper.io/db.v2"
    "upper.io/db.v2/mongo"
    util "github.com/gophergala2016/FriendzoneTeam/util/performer"
    ssh "github.com/gophergala2016/FriendzoneTeam/ssh"
)

type Scheduler struct {
    DmId string `json:"id_dm" bson:"id_dm"`
    Command string `json:"command" bson:"command"`
    UserId string `json:"user_id" bson:"user_id"`
    Status string `json:"status" bson:"status"`
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
    
    var settings = mongo.ConnectionURL{
        Address:  db.Host("ds049945.mongolab.com:49945"), // MongoDB hostname.
        Database: "socialgopher",            // Database name.
        User:     "friendzonedb",             // Optional user name.
        Password: "friendzonedb",             // Optional user password.
    }
    sess, err := db.Open(mongo.Adapter, settings)
    if err != nil {
        log.Fatalf("db.Open(): %q\n", err)
    }
    defer sess.Close()
    // Scheduler
    lCollection, err := sess.Collections()
    if len(lCollection) != 0 {
        log.Println("Existe la collection")
        schedulerCollection, err := sess.Collection("scheduler")
        if err != nil {
            log.Fatalf("Could not use collection: %q\n", err)
        }
        var scheduler []Scheduler
        for _, message := range dmResults {
            var res db.Result
            reg := new(Scheduler)
            res = schedulerCollection.Find().Where("id_dm = ?", message.IdStr,)
            err = res.One(reg)
            if err != nil {
                log.Fatalf("res.All(): %q\n", err)
            }else{
                // log.Println(err)
                // log.Println(message.IdStr)
                // log.Println(reg)
                
                // No existe el registro en la Base
                if reg.DmId == "" {
                    log.Println("No existe en la Base")
                    reg.DmId = message.IdStr
                    reg.Created_At = message.CreatedAt
                    reg.Status = "Queue"
                    reg.UserId = message.SenderScreenName
                    reg.Command = message.Text
                    schedulerCollection.Append(reg)
                }
            }
        }
        // Obtenemos todos los mensajes
        var results db.Result
        results = schedulerCollection.Find().Where("status = ?", "Queue")
        err = results.All(&scheduler)
        output, err := json.Marshal(scheduler)
        if err != nil {
            log.Printf("Error: %s", err.Error())
        }
        go func(){
            var comandos util.ArrComandos
            comandos = util.GetMessages(string(output))
            for _, comando := range comandos {
                if comando.Status {
                    out, err := ssh.Conekta("gophers", "gophers", "191.233.33.24", comando.Command)
                    if err != nil {
                        log.Fatalf("Run failed: %s", err)
                    }
                    if out == "" {
                        log.Fatalf("Output was empty for command: %s", comando.Command)
                    }
                }
            }
        }()
        fmt.Fprintf(w, string(output))
    }else {
        log.Println("No existe la collection")
        schedulerCollection, err := sess.Collection("scheduler")
        var scheduler []Scheduler
        for _, message := range dmResults {
            var reg Scheduler
            reg.DmId = message.IdStr
            reg.Created_At = message.CreatedAt
            reg.Status = "Queue"
            reg.UserId = message.SenderScreenName
            reg.Command = message.Text
            schedulerCollection.Append(reg)
        }
        
        // Obtenemos todos los mensajes
        var results db.Result
        results = schedulerCollection.Find()
        err = results.All(&scheduler)
        output, err := json.Marshal(scheduler)
        if err != nil {
            log.Printf("Error: %s", err.Error())
        }
        go func(){
           comandos := util.GetMessages(string(output)) 
           for _, comando := range comandos {
               fmt.Print(comando)
           }
        }()
        fmt.Fprintf(w, string(output))
    }
}

func Createupdate() {
    
}

func Command(){
    
}