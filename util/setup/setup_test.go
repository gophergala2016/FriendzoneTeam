package setup

import (
	"flag"
	"fmt"
	"strings"
	"testing"
)

func impresion() {
	fmt.Println("---")
	var spc = "spc scripts/go.sh $u@$h /tmp/go.sh"
	spc = strings.Replace(spc, "$u", *user, 1)
	spc = strings.Replace(spc, "$h", *host, 1)
	fmt.Println(spc)
	fmt.Println(*host)
	fmt.Println(*user)
	fmt.Println(*pwd)
	fmt.Println("---")
}

func Test_load1(t *testing.T) {
	load()
	impresion()
	if *host != "127.0.0.1" {
		t.Fail()
	}
}
func Test_load2(t *testing.T) {
	flag.Set("h", "192.168.1.254")
	load()
	impresion()
	if *host == "127.0.0.1" {
		t.Fail()
	}
}
func Test_load3(t *testing.T) {
	flag.Set("h", "192.168.1.3")
	flag.Set("u", "Chemasmas")
	load()
	impresion()
	if *host == "127.0.0.1" {
		t.Fail()
	}
}
func Test_load4(t *testing.T) {
	flag.Set("h", "192.168.1.4")
	flag.Set("u", "Chemasmas")
	flag.Set("p", "GopherGala")
	load()
	impresion()
	if *host == "127.0.0.1" {
		t.Fail()
	}
}
