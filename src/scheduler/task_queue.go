/* scheduler_queue.go - the implement of mini spider's scheduler */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
use container/list implement scheduler.
*/

package scheduler

import (
	"container/list"
	"crypto/md5"
	"errors"
	"sync"
)

import (
	"request"
)

type TaskQueue struct {
	queueLocker *sync.Mutex            // 队列锁
	taskQueue   *list.List             // 爬虫抓取队列
	hashTable   map[[md5.Size]byte]int // 已抓取请求的 hashtable，防止重复抓取
}

func NewTaskQueue() *TaskQueue {
	locker := new(sync.Mutex)
	queue := list.New()
	hash := make(map[[md5.Size]byte]int)
	return &TaskQueue{queueLocker: locker, taskQueue: queue, hashTable: hash}
}

func (t *TaskQueue) Push(req *request.Request) {
	t.queueLocker.Lock()
	defer t.queueLocker.Unlock()

	url := []byte(req.Request.URL.String())
	if len(url) <= 0 {
		return
	}
	key := md5.Sum(url)
	if _, ok := t.hashTable[key]; ok {
		return
	}

	t.taskQueue.PushBack(req)
	t.hashTable[key] = 1
}

func (t *TaskQueue) Pop() (*request.Request, error) {
	t.queueLocker.Lock()
	defer t.queueLocker.Unlock()

	if t.taskQueue.Len() <= 0 {
		return nil, nil
	}
	e := t.taskQueue.Front()
	t.taskQueue.Remove(e)

	req, ok := e.Value.(*request.Request)
	if ok {
		return req, nil
	}
	return nil, errors.New("pop request not the type of *request.Request")
}

func (t *TaskQueue) Count() int {
	t.queueLocker.Lock()
	defer t.queueLocker.Unlock()
	return t.taskQueue.Len()
}

func (t *TaskQueue) Empty() bool {
	t.queueLocker.Lock()
	defer t.queueLocker.Unlock()
	return t.taskQueue.Len() == 0
}
