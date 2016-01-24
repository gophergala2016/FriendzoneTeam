package performer

import (
	"encoding/json"
	"fmt"
	"github.com/creamdog/gonfig"
	"os"
	"regexp"
	"strings"
)

type Scheduler struct {
    DmId string `json:"id_dm" bson:"id_dm"`
    Command string `json:"command" bson:"command"`
    UserId string `json:"user_id" bson:"user_id"`
    Status string `json:"status" bson:"status"`
    Created_At string `json:"created_at" bson:"created_at"`
}

type ArrMessage []Scheduler
type Comandos struct {
	Status  bool
	Command string
}
type ArrComandos []Comandos

const ESPACIO = ""
const NOESPACIOS = "[\\S]+"
const NOMBRE_ARCHIVO = "[a-zA-z./_0-9-~]+"

var flag = false

func GetMessages(jsonS string) []Comandos {
	var mensajes ArrMessage
	data := []byte(jsonS)
	err := json.Unmarshal(data, &mensajes)
    if err != nil {
        fmt.Printf("GetMessages %s", err)
    }
	resultado := ProcessMessages(mensajes)
	// fmt.Println(resultado)
	return resultado
}

func ProcessMessages(messages []Scheduler) []Comandos {

	var lista = make([]Comandos, len(messages))
	for i := 0; i < len(messages); i++ {
		estatus, comando := interpretar(messages[i].Command)
		lista[i].Status = estatus
		lista[i].Command = comando
	}
	return lista
}

func interpretar(comando string) (bool, string) {

	var str string
	arrcadenas := strings.Split(comando, " ")
    // fmt.Printf("%s\n", arrcadenas)
    
	//Comandos
	f, err := os.Open("C:\\desarrollo\\ws_go\\src\\github.com\\gophergala2016\\FriendzoneTeam\\util\\performer\\comandos.json")
    if err != nil {
		fmt.Printf("Error abrir archivo: %s\n",err)
	}
	defer f.Close()

	dicc, err := gonfig.FromJson(f)

	if err != nil {
		fmt.Printf("Error abrir archivo: %s\n",err)
	}
    
	switch {
	case arrcadenas[0] == "create":
		str = create(dicc, arrcadenas)
	case arrcadenas[0] == "delete":
		str = delete(dicc, arrcadenas)
	case arrcadenas[0] == "move":
		str = move(dicc, arrcadenas)
	case arrcadenas[0] == "rename":
		str = rename(dicc, arrcadenas)
	case arrcadenas[0] == "server":
		str = server(dicc, arrcadenas)
	case strings.HasPrefix(arrcadenas[0], ":"):
		str = custom(arrcadenas)
	}

	if str != "" {
		return true, str
	}
	return false, "Comando Invalido"
}

func custom(arrcadenas []string) string {
	var str string
	str = strings.Join(arrcadenas, " ")
	str = strings.TrimPrefix(str, ":")
	return str
}

func server(dicc gonfig.Gonfig, arrcadenas []string) string {
	var str string
	switch {
	case arrcadenas[1] == "new":
		switch {
		case arrcadenas[2] == "go":
			str, _ = dicc.GetString("comandos/server/new/go/comando", nil)
		case arrcadenas[2] == "lamp":
			str, _ = dicc.GetString("comandos/server/new/lamp/comando", nil)
		case arrcadenas[2] == "lemp":
			str, _ = dicc.GetString("comandos/server/new/lemp/comando", nil)
		case arrcadenas[2] == "mean":
			str, _ = dicc.GetString("comandos/server/new/mean/comando", nil)
		}

	case arrcadenas[1] == "start":
		str, _ = dicc.GetString("comandos/server/start/comando", nil)
	case arrcadenas[1] == "restart":
		str, _ = dicc.GetString("comandos/server/restart/comando", nil)
	case arrcadenas[1] == "stop":
		str, _ = dicc.GetString("comandos/server/stop/comando", nil)
	}

	return completeRegExp(str, arrcadenas)
}

func rename(dicc gonfig.Gonfig, arrcadenas []string) string {
	var str string
	if len(arrcadenas) == 3 && testRegexp(NOMBRE_ARCHIVO, arrcadenas[1]) && testRegexp(NOMBRE_ARCHIVO, arrcadenas[2]) {
		if strings.HasSuffix(arrcadenas[1], "/") {
			str, _ = dicc.GetString("comandos/rename/v2/comando", nil)
		} else {
			str, _ = dicc.GetString("comandos/rename/v1/comando", nil)
		}
		return completeRegExp(str, arrcadenas)
	} else {
		str, _ = dicc.GetString("comandos/rename/v1/error", nil)
		return str
	}
}

func move(dicc gonfig.Gonfig, arrcadenas []string) string {
	var str string
	if len(arrcadenas) == 3 && testRegexp(NOMBRE_ARCHIVO, arrcadenas[1]) {
		if strings.HasSuffix(arrcadenas[1], "/") {
			str, _ = dicc.GetString("comandos/move/v2/comando", nil)
		} else {
			str, _ = dicc.GetString("comandos/move/v1/comando", nil)
		}
		return completeRegExp(str, arrcadenas)
	} else {
		str, _ = dicc.GetString("comandos/move/v1/error", nil)
		return str
	}
}

func delete(dicc gonfig.Gonfig, arrcadenas []string) string {
	var str string
	if len(arrcadenas) == 2 && testRegexp(NOMBRE_ARCHIVO, arrcadenas[1]) {
		if strings.HasSuffix(arrcadenas[1], "/") {
			str, _ = dicc.GetString("comandos/delete/v2/comando", nil)
		} else {
			str, _ = dicc.GetString("comandos/delete/v1/comando", nil)
		}
		return completeRegExp(str, arrcadenas)
	} else {
		str, _ = dicc.GetString("comandos/delete/v1/error", nil)
		return str
	}
}

func create(dicc gonfig.Gonfig, arrcadenas []string) string {
	var str string
	if len(arrcadenas) == 2 && testRegexp(NOMBRE_ARCHIVO, arrcadenas[1]) {
		if strings.HasSuffix(arrcadenas[1], "/") {
			str, _ = dicc.GetString("comandos/create/v3/comando", nil)
		} else {
			str, _ = dicc.GetString("comandos/create/v1/comando", nil)
		}
		return completeRegExp(str, arrcadenas)
	} else if len(arrcadenas) == 3 && testRegexp(NOMBRE_ARCHIVO, arrcadenas[1]) {
		if strings.HasSuffix(arrcadenas[1], "/") {
			str, _ = dicc.GetString("comandos/create/v4/comando", nil)
		} else {
			str, _ = dicc.GetString("comandos/create/v2/comando", nil)
		}
		return completeRegExp(str, arrcadenas)
	} else {
		str, _ = dicc.GetString("comandos/create/v1/error", nil)
		return str
	}
}

func testRegexp(exp, val string) bool {
	var reg *regexp.Regexp
	reg, _ = regexp.Compile(exp)
	return reg.MatchString(val)
}

func completeRegExp(ssh string, param []string) string {
	var tssh string

	if len(param) >= 3 {
		//Si hay un placeholder %2, remplaza ese
		//fmt.Println("ok")
		tssh = strings.Replace(ssh, "$1", param[2], -1)
		tssh = strings.Replace(tssh, "$2", param[1], -1)
		tssh = strings.Replace(tssh, "//", "/", -1)
	} else {
		tssh = strings.Replace(ssh, "$1", param[1], -1)
		tssh = strings.Replace(tssh, "//", "/", -1)
	}
	return tssh
}
