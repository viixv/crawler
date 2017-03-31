// Package result contains parsed result by PageProcesser.
// The result is processed by Pipeline.
package result

import (
	"github.com/viixv/crawler/core/commons/request"
)

// PageItems represents an entity save result parsed by PageProcesser and will be output at last.
type ResultItems struct {

	// The req is Request object that contains the parsed result, which saved in PageItems.
	req *request.Request

	// The items is the container of parsed result.
	items map[string]string

	// The skip represents whether send ResultItems to scheduler or not.
	skip bool
}

// NewResultItems returns initialized PageItems object.
func NewResultItems(req *request.Request) *ResultItems {
	items := make(map[string]string)
	return &ResultItems{req: req, items: items, skip: false}
}

// GetRequest returns request of PageItems
func (this *ResultItems) GetRequest() *request.Request {
	return this.req
}

// AddItem saves a KV result into PageItems.
func (this *ResultItems) AddItem(key string, item string) {
	this.items[key] = item
}

// GetItem returns value of the key.
func (this *ResultItems) GetItem(key string) (string, bool) {
	t, ok := this.items[key]
	return t, ok
}

// GetAll returns all the KVs result.
func (this *ResultItems) GetAll() map[string]string {
	return this.items
}

// SetSkip set skip true to make this page not to be processed by Pipeline.
func (this *ResultItems) SetSkip(skip bool) *ResultItems {
	this.skip = skip
	return this
}

// GetSkip returns skip label.
func (this *ResultItems) GetSkip() bool {
	return this.skip
}
