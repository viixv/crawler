// crawler master module
package crawler

import (
	"math/rand"
	"time"

	"github.com/viixv/crawler/core/commons/controller"
	"github.com/viixv/crawler/core/commons/log"
	"github.com/viixv/crawler/core/commons/page"
	"github.com/viixv/crawler/core/commons/request"
	"github.com/viixv/crawler/core/commons/result"
	"github.com/viixv/crawler/core/downloader"
	"github.com/viixv/crawler/core/pipeline"
	"github.com/viixv/crawler/core/processor"
	"github.com/viixv/crawler/core/scheduler"
)

type Crawler struct {
	taskname string

	pPageProcessor processor.PageProcessor

	pDownloader downloader.Downloader

	pScheduler scheduler.Scheduler

	pPiplelines []pipeline.Pipeline

	mc controller.GoroutineController

	threadnum uint

	exitWhenComplete bool

	// Sleeptype can be fixed or rand.
	startSleeptime uint
	endSleeptime   uint
	sleeptype      string
}

// Spider is scheduler module for all the other modules, like downloader, pipeline, scheduler and etc.
// The taskname could be empty string too, or it can be used in Pipeline for record the result crawled by which task;
func NewCrawler(pageinst processor.PageProcessor, taskname string) *Crawler {
	log.StraceInst().Open()

	ap := &Crawler{taskname: taskname, pPageProcessor: pageinst}

	// init filelog.
	ap.CloseFileLog()
	ap.exitWhenComplete = true
	ap.sleeptype = "fixed"
	ap.startSleeptime = 0

	// init spider
	if ap.pScheduler == nil {
		ap.SetScheduler(scheduler.NewQueueScheduler(false))
	}

	if ap.pDownloader == nil {
		ap.SetDownloader(downloader.NewHttpDownloader())
	}

	log.StraceInst().Println("** start crawler **")
	ap.pPiplelines = make([]pipeline.Pipeline, 0)

	return ap
}

func (this *Crawler) TaskName() string {
	return this.taskname
}

// Deal with one url and return the PageItems.
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
	// push url
	for _, req := range reqs {
		//req := request.NewRequest(u, respType, urltag, method, postdata, header, cookies)
		this.AddRequest(req)
	}

	pip := pipeline.NewPipelineCollector()
	this.AddPipeline(pip)

	this.Run()

	return pip.GetCollected()
}

func (this *Crawler) Run() {
	if this.threadnum == 0 {
		this.threadnum = 1
	}
	this.mc = controller.NewGoroutineControllerChan(this.threadnum)

	//init db  by sorawa

	for {
		req := this.pScheduler.Poll()

		// mc is not atomic
		if this.mc.Has() == 0 && req == nil && this.exitWhenComplete {
			log.StraceInst().Println("** executed callback **")
			this.pPageProcessor.Finish()
			log.StraceInst().Println("** end crawler **")
			break
		} else if req == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		this.mc.GetOne()

		// Asynchronous fetching
		go func(req *request.Request) {
			defer this.mc.FreeOne()
			//time.Sleep( time.Duration(rand.Intn(5)) * time.Second)
			log.StraceInst().Println("start crawl : " + req.GetUrl())
			this.pageProcess(req)
		}(req)
	}
	this.close()
}

func (this *Crawler) close() {
	this.SetScheduler(scheduler.NewQueueScheduler(false))
	this.SetDownloader(downloader.NewHttpDownloader())
	this.pPiplelines = make([]pipeline.Pipeline, 0)
	this.exitWhenComplete = true
}

func (this *Crawler) AddPipeline(p pipeline.Pipeline) *Crawler {
	this.pPiplelines = append(this.pPiplelines, p)
	return this
}

func (this *Crawler) SetScheduler(s scheduler.Scheduler) *Crawler {
	this.pScheduler = s
	return this
}

func (this *Crawler) GetScheduler() scheduler.Scheduler {
	return this.pScheduler
}

func (this *Crawler) SetDownloader(d downloader.Downloader) *Crawler {
	this.pDownloader = d
	return this
}

func (this *Crawler) GetDownloader() downloader.Downloader {
	return this.pDownloader
}

func (this *Crawler) SetThreadnum(i uint) *Crawler {
	this.threadnum = i
	return this
}

func (this *Crawler) GetThreadnum() uint {
	return this.threadnum
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

// The OpenFileLog initialize the log path and open log.
// If log is opened, error info or other useful info in spider will be logged in file of the filepath.
// Log command is mlog.LogInst().LogError("info") or mlog.LogInst().LogInfo("info").
// Spider's default log is closed.
// The filepath is absolute path.
func (this *Crawler) OpenFileLog(filePath string) *Crawler {
	log.InitFilelog(true, filePath)
	return this
}

// OpenFileLogDefault open file log with default file path like "WD/log/log.2014-9-1".
func (this *Crawler) OpenFileLogDefault() *Crawler {
	log.InitFilelog(true, "")
	return this
}

// The CloseFileLog close file log.
func (this *Crawler) CloseFileLog() *Crawler {
	log.InitFilelog(false, "")
	return this
}

// The OpenStrace open strace that output progress info on the screen.
// Spider's default strace is opened.
func (this *Crawler) OpenStrace() *Crawler {
	log.StraceInst().Open()
	return this
}

// The CloseStrace close strace.
func (this *Crawler) CloseStrace() *Crawler {
	log.StraceInst().Close()
	return this
}

// The SetSleepTime set sleep time after each crawl task.
// The unit is millisecond.
// If sleeptype is "fixed", the s is the sleep time and e is useless.
// If sleeptype is "rand", the sleep time is rand between s and e.
func (this *Crawler) SetSleepTime(sleeptype string, s uint, e uint) *Crawler {
	this.sleeptype = sleeptype
	this.startSleeptime = s
	this.endSleeptime = e
	if this.sleeptype == "rand" && this.startSleeptime >= this.endSleeptime {
		panic("startSleeptime must smaller than endSleeptime")
	}
	return this
}

func (this *Crawler) sleep() {
	if this.sleeptype == "fixed" {
		time.Sleep(time.Duration(this.startSleeptime) * time.Millisecond)
	} else if this.sleeptype == "rand" {
		sleeptime := rand.Intn(int(this.endSleeptime-this.startSleeptime)) + int(this.startSleeptime)
		time.Sleep(time.Duration(sleeptime) * time.Millisecond)
	}
}

func (this *Crawler) AddUrl(url string, respType string) *Crawler {
	req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	this.AddRequest(req)
	return this
}

func (this *Crawler) AddUrlEx(url string, respType string, headerFile string, proxyHost string) *Crawler {
	req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
	return this
}

func (this *Crawler) AddUrlWithHeaderFile(url string, respType string, headerFile string) *Crawler {
	req := request.NewRequestWithHeaderFile(url, respType, headerFile)
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

func (this *Crawler) AddUrlsWithHeaderFile(urls []string, respType string, headerFile string) *Crawler {
	for _, url := range urls {
		req := request.NewRequestWithHeaderFile(url, respType, headerFile)
		this.AddRequest(req)
	}
	return this
}

func (this *Crawler) AddUrlsEx(urls []string, respType string, headerFile string, proxyHost string) *Crawler {
	for _, url := range urls {
		req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
		this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
	}
	return this
}

// add Request to Schedule
func (this *Crawler) AddRequest(req *request.Request) *Crawler {
	if req == nil {
		log.LogInst().LogError("request is nil")
		return this
	} else if req.GetUrl() == "" {
		log.LogInst().LogError("request is empty")
		return this
	}
	this.pScheduler.Push(req)
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
		if err := recover(); err != nil { // do not affect other
			if strerr, ok := err.(string); ok {
				log.LogInst().LogError(strerr)
			} else {
				log.LogInst().LogError("pageProcess error")
			}
		}
	}()

	// download page
	for i := 0; i < 3; i++ {
		this.sleep()
		p = this.pDownloader.Download(req)
		if p.IsSucc() { // if fail retry 3 times
			break
		}

	}

	if !p.IsSucc() { // if fail do not need process
		return
	}

	this.pPageProcessor.Process(p)
	for _, req := range p.GetTargetRequests() {
		this.AddRequest(req)
	}

	// output
	if !p.GetSkip() {
		for _, pip := range this.pPiplelines {
			pip.Process(p.GetPageItems(), this)
		}
	}
}
