package notifo

import (
    "bufio";
    "encoding/base64";
    "fmt";
    "http";
    "io";
    "net";
    "os";
    "strings";
)

type badStringError struct {
    what string
    str  string
}

func (e *badStringError) String() string { return fmt.Sprintf("%s %q", e.what, e.str) }

func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

type readClose struct {
    io.Reader;
	io.Closer;
}

type NotifoApiConn struct {
    root         string;
    apiusername  string;
    apisecret    string;
    client      *http.ClientConn;
    verbose      bool;
}

func NewNotifoApiConn(apiusername string, apisecret string) *NotifoApiConn {
    return &NotifoApiConn {
        root:           "https://api.notifo.com/v1/",
        apiusername:    apiusername,
        apisecret:      apisecret,
        client:         nil,
        verbose:        false,
    }
}

func (api *NotifoApiConn) SubscribeUser(username string) (*http.Response, os.Error) {
    var data map[string]string;
    data["username"] = username;
    
    return api.submitRequest("subscribe_user", "POST", data);
}

func (api *NotifoApiConn) SendNotification(to string, msg string, label string,
        title string, uri string) (*http.Response, os.Error) {
            
    var data map[string]string;

    // TODO: Fail if to is blank.
    data["to"] = to;
    
    // TODO: Fail if msg is blank.
    data["msg"] = msg;
    if len(label) > 0 {
        data["label"] = label;
    }
    if len(title) > 0 {
        data["title"] = title;
    }
    if len(uri) > 0 {
        data["uri"] = uri;
    }
    
    return api.submitRequest("send_notification", "POST", data);    
}

func (api *NotifoApiConn) submitRequest(action string, method string,
        data map[string]string) (*http.Response, os.Error) {

    var request  *http.Request;
    var response *http.Response;
    var err       os.Error;

    rawurl := strings.Join([]string { api.root, action, makeQueryString(data) }, "");
    
    if request, err = prepareRequest(api.apiusername, api.apisecret, rawurl, method); err != nil {
        return nil, err;
    }
    
    if api.verbose {
        dump, _ := http.DumpRequest(request, true);
        print(string(dump));
    }
    
    if response, err = send(request); err != nil {
        return nil, err;
    }

    return response, nil;
}

func prepareRequest(username string, secret string, rawurl string, method string) (*http.Request, os.Error) {
    var request  http.Request;
    var url     *http.URL;
    var err      os.Error;
    
    if url, err = http.ParseURL(rawurl); err != nil {
        return nil, err;
    }
    
    userinfo := strings.Join([]string { username, secret }, ":");
    enc      := base64.URLEncoding;
 	encoded  := make([]byte, enc.EncodedLen(len(userinfo)));
 	
 	enc.Encode(encoded, []byte(userinfo));
 	
 	request.Header = make(map[string]string);
 	request.Header["Authorization"] = "Basic " + string(encoded);
    request.Method = method;
    request.URL = url;
    request.URL.Userinfo = userinfo;

    return &request, nil;
}

func makeQueryString(data map[string]string) string {
    args := "";
    sep  := "?";

    for key, value := range data {
        if len(args) > 0 {
            sep = "&";
        }
        
        args += fmt.Sprintf("%s%s=%s", sep, key, value);
    }
    
    return args; 
}

func send(req *http.Request) (resp *http.Response, err os.Error) {
    if req.URL.Scheme != "http" {
        return nil, &badStringError{"unsupported protocol scheme", req.URL.Scheme}
    }

    addr := req.URL.Host
    if !hasPort(addr) {
        addr += ":http"
    }
    info := req.URL.Userinfo
    if len(info) > 0 {
        enc := base64.URLEncoding
        encoded := make([]byte, enc.EncodedLen(len(info)))
        enc.Encode(encoded, []byte(info))
        if req.Header == nil {
                req.Header = make(map[string]string)
        }
        req.Header["Authorization"] = "Basic " + string(encoded)
    }
    conn, err := net.Dial("tcp", "", addr)
    if err != nil {
        return nil, err
    }

    err = req.Write(conn)
    if err != nil {
        conn.Close()
        return nil, err
    }

    reader := bufio.NewReader(conn)
    resp, err = http.ReadResponse(reader, req.Method)
    if err != nil {
        conn.Close()
        return nil, err
    }

    resp.Body = readClose{resp.Body, conn}

    return
}