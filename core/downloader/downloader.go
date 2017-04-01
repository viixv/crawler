package downloader

import (
	"github.com/viixv/crawler/core/commons/page"
	"github.com/viixv/crawler/core/commons/request"
)

type Downloader interface {
	Download(*request.Request) *page.Page
}
