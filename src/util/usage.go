package util

import (
	"fmt"
	"os"
)

const usageStr = `
Usage of mini_spider: ./mini_spider [options...]

Options:    
    -c="../conf/spider.conf" : Set  Config  Directory
    -l="../log/"             : Set  Log     Directory
    -v                       : Show Version And Exit

Example:
    ./mini_spider -c /path/to/mini_spider/conf/spider.conf -l /path/to/mini_spider/log/
`

func PrintUsageExit() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}
