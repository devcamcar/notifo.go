package main

import (
    "fmt";
    "http";
    "os";
    "notifo";
)

func main() {
    var response *http.Response;
    var err       os.Error;
    
    api := notifo.NewNotifoApiConn("username", "secret");
    
    if response, err = api.SubscribeUser("devcamcar"); err != nil {
        fmt.Printf("Error: %s\n", err.String());
    } else {
        fmt.Printf("Response status: %d\n", response.StatusCode);
    }
}
