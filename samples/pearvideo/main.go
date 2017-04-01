package main

import (
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/viixv/crawler/core/commons/page"
	"github.com/viixv/crawler/core/crawler"
	"github.com/viixv/crawler/core/pipeline"
)

type PageProcesser struct {
	popularReg *regexp.Regexp
	videoReg   *regexp.Regexp
}

func NewPageProcesser() *PageProcesser {
	p := PageProcesser{}
	p.popularReg, _ = regexp.Compile("http://www\\.pearvideo\\.com/popular")
	p.videoReg, _ = regexp.Compile("http://www\\.pearvideo\\.com/video_\\d+")
	return &p
}

func (this *PageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}

	query := p.GetHtmlParser()

	if this.popularReg.MatchString(p.GetRequest().Url) {
		var urls []string
		query.Find(".popularem.clearfix").Each(func(i int, s *goquery.Selection) {
			if href, e := s.Find(".popularembd.actplay").Attr("href"); e {
				urls = append(urls, "http://www.pearvideo.com/"+href)
			}
		})
		p.AddTargetRequests(urls, "html")
		return
	}

	if this.videoReg.MatchString(p.GetRequest().Url) {
		if title, e := query.Find("#share-to").Attr("data-title"); e {
			p.AddField("title", title)
		}

		if summary, e := query.Find("#share-to").Attr("data-summary"); e {
			p.AddField("summary", summary)
		}

		if picurl, e := query.Find("#share-to").Attr("data-picurl"); e {
			p.AddField("picurl", picurl)
		}
	}
}

func (this *PageProcesser) Finish() {
}

func main() {
	crawler.NewCrawler(NewPageProcesser(), "梨视频").
		AddUrl("http://www.pearvideo.com/popular", "html").
		AddPipeline(pipeline.NewConsolePipeline()).
		SetThreadnum(64).
		Run()
}
