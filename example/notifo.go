package main

import (
	"flag"
	"fmt"
	"github.com/devcamcar/notifo.go"
	"os"
	"strings"
)

var apiuser, apisecret, to, label, title, url string
var verbose bool

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Verbose response")
	flag.StringVar(&apiuser, "user", "", "Your API username")
	flag.StringVar(&apisecret, "secret", "", "Your API secret")
	flag.StringVar(&to, "to", "", "recipient")
	flag.StringVar(&label, "label", "notifo.go", "msg label")
	flag.StringVar(&title, "title", "", "msg title")
	flag.StringVar(&url, "url", "", "A URL to send")
}

func main() {
	flag.Parse()

	n := notifo.New(apiuser, apisecret)
	msg := strings.Join(flag.Args(), " ")

	rv, err := n.SendNotification(to, msg, label, title, url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message:  %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Response:  %#v\n", rv)
	}
}
