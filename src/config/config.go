/* config.go - the config of mini spider */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
This file contains the struct 'Config' of mini spider and check method of the struct.
*/

package config

import (
	"errors"
	"path/filepath"
)

import (
	"code.google.com/p/gcfg"
)

import (
	"util"
)

type Config struct {
	Spider struct {
		UrlListFile     string // 种子文件路径
		OutputDirectory string // 下载目录
		MaxDepth        uint   // 最大抓取深度
		CrawlInterval   uint   // 抓取间隔
		CrawlTimeout    uint   // 抓取超时
		TargetUrl       string // 目标文件正则
		ThreadCount     uint   // 抓取routine数
	}
}

func (c *Config) Valid() (bool, error) {
	if !util.IsFileExist(c.Spider.UrlListFile) {
		return false, errors.New("UrlListFile file not exitst: " + c.Spider.UrlListFile)
	}
	if !util.IsDirExist(c.Spider.OutputDirectory) {
		return false, errors.New("OutputDirectory not exitst: " + c.Spider.OutputDirectory)
	}
	if c.Spider.CrawlInterval < 0 {
		return false, errors.New("CrawlInterval must >= 0")
	}
	if c.Spider.CrawlTimeout <= 0 {
		return false, errors.New("CrawlTimeout must > 0")
	}
	if c.Spider.TargetUrl == "" {
		return false, errors.New("TargetUrl is empty")
	}
	if c.Spider.ThreadCount <= 0 {
		return false, errors.New("ThreadCount must > 0")
	}
	return true, nil
}

/*
* GetConfigFromFile - get spider config from file
*
* PARAMS:
*   - file: config file path
*
* RETURNS:
*   *Config, nil, if succeed
*   nil, error,   if fail
 */
func GetConfigFromFile(file string) (*Config, error) {
	var config Config
	err := gcfg.ReadFileInto(&config, file)
	if err != nil {
		return nil, err
	}

	var configDir string
	configDir, err = filepath.Abs(filepath.Dir(file))
	if !util.IsFileExist(config.Spider.UrlListFile) {
		urlListFile := configDir + "/" + config.Spider.UrlListFile
		if util.IsFileExist(urlListFile) {
			config.Spider.UrlListFile = urlListFile
		}
	}
	if !util.IsDirExist(config.Spider.OutputDirectory) {
		outputDir := configDir + "/" + config.Spider.UrlListFile
		if util.IsDirExist(outputDir) {
			config.Spider.OutputDirectory = outputDir
		}
	}

	_, err = config.Valid()
	if err != nil {
		return nil, err
	}
	return &config, nil
}
