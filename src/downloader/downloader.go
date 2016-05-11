/* downloader.go - the interface of mini spider's downloader */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
This file contains the interface 'Downloader' and list method of the interface.
*/

package downloader

import (
	"page"
	"request"
)

type Downloader interface {
	Download(req *request.Request) (*page.Page, error)
}
