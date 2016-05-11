
# mini spider

## spider component

1. request    代表一个请求，即一个 Url
2. downloader 代表下载器，将 request 下载下来，用响应生成 page
3. page       代表一个页面，即 Url 下载下来的页面
4. parser     代表解析器，对下载好的 page，进行解析（抓取下一次爬取的种子 和 target）
5. manager    负责管理 爬虫
6. scheduler  作为一个抓取队列（内部有个简单去重，对抓取过的 url 做记录）
7. spider     爬虫engine

## 代码目录结构介绍

```
$ ls
README.md  bin        conf       data      log        output     pkg        src
```

1. bin    可执行程序目录
2. conf   爬虫配置文件目录 conf/spider.conf
3. data   种子文件目录 data/url.data
4. log    抓取日志目录
5. output 定向抓取 target 输出目录
6. src    源代码 



```
$ ls src
config     downloader main       manager    page       parser     request    scheduler  spider     util
```

1. config      配置的实现，读取配置
2. downloader  下载器，实现请求的下载
3. main        mini spider 主程序
4. manager     管理routine
5. page        下载后页面
6. parser      解析器，可以自定义解析器，满足个性化需求
7. request     请求的实现，表示一个请求
8. scheduler   抓取队列，使用list实现


## 环境设置

```
export PROJPATH=/home/work/local/mini_spider
export GOPATH=/home/work/libs/Go/golang-lib:$PROJPATH
export GOBIN=$PROJPATH/bin
export PATH=$PATH:$GOBIN
```

## 安装运行

```
$ cd project_dir/bin
$ go build ../src/main/mini_spider.go
$ ./mini_spider
```

## 输出结果

直接运行的话, 使用 data/seed.data 中设置的抓取起点为 http://news.baidu.com/.

同时, 使用默认的 regexp_parser.go 来解析页面, regexp_parser.go 根据 conf/spider.conf 中 targetUrl 指定的正则, 来决定下次抓取页面.

conf/spider.conf 中 targetUrl 设置为 .*.(htm|html)$, 所以只抓取以 htm 或 html 结尾的页面.

你可以实现自己的 parser, 控制页面的分析和下次抓取页面.

抓取结果存储在 output 中.
