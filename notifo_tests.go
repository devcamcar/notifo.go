package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"notifo"
)

func main() {
	var response *http.Response
	var err error

	api := notifo.New("gotest", "a25c4f206494150bddf2e716705c8bedcad0cb16")
	//api.SetEndpoint("http://localhost:8000/v1/");

	if response, err = api.SubscribeUser("devcamcar"); err != nil {
		log.Stderr("ERROR")
		log.Stderr(err)
	} else {
		dump, _ := httputil.DumpResponse(response, true)
		log.Stdout(response.StatusCode)
		log.Stdout(string(dump))
	}
}
