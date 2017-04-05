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
	hdUrl      *regexp.Regexp
	sdUrl      *regexp.Regexp
	ldUrl      *regexp.Regexp
}

func NewPageProcesser() *PageProcesser {
	p := PageProcesser{}
	p.popularReg = regexp.MustCompile("http://www\\.pearvideo\\.com/popular")
	p.videoReg = regexp.MustCompile("http://www\\.pearvideo\\.com/video_\\d+")
	p.hdUrl = regexp.MustCompile("hdUrl=\"(.*?)\"")
	p.sdUrl = regexp.MustCompile("sdUrl=\"(.*?)\"")
	p.ldUrl = regexp.MustCompile("ldUrl=\"(.*?)\"")
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
		scriptText := query.Find(".details-main.vertical-details.cmmain script").Text()
		if hdUrls := this.hdUrl.FindStringSubmatch(scriptText); len(hdUrls) > 1 {
			p.AddField("hdUrl", hdUrls[1])
		}
		if sdUrls := this.sdUrl.FindStringSubmatch(scriptText); len(sdUrls) > 1 {
			p.AddField("sdUrl", sdUrls[1])
		}
		if ldUrls := this.ldUrl.FindStringSubmatch(scriptText); len(ldUrls) > 1 {
			p.AddField("ldUrl", ldUrls[1])
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
