package httpcheck

import (
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type HttpCheckClient struct {
}

type Response struct {
	statusCode int
	data       io.ReadCloser
}

func (hcc *HttpCheckClient) Check(proxy interface{}) (*Response, error) {
	var client http.Client

	switch pr := proxy.(type) {
	case func(string, string) (net.Conn, error):
		client = http.Client{
			Transport: &http.Transport{
				Dial: pr,
			},
		}

	case func(*http.Request) (*url.URL, error):
		client = http.Client{
			Transport: &http.Transport{
				Proxy: pr,
			},
		}

	default:
		logrus.Error("Invalid proxy %v", proxy)
	}

	resp, err := client.Get("http://api.ipify.org")
	if err != nil {
		logrus.Errorf("Proxy invalid - %v", err)
		return nil, err
	}

	return &Response{
		statusCode: resp.StatusCode,
		data:       resp.Body,
	}, nil
}

func (r *Response) IsSuccess() bool {
	return r.statusCode == 200
}

func (r *Response) GetStatusCodeRaw() int {
	return r.statusCode
}
