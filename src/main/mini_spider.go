/* mini_spider.go - the main of spider */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
This is the main of spider.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
)

import (
	"config"
	"downloader"
	"parser"
	"spider"
	"util"
)

var (
	showVersion bool
	confFile    string
	logDir      string
)

func init() {
	flag.BoolVar(&showVersion, "v", false, "Show Version And Exit")
	flag.StringVar(&confFile, "c", "../conf/spider.conf", "Set Config Directory")
	flag.StringVar(&logDir, "l", "../log/", "Set Log Directory")
}

func main() {
	flag.Usage = util.PrintUsageExit
	flag.Parse()
	if showVersion {
		spider.PrintSpiderVersionExit()
	}

	// config
	confFile = path.Clean(confFile)
	cfg, err := config.GetConfigFromFile(confFile)
	if err != nil {
		fmt.Println("Parse Config Err: " + err.Error())
		os.Exit(1)
	}

	// log
	logDir = path.Clean(logDir)
	if !util.IsDirExist(logDir) {
		if ok, err := util.Mkdir(logDir); !ok {
			fmt.Println("Make Log Dir Err: " + err.Error())
			os.Exit(1)
		}
	}
	util.LogInit("mini_spider", "INFO", logDir, true, "midnight", 5)

	// new spider
	runtime.GOMAXPROCS(runtime.NumCPU())
	spider := spider.NewSpider(cfg.Spider.ThreadCount, cfg.Spider.MaxDepth, cfg.Spider.CrawlInterval)

	// add spider start urls
	var startUrls map[string]bool
	startUrls, err = util.GetSeedFromFile(cfg.Spider.UrlListFile, false)
	if err != nil {
		util.Logger.Error(err)
		os.Exit(1)
	}
	if len(startUrls) <= 0 {
		util.Logger.Warn("no url in file: " + cfg.Spider.UrlListFile)
		os.Exit(1)
	}
	spider.AddStartUrls(startUrls)

	// init spider's parser
	p := parser.NewRegexpParser(cfg.Spider.TargetUrl)
	spider.AddParser(p)

	// init spider's downloader
	d := downloader.NewHttpDownloader(cfg.Spider.CrawlTimeout, cfg.Spider.OutputDirectory)
	spider.AddDownloader(d)

	// begin crawl
	spider.Start()
}
