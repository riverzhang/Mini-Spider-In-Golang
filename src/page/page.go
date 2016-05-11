/* page.go - the implement of request's response */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
This file contains the struct 'Page', means request's response.
*/

package page

import (
	"fmt"
	"net/http"
	"net/url"
)

import (
	"request"
)

type Page struct {
	request          *request.Request   // request
	respStateCode    int                // response statue code
	respHeader       http.Header        // response header
	respCookies      []*http.Cookie     // response cookies
	respBody         []byte             // response body
	nextCrawlRequest []*request.Request // 继续抓取的请求
	isParse          bool               // 是否对页面进行 parse
}

func NewPage(req *request.Request) *Page {
	return &Page{request: req}
}

func (p *Page) Success() bool {
	if p.respStateCode >= http.StatusOK && p.respStateCode < http.StatusMultipleChoices {
		return true
	}
	return false
}

func (p *Page) GetNextCrawlRequest() []*request.Request {
	return p.nextCrawlRequest
}

func (p *Page) SetRespBody(body []byte) {
	p.respBody = body
}

func (p *Page) StrRespBody() string {
	return string(p.respBody)
}

func (p *Page) RespBody() []byte {
	return p.respBody
}

func (p *Page) SetRespStatusCode(status int) {
	p.respStateCode = status
}

func (p *Page) AddNextCrawlRequest(req *request.Request) {
	p.nextCrawlRequest = append(p.nextCrawlRequest, req)
}

func (p *Page) SetRespCookies(cookies []*http.Cookie) {
	p.respCookies = cookies
}

func (p *Page) SetRespHeader(header http.Header) {
	p.respHeader = header
}

func (p *Page) RespHeader(key string) string {
	return p.respHeader.Get(key)
}

func (p *Page) SetIsParse(isParse bool) {
	p.isParse = isParse
}

func (p *Page) IsParse() bool {
	return p.isParse
}

func (p *Page) URL() *url.URL {
	return p.request.URL
}

func (p *Page) Url() string {
	return p.request.Url()
}

func (p *Page) Depth() uint {
	return p.request.Depth()
}

func (p *Page) String() string {
	return fmt.Sprintf("Req: %s, Resp Code: %d, Need Parse: %t",
		p.request.String(),
		p.respStateCode,
		p.isParse)
}
