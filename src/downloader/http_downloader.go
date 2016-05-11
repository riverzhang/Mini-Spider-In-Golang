/* http_downloader.go - the implement of http downloader */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
This file contains the struct 'HttpDownloader', implement the downloader interface.
*/

package downloader

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

import (
	"page"
	"request"
	"util"
)

// content-type 前5位，作为类型标识
const (
	CONTENT_TYPE_PREFIX_LEN = 5
)

// content-type 前5位, 如果在 map 中, 则在 parse 阶段不解析（图片、视频等无需解析）
var contentTypeNotPrase = map[string]bool{
	"video": true,
	"audio": true,
	"image": true,
	"messa": true, // message
	"appli": true, // application
}

type HttpDownloader struct {
	client       *http.Client  // http client, reuse client
	crawlTimeout time.Duration // 抓取超时
	outputDir    string        // 下载目录
}

/*
* NewHttpDownloader - new http downloader
*
* PARAMS:
*   - timeout: http client timeout
*   - downloadDir: download file save dir
*
* RETURNS:
*   *HttpDownloader
 */
func NewHttpDownloader(timeout uint, downloadDir string) *HttpDownloader {
	crawlTimeout := time.Duration(timeout) * time.Second
	defaultClient := &http.Client{
		Timeout: crawlTimeout,
	}
	return &HttpDownloader{
		client:       defaultClient,
		crawlTimeout: crawlTimeout,
		outputDir:    downloadDir,
	}
}

func (h *HttpDownloader) SetClient(client *http.Client) {
	h.client = client
}

func (h *HttpDownloader) SetTimeout(timeout int) *HttpDownloader {
	h.crawlTimeout = time.Second * time.Duration(timeout)
	h.client.Timeout = h.crawlTimeout
	return h
}

func (h *HttpDownloader) SetOutputDir(dir string) *HttpDownloader {
	if dir != "" {
		h.outputDir = dir
	}
	return h
}

/*
* Download - use http implement the downloader interface
*
* PARAMS:
*   - req: represent a request
*
* RETURNS:
*    *page.Page, nil   if success
*    nil, error        if fail
 */
func (h *HttpDownloader) Download(req *request.Request) (*page.Page, error) {
	util.Logger.Info("##### DOWNLOAD REQ #####", req)

	resp, err := h.client.Do(req.Request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	p := page.NewPage(req)
	p.SetRespHeader(resp.Header)
	p.SetRespCookies(resp.Cookies())
	p.SetRespStatusCode(resp.StatusCode)
	p.SetRespBody(content)

	// content type 为 图片、语音、视频等，不解析 page 内容
	contentType := http.DetectContentType(content)
	miniType := strings.ToLower(contentType[:CONTENT_TYPE_PREFIX_LEN])
	if _, found := contentTypeNotPrase[miniType]; found {
		p.SetIsParse(false)
	} else {
		p.SetIsParse(true)
	}

	/*
	 * 是否下载, 将需下载的文件保存到本地
	 * 写入失败，记录 Log，但正常返回
	 */
	if req.IsDownload() && p.Success() {
		saveFile := url.QueryEscape(req.Url())
		out, err := os.Create(h.outputDir + "/" + saveFile)
		if err != nil {
			util.Logger.Warn("downloader create file error", saveFile, err)
			return p, nil
		}
		defer out.Close()
		_, err = out.Write(content)
		if err != nil {
			util.Logger.Warn("downloader write resp body error", err)
		}
	}

	return p, nil
}
