// Package pipeline is the persistent and offline process part of crawler.
package pipeline

import (
	"github.com/viixv/crawler/core/commons/result"
	"github.com/viixv/crawler/core/commons/task"
)

// The interface Pipeline can be implemented to customize ways of persistent.
type Pipeline interface {
	// The Process implements result persistent.
	// The items has the result be crawled.
	// The t has informations of this crawl task.
	Process(items *result.ResultItems, t task.Task)
}

// The interface CollectPipeline recommend result in process's memory temporarily.
type CollectPipeline interface {
	Pipeline

	// The GetCollected returns result saved in in process's memory temporarily.
	GetCollected() []*result.ResultItems
}
