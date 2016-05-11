/* parser.go - the interface of mini spider's parser */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
define the interface of parser(parser response body)
*/

package parser

import (
	"page"
)

type Parser interface {
	Parse(page *page.Page) error // 解析页面内容
}
