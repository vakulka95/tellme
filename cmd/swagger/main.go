package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/tidwall/pretty"

	"gitlab.com/tellmecomua/tellme.api/pkg/swagger"
)

func main() {
	bb, err := swagger.GetApidocsJSON()
	if err != nil {
		log.Fatalf("failed to get apidocs json: %v", err)
	}

	err = ioutil.WriteFile("./apidocs/swagger.json", pretty.Pretty(bb), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to write apidocs file: %v", err)
	}
}
