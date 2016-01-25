package service

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	ssh "github.com/gophergala2016/FriendzoneTeam/ssh"
	util "github.com/gophergala2016/FriendzoneTeam/util/performer"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"upper.io/db.v2"
	"upper.io/db.v2/mongo"
)

type Scheduler struct {
	DmId       string `json:"id_dm" bson:"id_dm"`
	Command    string `json:"command" bson:"command"`
	UserId     string `json:"user_id" bson:"user_id"`
	Status     string `json:"status" bson:"status"`
	Created_At string `json:"created_at" bson:"created_at"`
}

// Regresa los DMs que se han recibido en base a la fecha actual
func RevisarDM(w http.ResponseWriter, r *http.Request) {
	// Verifica si ha obtenido nuevos DMs
	anaconda.SetConsumerKey("lCqA4GsOhivuJumCMynVuOI2B")
	anaconda.SetConsumerSecret("B3XOc0n1FLw0faxl2SSCC7xNFxAAWdPnL7shzLj0Sq3l3OUqvE")
	api := anaconda.NewTwitterApi("4804703832-9fNj7vcJBobyDdbYhDPKOfFbTJzXkq64VFr99Qr", "aK3dWB3HoM01p79UFRTD8eQh9SYCwElr1RCicqe3imHBf") // Revisamos que ejecute bien la conexion con las credenciales
	if api.Credentials == nil {
		log.Println("Twitter Api client has empty (nil) credentials")
	}
	uri := url.Values{}
	dmResults, err := api.GetDirectMessages(uri)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	//Setting de la coneccion a mongo
	var settings = mongo.ConnectionURL{
		Address:  db.Host("ds049945.mongolab.com:49945"), // MongoDB hostname.
		Database: "socialgopher",                         // Database name.
		User:     "friendzonedb",                         // Optional user name.
		Password: "friendzonedb",                         // Optional user password.
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
			res = schedulerCollection.Find().Where("id_dm = ?", message.IdStr)
			err = res.One(reg)
			if err != nil {
				log.Fatalf("res.All(): %q\n", err)
			} else {
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
		go func(col db.Collection, json string) {
			var comandos util.ArrComandos
			comandos = util.GetMessages(json)
			for _, comando := range comandos {
				if comando.Status {
					status := "Success"
					out, err := ssh.Conekta("gophers", "gophers", "191.233.33.24", comando.Command)
					if err != nil {
						log.Fatalf("Run failed: %s", err)
						status = "Failed"
					}
					if out == "" {
						log.Fatalf("Output was empty for command: %s", comando.Command)
						status = "Failed"
					}
					var wg sync.WaitGroup
                    fmt.Println(status + "\n")
                    fmt.Println(out + "\n")
                    wg.Add(1)
                    /*if strings.HasSuffix(out, ".sh") {
                        wg.Add(2)
                    }else {
                        wg.Add(1)
                    }
					log.Printf("%s", err2)
					/*go func() {
						defer wg.Done()   
						var res2 db.Result
						reg2 := new(Scheduler)
						res2 = col.Find().Where("command = ?", comando.Command)
						err2 := res2.One(reg2)
						reg2.Status = status
						err2 = res2.Update(reg2)
						if err2 != nil {
							log.Printf("%s", err2)
						}
					*ssh/     }()*/
					if strings.HasSuffix(out, ".sh") {
						go func(file string) {
							defer wg.Done()

							var cmds [3]string
							cmds[0] = fmt.Sprintf("scp scripts/%s %s@%s /tmp/%s", file, "root", "191.233.33.24", file)
							cmds[1] = fmt.Sprintf("chmod +x /tmp/%s", file)
							cmds[2] = fmt.Sprintf("./%s", file)

							for _, cmd := range cmds {
								out, err := ssh.Conekta("gophers", "gophers", "191.233.33.24", cmd)
								if err != nil {
									log.Fatalf("Run failed: %s", err)
								}
								if out == "" {
									log.Fatalf("Output was empty for command: %s", cmd)
								}
							}

						}(out)
					}
					wg.Wait()
				}

			}
		}(schedulerCollection, string(output))
		fmt.Fprintf(w, string(output))
	} else {
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
		go func() {
			var comandos util.ArrComandos
			comandos = util.GetMessages(string(output))
			for _, comando := range comandos {
				if comando.Status {
					status := "Success"
					out, err := ssh.Conekta("gophers", "gophers", "191.233.33.24", comando.Command)
					if err != nil {
						log.Fatalf("Run failed: %s", err)
						status = "Failed"
					}
					if out == "" {
						log.Fatalf("Output was empty for command: %s", comando.Command)
						status = "Failed"
					}
					var res db.Result
					reg := new(Scheduler)
					res = schedulerCollection.Find().Where("command = ?", comando.Command)
					err = res.One(reg)
					reg.Status = status
					err = res.Update(reg)
					if err != nil {
						fmt.Printf("%s", err)
					}

					var wg sync.WaitGroup
					wg.Add(1)
					if strings.HasSuffix(out, ".sh") {
						go func(file string) {
							defer wg.Done()

							var cmds [3]string
							cmds[0] = fmt.Sprintf("scp scripts/%s %s@%s /tmp/%s", file, "root", "191.233.33.24", file)
							cmds[1] = fmt.Sprintf("chmod +x /tmp/%s", file)
							cmds[2] = fmt.Sprintf("./%s", file)

							for _, cmd := range cmds {
								out, err := ssh.Conekta("gophers", "gophers", "191.233.33.24", cmd)
								if err != nil {
									log.Fatalf("Run failed: %s", err)
								}
								if out == "" {
									log.Fatalf("Output was empty for command: %s", cmd)
								}
							}

						}(out)
					}
					wg.Wait()
				}
			}
		}()
		fmt.Fprintf(w, string(output))
	}
}

func Createupdate() {

}

func Command() {

}
