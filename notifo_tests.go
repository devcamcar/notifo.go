package main

import (
    "http";
    "log";
    "os";
    "notifo";
)

// TODO(devcamcar): support https in client lib

func main() {
    var response *http.Response;
    var err       os.Error;
    
    api := notifo.NewNotifoApiConn("gotest", "a25c4f206494150bddf2e716705c8bedcad0cb16");
    
    if response, err = api.SubscribeUser("devcamcar"); err != nil {
        log.Stderr(err)
    } else {
        dump, _ := http.DumpResponse(response, true)
        log.Stdout(string(dump))
    }
}        
