package request

import (
	"net/http"
)

type Request struct {
	Url           string
	RespType      string
	Method        string
	PostData      string
	UrlTag        string
	Header        http.Header
	Cookies       []*http.Cookie
	ProxyHost     string
	checkRedirect func(req *http.Request, via []*http.Request) error
	Meta          interface{}
}

func NewRequest(url string, respType string, urlTag string, method string,
	postdata string, header http.Header, cookies []*http.Cookie,
	checkRedirect func(req *http.Request, via []*http.Request) error,
	meta interface{}) *Request {
	return &Request{url, respType, method, postdata, urlTag, header, cookies, "", checkRedirect, meta}
}

func NewRequestWithProxy(url string, respType string, urltag string, method string,
	postdata string, header http.Header, cookies []*http.Cookie, proxyHost string,
	checkRedirect func(req *http.Request, via []*http.Request) error,
	meta interface{}) *Request {
	return &Request{url, respType, method, postdata, urltag, header, cookies, proxyHost, checkRedirect, meta}
}

func (this *Request) AddProxyHost(host string) *Request {
	this.ProxyHost = host
	return this
}

func (this *Request) GetUrl() string {
	return this.Url
}

func (this *Request) GetUrlTag() string {
	return this.UrlTag
}

func (this *Request) GetMethod() string {
	return this.Method
}

func (this *Request) GetPostdata() string {
	return this.PostData
}

func (this *Request) GetHeader() http.Header {
	return this.Header
}

func (this *Request) GetCookies() []*http.Cookie {
	return this.Cookies
}

func (this *Request) GetProxyHost() string {
	return this.ProxyHost
}

func (this *Request) GetResponceType() string {
	return this.RespType
}

func (this *Request) GetRedirectFunc() func(req *http.Request, via []*http.Request) error {
	return this.checkRedirect
}

func (this *Request) GetMeta() interface{} {
	return this.Meta
}
