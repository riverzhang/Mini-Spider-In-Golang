/* spider.go - the engine of mini spider */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
spider is the engine of mini spider, coordinate all the component.
*/

package spider

import (
	"fmt"
	"os"
	"strings"
	"time"
)

import (
	"downloader"
	"manager"
	"parser"
	"request"
	"scheduler"
	"util"
)

const (
	VERSION = "0.0.1" // spider version

)

type Spider struct {
	routineNum     uint                    // 爬虫数量
	startUrls      []string                // 抓取起始Urls
	crawlMaxDepth  uint                    // 最大抓取深度
	crawlInterval  time.Duration           // 抓取间隔
	routineManager *manager.RoutineManager // 爬虫管理者
	downloadQueue  scheduler.Scheduler     // 爬取任务队列
	parser         parser.Parser           // 页面解析器
	downloader     downloader.Downloader   // 页面下载器
}

func NewSpider(routineNum uint, maxDepth uint, interval uint) *Spider {
	crawlInterval := time.Duration(interval) * time.Second
	manager := manager.NewRoutineManager(routineNum)
	queue := scheduler.NewTaskQueue()
	return &Spider{
		routineNum:     routineNum,
		crawlMaxDepth:  maxDepth,
		crawlInterval:  crawlInterval,
		routineManager: manager,
		downloadQueue:  queue,
	}
}

func (s *Spider) AddStartUrls(urls map[string]bool) {
	for url, isDownload := range urls {
		req := request.NewRequest(request.GET, url, isDownload, util.SEED_START_DEPTH, nil)
		if req != nil {
			s.addStartRequest(req)
		}
	}
}

func (s *Spider) addStartRequest(req *request.Request) bool {
	if s.addRequest(req) {
		s.startUrls = append(s.startUrls, req.Url())
		return true
	}
	return false
}

func (s *Spider) addRequest(req *request.Request) bool {
	if req == nil || !req.Valid() {
		return false
	}
	s.downloadQueue.Push(req)
	return true
}

func (s *Spider) AddDownloader(downloader downloader.Downloader) {
	s.downloader = downloader
}

func (s *Spider) AddParser(parser parser.Parser) {
	s.parser = parser
}

func (s *Spider) checkParserDownloader() bool {
	if s.downloader == nil {
		util.Logger.Error("downloader = nil, init downloader first")
		return false
	}
	if s.parser == nil {
		util.Logger.Error("parser = nil, init parser first")
		return false
	}
	return true
}

func (s *Spider) Start() {
	canRun := s.checkParserDownloader()
	if !canRun {
		os.Exit(1)
	}
	util.Logger.Info(s)

	/**
	 * 主线程为 master（不会阻塞）, 新起 routine 为 worker.
	 * 当抓取队列中有 待抓取请求时，master 分配 worker 去处理(parse download ...)。
	 * master 控制 worker routine 的数量。worker 最大并发数为 routineNum。
	 * worker 处理好请求后，把下次待抓取的 req push 到队列中。 worker 一直在轮训抓取队列，当队列中没有待处理的请求时退出。
	 */
	for {
		// check spider is finish
		if s.isFinishCrawl() {
			break
		}

		// spider busy, but queue empty, wait for request
		if s.downloadQueue.Empty() {
			continue
		}

		// master control worker
		bolGet := s.routineManager.GetOne()
		if !bolGet {
			continue
		}

		go func() {
			defer s.routineManager.FreeOne()

			// worker process request until queue is empty
			for {
				req, err := s.downloadQueue.Pop()
				if err != nil {
					util.Logger.Error("pop request from queue err: ", err)
					continue
				}
				if req == nil {
					break
				}
				util.Logger.Info("##### POP REQ #####", req.Url())

				s.process(req)
				s.routineSleep()
			}

		}()
	}
}

func (s *Spider) process(req *request.Request) {
	// download page
	page, err := s.downloader.Download(req)
	if err != nil || page == nil || !page.Success() {
		util.Logger.Error("download request err: ", req, err)
		return
	}

	// parse page
	if page.Depth() >= s.crawlMaxDepth {
		page.SetIsParse(false) // 页面深度 >= 阀值，不进行 parse(sub request 的深度+1，超过阀值）
	}
	err = s.parser.Parse(page)
	if err != nil {
		util.Logger.Error("parse page err: ", page, err)
		return
	}

	// add next crawl request
	for _, req := range page.GetNextCrawlRequest() {
		if req.Depth() > s.crawlMaxDepth {
			continue // 新请求深度 超过 阀值，不 push 到 抓取队列
		}
		s.downloadQueue.Push(req)
	}
}

func (s *Spider) isFinishCrawl() bool {
	if s.downloadQueue.Count() == 0 && s.routineManager.Used() == 0 {
		util.Logger.Info("request scheduler is empty, spider are idle, finish crawl, exist.")
		time.Sleep(100 * time.Millisecond)
		return true
	}
	return false
}

func (s *Spider) routineSleep() {
	// 控制 单routine 抓取间隔
	time.Sleep(s.crawlInterval)
}

func (s *Spider) String() string {
	return fmt.Sprintf(
		"spider routine num: %d, max depth: %d, crawl interval: %02.2f seconds, start urls: %s.",
		s.routineNum,
		s.crawlMaxDepth,
		s.crawlInterval.Seconds(),
		strings.Join(s.startUrls, ","))
}

func PrintSpiderVersionExit() {
	fmt.Printf("mini spider version %s\n", VERSION)
	os.Exit(0)
}
