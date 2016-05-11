package scheduler

import (
	"strconv"
	"testing"
)

import (
	"request"
	"util"
)

// go test
// go test -test.bench=".*"

// 初始化的队列 为空
func TestQueueEmpty(t *testing.T) {
	queue := NewTaskQueue()
	if queue.Count() != 0 {
		t.Error("new scheduler queue not empty")
	}
}

// 队列去重
func TestUniqQueue(t *testing.T) {
	queue := NewTaskQueue()

	url := "https://www.baidu.com"
	req := request.NewRequest(request.GET, url, false, util.SEED_START_DEPTH, nil)

	queue.Push(req)
	if queue.Count() != 1 {
		t.Error("scheduler queue push one, count != 1")
	}

	queue.Push(req)
	if queue.Count() != 1 {
		t.Error("scheduler queue push same item, count != 1")
	}

	queue.Pop()
	if queue.Count() != 0 {
		t.Error("scheduler queue should empty")
	}
}

// 队列非重 req 计数
func TestPushQueue(t *testing.T) {
	queue := NewTaskQueue()

	url := "https://www.baidu.com"
	req := request.NewRequest(request.GET, url, false, util.SEED_START_DEPTH, nil)
	queue.Push(req)
	if queue.Count() != 1 {
		t.Error("scheduler queue push one, count != 1")
	}

	url = "http://www.qq.com/"
	req = request.NewRequest(request.GET, url, false, util.SEED_START_DEPTH, nil)
	queue.Push(req)
	if queue.Count() != 2 {
		t.Error("scheduler queue push two, count != 2")
	}
}

// 压力测试
func Benchmark_Push(b *testing.B) {
	b.StopTimer()
	queue := NewTaskQueue()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		url := "https://www.baidu.com/" + strconv.Itoa(b.N)
		req := request.NewRequest(request.GET, url, false, util.SEED_START_DEPTH, nil)
		queue.Push(req)
	}
}
