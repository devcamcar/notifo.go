package notifo

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type NotifoApiClient struct {
	endpoint    string
	apiusername string
	apisecret   string
}

func New(apiusername, apisecret string) *NotifoApiClient {
	return &NotifoApiClient{
		endpoint:    "https://api.notifo.com/v1/",
		apiusername: apiusername,
		apisecret:   apisecret,
	}
}

func (api *NotifoApiClient) SetEndpoint(endpoint string) {
	api.endpoint = endpoint
}

func (api *NotifoApiClient) SubscribeUser(username string) (*http.Response, error) {
	data := make(map[string]string)
	data["username"] = username

	return api.submitRequest("subscribe_user", "POST", data)
}

func (api *NotifoApiClient) SendNotification(to string, msg string, label string,
	title string, uri string) (*http.Response, error) {
	data := make(map[string]string)

	if to == "" {
		return nil, errors.New("'to' must not be blank")
	}
	data["to"] = to

	if msg == "" {
		return nil, errors.New("'msg' must not be blank")
	}
	data["msg"] = msg
	if label != "" {
		data["label"] = label
	}
	if title != "" {
		data["title"] = title
	}
	if uri != "" {
		data["uri"] = uri
	}

	return api.submitRequest("send_notification", "POST", data)
}

func (api *NotifoApiClient) submitRequest(action, method string, params map[string]string) (*http.Response, error) {
	vals := url.Values{}
	for k, v := range params {
		vals.Add(k, v)
	}

	client := http.Client{}

	req, err := http.NewRequest(method, strings.Join([]string{api.endpoint, action}, "/"),
		strings.NewReader(vals.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(api.apiusername, api.apisecret)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
