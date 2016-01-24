package performer

import (
	"fmt"
	"regexp"
	//	"io/ioutil"
	"strings"
	//    "log"
	//    "time"
	//    "net/http"
	//    "net/url"
	"encoding/json"
	//	"github.com/ChimeraCoder/anaconda"
	"github.com/creamdog/gonfig"
	service "github.com/gophergala2016/FriendzoneTeam/services"
	"os"
	//    util "github.com/gophergala2016/FriendzoneTeam/util/dateformat"
)

//tipo
//type arrMessage []anaconda.DirectMessage
type arrMessage []service.Scheduler
type comandos struct {
	comando       string
	uso           string
	instrucciones []string
}
type arrComandos []comandos

const ESPACIO = ""
const NOESPACIOS = "[\\S]+"
const NOMBRE_ARCHIVO = "[a-zA-z./_0-9-~]+"

var flag = false

//func getMessages(jsonS string) (mesages []anaconda.DirectMessage) {
func getMessages(jsonS string) (mesages []service.Scheduler) {
	var mensajes arrMessage
	data := []byte(jsonS)
	json.Unmarshal(data, &mensajes)
	//fmt.Println(mensajes)
	//for i := 0; i < len(mensajes); i++ {
	processMessages(mensajes)
	//}
	return
}

//func processMessages(messages []anaconda.DirectMessage) {
func processMessages(messages []service.Scheduler) {

	for i := 0; i < len(messages); i++ {
		//estatus, comando := interpretar(messages[i].Text)
		estatus, comando := interpretar(messages[i].Command)
		fmt.Println("----")
		fmt.Println(messages[i].Command)
		fmt.Println(estatus)
		fmt.Println(comando)
		fmt.Println("----")
	}
	return
}

func interpretar(comando string) (status bool, ssh string) {

	var str string

	arrcadenas := strings.Split(comando, " ")

	//Comandos
	f, _ := os.Open("comandos.json")

	defer f.Close()

	dicc, err := gonfig.FromJson(f)

	//var cli comandos

	if err != nil {
		//fmt.Println("Se murio")
		fmt.Println(err)
	}

	//fmt.Println(arrcadenas[0])
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
	return false, str
}

func custom(arrcadenas []string) (comando string) {
	var str string
	str = strings.Join(arrcadenas, " ")
	str = strings.TrimPrefix(str, ":")
	return str
}

func server(dicc gonfig.Gonfig, arrcadenas []string) (comando string) {
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

func rename(dicc gonfig.Gonfig, arrcadenas []string) (comando string) {
	var str string
	if testRegexp(NOMBRE_ARCHIVO, arrcadenas[1]) {
		if strings.HasSuffix(arrcadenas[1], "/") {
			str, _ = dicc.GetString("comandos/rename/v2/comando", nil)
		} else {
			str, _ = dicc.GetString("comandos/rename/v1/comando", nil)
		}
		return completeRegExp(str, arrcadenas)
	} else {
		return ""
	}
}

func move(dicc gonfig.Gonfig, arrcadenas []string) (comando string) {
	var str string
	if testRegexp(NOMBRE_ARCHIVO, arrcadenas[1]) {
		if strings.HasSuffix(arrcadenas[1], "/") {
			str, _ = dicc.GetString("comandos/move/v2/comando", nil)
		} else {
			str, _ = dicc.GetString("comandos/move/v1/comando", nil)
		}
		return completeRegExp(str, arrcadenas)
	} else {
		return ""
	}
}

func delete(dicc gonfig.Gonfig, arrcadenas []string) (comando string) {
	var str string
	if testRegexp(NOMBRE_ARCHIVO, arrcadenas[1]) {
		if strings.HasSuffix(arrcadenas[1], "/") {
			str, _ = dicc.GetString("comandos/delete/v2/comando", nil)
		} else {
			str, _ = dicc.GetString("comandos/delete/v1/comando", nil)
		}
		return completeRegExp(str, arrcadenas)
	} else {
		return ""
	}
}

func create(dicc gonfig.Gonfig, arrcadenas []string) (comando string) {
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
		return ""
	}
}

func testRegexp(exp, val string) (valor bool) {
	var reg *regexp.Regexp
	reg, _ = regexp.Compile(exp)
	return reg.MatchString(val)
}

func completeRegExp(ssh string, param []string) (truessh string) {
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
	/*
		fmt.Println("---")
		fmt.Println(strings.Contains(ssh, "$2"))
		fmt.Println(len(param) >= 3)
		fmt.Println(len(param))
		fmt.Println(ssh)
		fmt.Println(tssh)
		fmt.Println("---")
	*/
	return tssh
}
