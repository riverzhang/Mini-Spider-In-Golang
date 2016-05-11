package downloader

import (
	"strings"
	"testing"
)

import (
	"request"
	"util"
)

func Test_HTTPDownload(t *testing.T) {

	url := "https://www.baidu.com"
	req := request.NewRequest(request.GET, url, false, util.SEED_START_DEPTH, nil)

	var timeout uint
	timeout = 2
	outputDir := ""
	d := NewHttpDownloader(timeout, outputDir)
	p, err := d.Download(req)
	if err != nil {
		t.Error("http download failed!")
	}

	strBody := p.StrRespBody()
	if !strings.Contains(strBody, "baidu") {
		t.Error("download html page failed!")
	}
}
