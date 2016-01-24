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
	"github.com/ChimeraCoder/anaconda"
	"github.com/creamdog/gonfig"
	"os"
	//    util "github.com/gophergala2016/FriendzoneTeam/util/dateformat"
)

//tipo
type arrMessage []anaconda.DirectMessage
type comandos struct {
	comando       string
	uso           string
	instrucciones []string
}

type arrComandos []comandos

const ESPACIO = ""
const NOESPACIOS = "[\\S]+"
const NOMBRE_ARCHIVO = "[a-zA-z./_0-9-~]+|(\"[a-zA-z./_0-9-\\s~]+\")"

var flag = false

func getMessages(jsonS string) (mesages []anaconda.DirectMessage) {
	var mensajes arrMessage
	data := []byte(jsonS)
	json.Unmarshal(data, &mensajes)
	//fmt.Println(mensajes)
	//for i := 0; i < len(mensajes); i++ {
	processMessages(mensajes)
	//}
	return
}

func processMessages(messages []anaconda.DirectMessage) {

	for i := 0; i < len(messages); i++ {
		estatus, comando := interpretar(messages[i].Text)
		fmt.Println("----")
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
		fmt.Println("Se murio")
		fmt.Println(err)
	}

	//fmt.Println(arrcadenas[0])
	switch {
	case arrcadenas[0] == "create":
		if len(arrcadenas) == 2 && testRegexp(NOMBRE_ARCHIVO, arrcadenas[1]) {
			str, _ = dicc.GetString("comandos/create/v1/comando", nil)
			//fmt.Println("Exito")
			//fmt.Println(str)
			return true, completeRegExp(str, arrcadenas)
		} else if len(arrcadenas) == 3 && testRegexp(NOMBRE_ARCHIVO, arrcadenas[1]) {
			str, _ = dicc.GetString("comandos/create/v2/comando", nil)
			//fmt.Println("Exito")
			//fmt.Println(str)
			return true, completeRegExp(str, arrcadenas)
		} else {
			return false, ""
			//fmt.Println("MAL")
		}
	}

	//fmt.Println(str)

	//	for i := 0; i < len(arrcadenas); i++ {
	//		fmt.Println(arrcadenas[i])
	//	}
	//comment

	return false, ""
}

func testRegexp(exp, val string) (valor bool) {
	var reg *regexp.Regexp
	reg, _ = regexp.Compile(exp)
	return reg.MatchString(val)
}

func completeRegExp(ssh string, param []string) (truessh string) {
	var tssh string
	if strings.Contains(ssh, "$2") && len(param) >= 3 {
		//Si hay un placeholder %2, remplaza ese
		//fmt.Println("ok")
		tssh = strings.Replace(ssh, "$1", param[2], 1)
		tssh = strings.Replace(tssh, "$2", param[1], 1)
		tssh = strings.Replace(tssh, "//", "/", -1)
	} else {
		tssh = strings.Replace(tssh, "$1", param[1], 1)
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
