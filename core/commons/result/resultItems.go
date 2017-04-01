package result

import (
	"github.com/viixv/crawler/core/commons/request"
)

type ResultItems struct {
	req   *request.Request
	items map[string]string
	skip  bool
}

func NewResultItems(req *request.Request) *ResultItems {
	items := make(map[string]string)
	return &ResultItems{req: req, items: items, skip: false}
}

func (res *ResultItems) GetRequest() *request.Request {
	return res.req
}

func (res *ResultItems) AddItem(key string, item string) {
	res.items[key] = item
}

func (res *ResultItems) GetItem(key string) (string, bool) {
	t, ok := res.items[key]
	return t, ok
}

func (res *ResultItems) GetAll() map[string]string {
	return res.items
}

func (res *ResultItems) SetSkip(skip bool) *ResultItems {
	res.skip = skip
	return res
}

func (res *ResultItems) GetSkip() bool {
	return res.skip
}
