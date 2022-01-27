package atsgolangdemo

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Request struct {
	http          *http.Client
	ApiURL        string
	OperatorToken string
	SecretKey     string
}

func NewRequest(ApiURL, OperatorToken, SecretKey string) Request {
	return Request{ApiURL: ApiURL, OperatorToken: OperatorToken, SecretKey: SecretKey}
}

type Response struct {
	Code        int    `json:"code"`
	Message     string `json:"msg"`
	RequestID   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
}

func (r *Request) Post(path string, params interface{}, timeout int) (Response, error) {
	var apiResp Response
	headers := map[string]string{
		"User-Agent":     "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; .NET CLR 1.1.4322)",
		"Content-Type":   "application/json",
		"operator-token": r.OperatorToken,
		"secret-key":     r.SecretKey,
	}

	var payload io.Reader
	if params != nil {
		if strings.HasPrefix(reflect.TypeOf(params).String(), "map") {
			bs, _ := json.Marshal(params)
			payload = bytes.NewReader(bs)
		} else if reflect.TypeOf(params).String() == "string" {
			payload = bytes.NewReader([]byte(params.(string)))
		}
	}

	creq, _ := http.NewRequest("POST", r.ApiURL+path, payload)
	for k, v := range headers {
		creq.Header.Set(k, v)
	}
	crep, err := (&http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout))
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * time.Duration(timeout),
		},
	}).Do(creq)
	if err != nil {
		return apiResp, err
	}
	defer crep.Body.Close()
	if crep.StatusCode == 200 {
		rawBody, err := ioutil.ReadAll(crep.Body)
		if err != nil {
			return apiResp, err
		}
		if err := json.Unmarshal(rawBody, &apiResp); err != nil {
			return apiResp, err
		}
		return apiResp, nil
	}
	return apiResp, err
}
