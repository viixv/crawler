package crawler

import (
	"log"
	"math/rand"
	"time"

	"github.com/viixv/crawler/core/commons/controller"
	"github.com/viixv/crawler/core/commons/page"
	"github.com/viixv/crawler/core/commons/request"
	"github.com/viixv/crawler/core/commons/result"
	"github.com/viixv/crawler/core/downloader"
	"github.com/viixv/crawler/core/pipeline"
	"github.com/viixv/crawler/core/processor"
	"github.com/viixv/crawler/core/scheduler"
)

type Crawler struct {
	cController      controller.GoroutineController
	cDownloader      downloader.Downloader
	cScheduler       scheduler.Scheduler
	exitWhenComplete bool
	goroutines       uint
	pageProcessor    processor.PageProcessor
	pipelines        []pipeline.Pipeline
	sleepType        string
	startSleepTime   uint
	endSleepTime     uint
	taskName         string
}

func NewCrawler(pageProcessor processor.PageProcessor, taskName string) *Crawler {
	crawler := Crawler{taskName: taskName, pageProcessor: pageProcessor}
	crawler.exitWhenComplete = true
	crawler.sleepType = "fixed"
	crawler.startSleepTime = 0
	if crawler.cScheduler == nil {
		crawler.SetScheduler(scheduler.NewQueueScheduler(false))
	}
	if crawler.cDownloader == nil {
		crawler.SetDownloader(downloader.NewHttpDownloader())
	}
	crawler.pipelines = make([]pipeline.Pipeline, 0)
	log.Println("Crawler initialization complete.")
	return &crawler
}

func (this *Crawler) TaskName() string {
	return this.taskName
}

func (this *Crawler) Get(url string, respType string) *result.ResultItems {
	req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	return this.GetByRequest(req)
}

// Deal with several urls and return the PageItems slice.
func (this *Crawler) GetAll(urls []string, respType string) []*result.ResultItems {
	for _, u := range urls {
		req := request.NewRequest(u, respType, "", "GET", "", nil, nil, nil, nil)
		this.AddRequest(req)
	}

	pip := pipeline.NewPipelineCollector()
	this.AddPipeline(pip)

	this.Run()

	return pip.GetCollected()
}

// Deal with one url and return the PageItems with other setting.
func (this *Crawler) GetByRequest(req *request.Request) *result.ResultItems {
	var reqs []*request.Request
	reqs = append(reqs, req)
	items := this.GetAllByRequest(reqs)
	if len(items) != 0 {
		return items[0]
	}
	return nil
}

// Deal with several urls and return the PageItems slice
func (this *Crawler) GetAllByRequest(reqs []*request.Request) []*result.ResultItems {
	for _, req := range reqs {
		this.AddRequest(req)
	}

	pip := pipeline.NewPipelineCollector()
	this.AddPipeline(pip)

	this.Run()

	return pip.GetCollected()
}

func (this *Crawler) Run() {
	if this.goroutines == 0 {
		this.goroutines = 1
	}
	this.cController = controller.NewGoroutineControllerChan(this.goroutines)

	for {
		req := this.cScheduler.Poll()
		if this.cController.Has() == 0 && req == nil && this.exitWhenComplete {
			this.pageProcessor.Finish()
			log.Println("Crawling complete.")
			break
		} else if req == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		this.cController.GetOne()
		go func(req *request.Request) {
			defer this.cController.FreeOne()
			log.Println("start crawl : " + req.GetUrl())
			this.pageProcess(req)
		}(req)
	}
	this.close()
}

func (this *Crawler) close() {
	this.SetScheduler(scheduler.NewQueueScheduler(false))
	this.SetDownloader(downloader.NewHttpDownloader())
	this.pipelines = make([]pipeline.Pipeline, 0)
	this.exitWhenComplete = true
}

func (this *Crawler) AddPipeline(p pipeline.Pipeline) *Crawler {
	this.pipelines = append(this.pipelines, p)
	return this
}

func (this *Crawler) SetScheduler(s scheduler.Scheduler) *Crawler {
	this.cScheduler = s
	return this
}

func (this *Crawler) GetScheduler() scheduler.Scheduler {
	return this.cScheduler
}

func (this *Crawler) SetDownloader(d downloader.Downloader) *Crawler {
	this.cDownloader = d
	return this
}

func (this *Crawler) GetDownloader() downloader.Downloader {
	return this.cDownloader
}

func (this *Crawler) SetThreadnum(i uint) *Crawler {
	this.goroutines = i
	return this
}

func (this *Crawler) GetThreadnum() uint {
	return this.goroutines
}

// If exit when each crawl task is done.
// If you want to keep spider in memory all the time and add url from outside, you can set it true.
func (this *Crawler) SetExitWhenComplete(e bool) *Crawler {
	this.exitWhenComplete = e
	return this
}

func (this *Crawler) GetExitWhenComplete() bool {
	return this.exitWhenComplete
}

func (this *Crawler) SetSleepTime(sleeptype string, s uint, e uint) *Crawler {
	this.sleepType = sleeptype
	this.startSleepTime = s
	this.endSleepTime = e
	if this.sleepType == "rand" && this.startSleepTime >= this.endSleepTime {
		panic("startSleeptime must smaller than endSleeptime")
	}
	return this
}

func (this *Crawler) sleep() {
	if this.sleepType == "fixed" {
		time.Sleep(time.Duration(this.startSleepTime) * time.Millisecond)
	} else if this.sleepType == "rand" {
		sleeptime := rand.Intn(int(this.endSleepTime-this.startSleepTime)) + int(this.startSleepTime)
		time.Sleep(time.Duration(sleeptime) * time.Millisecond)
	}
}

func (this *Crawler) AddUrl(url string, respType string) *Crawler {
	req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	this.AddRequest(req)
	return this
}

func (this *Crawler) AddUrls(urls []string, respType string) *Crawler {
	for _, url := range urls {
		req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
		this.AddRequest(req)
	}
	return this
}

// add Request to Schedule
func (this *Crawler) AddRequest(req *request.Request) *Crawler {
	if req == nil {
		log.Println("request is nil")
		return this
	} else if req.GetUrl() == "" {
		log.Println("request is empty")
		return this
	}
	this.cScheduler.Push(req)
	return this
}

//
func (this *Crawler) AddRequests(reqs []*request.Request) *Crawler {
	for _, req := range reqs {
		this.AddRequest(req)
	}
	return this
}

// core processer
func (this *Crawler) pageProcess(req *request.Request) {
	var p *page.Page
	defer func() {
		if err := recover(); err != nil {
			if strerr, ok := err.(string); ok {
				log.Println(strerr)
			} else {
				log.Println("pageProcess error")
			}
		}
	}()

	for i := 0; i < 3; i++ {
		this.sleep()
		p = this.cDownloader.Download(req)
		if p.IsSucc() {
			break
		}
	}

	if !p.IsSucc() {
		return
	}

	this.pageProcessor.Process(p)
	for _, req := range p.GetTargetRequests() {
		this.AddRequest(req)
	}

	if !p.GetSkip() {
		for _, pipe := range this.pipelines {
			pipe.Process(p.GetPageItems(), this)
		}
	}
}
