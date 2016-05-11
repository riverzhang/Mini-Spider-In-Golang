/* regexp_parser.go - the implement of mini spider's parser */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
default parser of mini spider, you can implement your parser by demand
parser response body, get the next crawl request , and get target by regexp
*/

package parser

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

import (
	"page"
	"request"
	"util"
)

type RegexpParser struct {
	targetRegexp *regexp.Regexp
}

func NewRegexpParser(expr string) *RegexpParser {
	target := regexp.MustCompile(expr)
	return &RegexpParser{targetRegexp: target}
}

/*
* Parse - parse response
*
* PARAMS:
*   - page: the downloaded page
*
* RETURNS:
*   nil
 */
func (r *RegexpParser) Parse(page *page.Page) error {
	if page == nil || !page.IsParse() {
		return nil
	}
	strBody := page.StrRespBody()
	if strBody == "" {
		return nil
	}
	util.Logger.Info("##### PARSE REQ #####", page)

	// body converter to utf-8
	bodyReader := strings.NewReader(strBody)
	strUTF8Body, err := util.Charset2UTF8(page.RespHeader("Content-Type"), bodyReader)
	if err == nil {
		// use utf8 reader, or use original body reader
		bodyReader = strings.NewReader(strUTF8Body)
	}

	// parse body
	root, err := html.Parse(bodyReader)
	if err != nil {
		return err
	}
	r.ParseByNode(page, root)
	return nil
}

func (r *RegexpParser) ParseByNode(page *page.Page, n *html.Node) {

	// 找出超链接，作为下次抓取的对象
	// 超链接满足 target 条件，设置请求属性为下载
	if n.Type == html.ElementNode {
		for _, arr := range n.Attr {
			if arr.Key == "href" || arr.Key == "src" {
				refUrl := arr.Val
				if strings.HasPrefix(arr.Val, "javascript") {
					refUrl = util.GetUrlFromJs(refUrl)
				}
				absUrl, err := util.RelativeUrl2Abs(page.URL(), refUrl)
				if err != nil {
					util.Logger.Warn("Relative Url 2 Abs Err: ", refUrl, err)
					continue
				}

				req := request.NewRequest(request.GET, absUrl, false, page.Depth()+1, nil)
				req.Header.Add("Referer", page.Url())
				if r.targetRegexp.MatchString(absUrl) {
					req.SetIsDownload(true)
				}
				page.AddNextCrawlRequest(req)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.ParseByNode(page, c)
	}
}
