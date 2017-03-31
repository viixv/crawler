package pipeline

import (
	"os"

	"github.com/viixv/crawler/core/commons/result"
	"github.com/viixv/crawler/core/commons/task"
)

type FilePipeline struct {
	pFile *os.File

	path string
}

func NewFilePipeline(path string) *FilePipeline {
	pFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic("File '" + path + "' in PipelineFile open failed.")
	}
	return &FilePipeline{path: path, pFile: pFile}
}

func (this *FilePipeline) Process(items *result.ResultItems, t task.Task) {
	this.pFile.WriteString("----------------------------------------------------------------------------------------------\n")
	this.pFile.WriteString("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
	this.pFile.WriteString("Crawled result : \n")
	for key, value := range items.GetAll() {
		this.pFile.WriteString(key + "\t:\t" + value + "\n")
	}
}
