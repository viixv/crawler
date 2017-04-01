package scheduler

import (
	"container/list"
	"crypto/md5"
	"sync"

	"github.com/viixv/crawler/core/commons/request"
)

type QueueScheduler struct {
	mutex sync.Mutex
	rm    bool
	rmKey map[[md5.Size]byte]*list.Element
	queue *list.List
}

func NewQueueScheduler(rmDuplicate bool) *QueueScheduler {
	queue := list.New()
	rmKey := make(map[[md5.Size]byte]*list.Element)
	return &QueueScheduler{rm: rmDuplicate, queue: queue, rmKey: rmKey}
}

func (this *QueueScheduler) Push(req *request.Request) {
	this.mutex.Lock()
	var key [md5.Size]byte
	if this.rm {
		key = md5.Sum([]byte(req.GetUrl()))
		if _, ok := this.rmKey[key]; ok {
			this.mutex.Unlock()
			return
		}
	}
	e := this.queue.PushBack(req)
	if this.rm {
		this.rmKey[key] = e
	}
	this.mutex.Unlock()
}

func (this *QueueScheduler) Poll() *request.Request {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.queue.Len() <= 0 {
		return nil
	}
	e := this.queue.Front()
	req := e.Value.(*request.Request)
	key := md5.Sum([]byte(req.GetUrl()))
	this.queue.Remove(e)
	if this.rm {
		delete(this.rmKey, key)
	}
	return req
}

func (this *QueueScheduler) Count() int {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.queue.Len()
}
