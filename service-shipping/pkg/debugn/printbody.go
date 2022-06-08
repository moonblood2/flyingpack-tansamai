package debugn

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func PrintBody(body io.ReadCloser) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(bytes))
}
