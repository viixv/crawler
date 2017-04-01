package pipeline

import (
	"fmt"
	"sync"

	"github.com/viixv/crawler/core/commons/result"
	"github.com/viixv/crawler/core/commons/task"
)

type ConsolePipeline struct {
	mutex sync.Mutex
}

func NewConsolePipeline() *ConsolePipeline {
	return &ConsolePipeline{}
}

func (p *ConsolePipeline) Process(items *result.ResultItems, t task.Task) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	fmt.Println("***************************************************************")
	fmt.Println("Crawled url:\t" + items.GetRequest().GetUrl() + "\n")
	fmt.Println("Crawled result:")
	for key, value := range items.GetAll() {
		fmt.Println(key + "\t:\t" + value)
	}
	fmt.Println()
}
