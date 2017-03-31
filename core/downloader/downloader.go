// Package downloader is the main module of crawler for download page.
package downloader

import (
	"github.com/viixv/crawler/core/commons/page"
	"github.com/viixv/crawler/core/commons/request"
)

// The Downloader interface.
// You can implement the interface by implement function Download.
// Function Download need to return Page instance pointer that has request result downloaded from Request.
type Downloader interface {
	Download(req *request.Request) *page.Page
}
