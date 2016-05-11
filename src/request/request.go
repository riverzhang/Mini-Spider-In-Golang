/* request.go - the crawl object of mini spider */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
request are crawl by spider. request used by http client.
*/

package request

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// 请求方法
const (
	GET  = "GET"
	POST = "POST"
	HEAD = "HEAD"
	PUT  = "PUT"
)

type Request struct {
	*http.Request      // 继承 http.Request
	retryTimes    int  // 重试次数
	isDownload    bool // 是否下载页面
	depth         uint // 深度
}

func NewRequest(method string, url string, isDownload bool, depth uint, data url.Values) *Request {
	if url == "" {
		return nil
	}
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil
	}
	return &Request{Request: req, retryTimes: 0, isDownload: isDownload, depth: depth}
}

func (r *Request) RetryTimes() int {
	return r.retryTimes
}

func (r *Request) IsDownload() bool {
	return r.isDownload
}

func (r *Request) SetIsDownload(isdownload bool) {
	r.isDownload = isdownload
}

func (r *Request) Depth() uint {
	return r.depth
}

func (r *Request) Url() string {
	return r.Request.URL.String()
}

func (r *Request) GetRequest() *http.Request {
	return r.Request
}

func (r *Request) Valid() bool {
	if r.Url() == "" {
		return false
	}
	return true
}

func (r *Request) String() string {
	return fmt.Sprintf(
		"Request Method: %s, Url: %s, Retry Times: %d, Depth: %d, Download: %t",
		r.Request.Method,
		r.Request.URL.String(),
		r.retryTimes,
		r.depth,
		r.isDownload)
}
