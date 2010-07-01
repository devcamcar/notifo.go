package notifo

import (
    "http"
    "os"
    "restful"
)

type NotifoApiClient struct {
    endpoint    string;
    apiusername string;
    apisecret   string;
}

func New(apiusername, apisecret string) (*NotifoApiClient) {
    return &NotifoApiClient {
        endpoint:       "https://api.notifo.com/v1/",
        apiusername:    apiusername,
        apisecret:      apisecret,
    }
}

func (api *NotifoApiClient) SetEndpoint(endpoint string) {
    api.endpoint = endpoint;
}

func (api *NotifoApiClient) SubscribeUser(username string) (*http.Response, os.Error) {
    data := make(map[string]string);
    data["username"] = username;
    
    return api.submitRequest("subscribe_user", "POST", data);
}

func (api *NotifoApiClient) SendNotification(to string, msg string, label string,
        title string, uri string) (*http.Response, os.Error) {    
    data := make(map[string]string);

    // TODO(devcamcar): Fail if to is blank.
    data["to"] = to;
    
    // TODO(devcamcar): Fail if msg is blank.
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

func (api *NotifoApiClient) submitRequest(action, method string, params map[string]string) (*http.Response, os.Error ) {
    var resp *http.Response
    var err   os.Error

    client := &restful.RestClient {
        Endpoint: api.endpoint,
        UserInfo: api.apiusername + ":" + api.apisecret,
    }
    
    if resp, err = client.SubmitRequest(action, method, nil, params); err != nil {
       return nil, err
    } 
    
    return resp, nil
}
