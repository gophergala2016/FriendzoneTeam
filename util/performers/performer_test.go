package performer

import (
	"io/ioutil"
	//"strings"
	"testing"
)

func Test_getMessages(t *testing.T) {
	//Leo el dato de prueba de un archivo
	data, _ := ioutil.ReadFile("test.json")
	dataS := string(data)
	getMessages(dataS)

	//if strings.Compare(dataS, "") {
	//t.Error(...)
	//}
}
