package setup

import (
	"flag"
	"fmt"
	"strings"
)

var Host *string
var User *string
var Pwd *string
var Spc string

//Valores Default

func init() {
	host = flag.String("h", "127.0.0.1", "Host")
	user = flag.String("u", "dummy", "User")
	pwd = flag.String("p", "123456", "Password")
	spc = "spc scripts/go.sh $u@$h /tmp/go.sh"
}

func Load() (string, string, string) {
	flag.Parse()
	spc = strings.Replace(spc, "$u", *user, 1)
	spc = strings.Replace(spc, "$h", *host, 1)
	fmt.Println(spc)
	return *host, *user, *pwd
}
