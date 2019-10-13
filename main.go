package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	handler "github.com/m-buntoro/go-face-webserver/api/http"
	"github.com/m-buntoro/go-face-webserver/pkg/config"
)

func main() {

	var config config.Config

	out, err := ioutil.ReadFile("files/config.json")
	if err != nil {
		log.Println(err)
		return
	}
	if err := json.Unmarshal(out, &config); err != nil {
		log.Println(err)
		return
	}

	rtr := handler.New(handler.Config{
		Conf: config,
	})

	log.Println("Begining serve on port:", config.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), rtr.Handler())
	if err != nil {
		log.Println(err)
	}

	defer func() {
		log.Println("*laugh as program dies*")
		rtr.Close()
	}()
}
