package pipeline

import (
	"github.com/viixv/crawler/core/commons/result"
	"github.com/viixv/crawler/core/commons/task"
)

type PipelineCollector struct {
	collector []*result.ResultItems
}

func NewPipelineCollector() *PipelineCollector {
	collector := make([]*result.ResultItems, 0)
	return &PipelineCollector{collector: collector}
}

func (this *PipelineCollector) Process(items *result.ResultItems, t task.Task) {
	this.collector = append(this.collector, items)
}

func (this *PipelineCollector) GetCollected() []*result.ResultItems {
	return this.collector
}
