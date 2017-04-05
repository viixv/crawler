// Package page 保存Downloader下载结果，并且用于Processer解析页面。
package page

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/viixv/crawler/core/commons/request"
	"github.com/viixv/crawler/core/commons/result"
)

// Page代表了从Downloader下载到的一个页面，可能是HTML，也可能是JSON或者其他文本格式的内容。
type Page struct {
	// isFail 表示下载页面是否失败。
	isFail bool
	// 错误信息
	errorMsg string

	// Request 相关信息。
	req *request.Request

	// 结果文本
	body string

	header  http.Header
	cookies []*http.Cookie

	// The docParser is a pointer of goquery boject that contains html result.
	docParser *goquery.Document

	// The jsonMap is the json result.
	jsonMap *simplejson.Json

	// The pItems is object for save Key-Values in PageProcesser.
	// And pItems is output in Pipline.
	pItems *result.ResultItems

	// The targetRequests is requests to put into Scheduler.
	targetRequests []*request.Request
}

// NewPage returns initialized Page object.
func NewPage(req *request.Request) *Page {
	return &Page{pItems: result.NewResultItems(req), req: req}
}

// SetHeader save the header of http responce
func (this *Page) SetHeader(header http.Header) {
	this.header = header
}

// GetHeader returns the header of http responce
func (this *Page) GetHeader() http.Header {
	return this.header
}

// SetHeader save the cookies of http responce
func (this *Page) SetCookies(cookies []*http.Cookie) {
	this.cookies = cookies
}

// GetHeader returns the cookies of http responce
func (this *Page) GetCookies() []*http.Cookie {
	return this.cookies
}

// IsSucc test whether download process success or not.
func (this *Page) IsSucc() bool {
	return !this.isFail
}

// Errormsg show the download error message.
func (this *Page) Errormsg() string {
	return this.errorMsg
}

// SetStatus save status info about download process.
func (this *Page) SetStatus(isFail bool, errorMsg string) {
	this.isFail = isFail
	this.errorMsg = errorMsg
}

// AddField saves KV string pair to PageItems preparing for Pipeline
func (this *Page) AddField(key string, value string) {
	this.pItems.AddItem(key, value)
}

// GetPageItems returns PageItems object that record KV pair parsed in PageProcesser.
func (this *Page) GetPageItems() *result.ResultItems {
	return this.pItems
}

// SetSkip set label "skip" of PageItems.
// PageItems will not be saved in Pipeline wher skip is set true
func (this *Page) SetSkip(skip bool) {
	this.pItems.SetSkip(skip)
}

// GetSkip returns skip label of PageItems.
func (this *Page) GetSkip() bool {
	return this.pItems.GetSkip()
}

// SetRequest saves request oject of this page.
func (this *Page) SetRequest(r *request.Request) *Page {
	this.req = r
	return this
}

// GetRequest returns request oject of this page.
func (this *Page) GetRequest() *request.Request {
	return this.req
}

// GetUrlTag returns name of url.
func (this *Page) GetUrlTag() string {
	return this.req.GetUrlTag()
}

// AddTargetRequest adds one new Request waitting for crawl.
func (this *Page) AddTargetRequest(url string, respType string) *Page {
	this.targetRequests = append(this.targetRequests, request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil))
	return this
}

// AddTargetRequests adds new Requests waitting for crawl.
func (this *Page) AddTargetRequests(urls []string, respType string) *Page {
	for _, url := range urls {
		this.AddTargetRequest(url, respType)
	}
	return this
}

// AddTargetRequestWithProxy adds one new Request waitting for crawl.
func (this *Page) AddTargetRequestWithProxy(url string, respType string, proxyHost string) *Page {

	this.targetRequests = append(this.targetRequests, request.NewRequestWithProxy(url, respType, "", "GET", "", nil, nil, proxyHost, nil, nil))
	return this
}

// AddTargetRequestsWithProxy adds new Requests waitting for crawl.
func (this *Page) AddTargetRequestsWithProxy(urls []string, respType string, proxyHost string) *Page {
	for _, url := range urls {
		this.AddTargetRequestWithProxy(url, respType, proxyHost)
	}
	return this
}

// AddTargetRequest adds one new Request waitting for crawl.
// The respType is "html" or "json" or "jsonp" or "text".
// The urltag is name for marking url and distinguish different urls in PageProcesser and Pipeline.
// The method is POST or GET.
// The postdata is http body string.
// The header is http header.
// The cookies is http cookies.
func (this *Page) AddTargetRequestWithParams(req *request.Request) *Page {
	this.targetRequests = append(this.targetRequests, req)
	return this
}

// AddTargetRequests adds new Requests waitting for crawl.
func (this *Page) AddTargetRequestsWithParams(reqs []*request.Request) *Page {
	for _, req := range reqs {
		this.AddTargetRequestWithParams(req)
	}
	return this
}

// GetTargetRequests returns the target requests that will put into Scheduler
func (this *Page) GetTargetRequests() []*request.Request {
	return this.targetRequests
}

// SetBodyStr saves plain string crawled in Page.
func (this *Page) SetBodyStr(body string) *Page {
	this.body = body
	return this
}

// GetBodyStr returns plain string crawled.
func (this *Page) GetBodyStr() string {
	return this.body
}

// SetHtmlParser saves goquery object binded to target crawl result.
func (this *Page) SetHtmlParser(doc *goquery.Document) *Page {
	this.docParser = doc
	return this
}

// GetHtmlParser returns goquery object binded to target crawl result.
func (this *Page) GetHtmlParser() *goquery.Document {
	return this.docParser
}

// GetHtmlParser returns goquery object binded to target crawl result.
func (this *Page) ResetHtmlParser() *goquery.Document {
	r := strings.NewReader(this.body)
	var err error
	this.docParser, err = goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	return this.docParser
}

// SetJson saves json result.
func (this *Page) SetJson(js *simplejson.Json) *Page {
	this.jsonMap = js
	return this
}

// SetJson returns json result.
func (this *Page) GetJson() *simplejson.Json {
	return this.jsonMap
}
