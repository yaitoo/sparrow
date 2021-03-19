package micro

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
)

//DefaultHTTPClient 默认HTTPClient实例
var DefaultHTTPClient = NewHTTPClient()

//HTTPClientOption HTTP请求参数
type HTTPClientOption func(c *HTTPClient)

//HTTPClient Http请求客户端, 可全局共享一个对象实例
type HTTPClient struct {
	// //URL  目标地址
	// URL string
	//Transport 网络参数
	transport *http.RoundTripper
	// Header    map[string]string
}

//NewHTTPClient 创建HTTPClient对象
func NewHTTPClient(options ...HTTPClientOption) *HTTPClient {
	c := &HTTPClient{}

	for _, opt := range options {
		opt(c)
	}

	if c.transport == nil {
		c.transport = &http.DefaultTransport
	}

	return c
}

//Do 执行请求
func (c *HTTPClient) Do(ctx context.Context, method, url string, data string, header map[string]string) ([]byte, int, error) {

	//	cookieJar, _ := cookiejar.New(nil)

	client := http.Client{
		Transport: *c.transport,
		//		Jar:       cookieJar,
	}

	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer func() { req.Close = true }()

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return b, resp.StatusCode, nil
	}

	return nil, resp.StatusCode, Throwf(ctx, ErrUnknown, "%s", string(b))

}
