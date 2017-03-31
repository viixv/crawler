package pipeline

import (
	"github.com/viixv/crawler/core/commons/result"
	"github.com/viixv/crawler/core/commons/task"
)

type ConsolePipeline struct {
}

func NewConsolePipeline() *ConsolePipeline {
	return &ConsolePipeline{}
}

func (this *ConsolePipeline) Process(items *result.ResultItems, t task.Task) {
	println("----------------------------------------------------------------------------------------------")
	println("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
	println("Crawled result : ")
	for key, value := range items.GetAll() {
		println(key + "\t:\t" + value)
	}
}
