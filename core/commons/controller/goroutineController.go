package controller

// ResourceManageChan inherits the ResourceManage interface.
// In spider, ResourceManageChan manage resource of Coroutine to crawl page.
type GoroutineControllerChan struct {
	capnum uint
	mc     chan uint
}

// NewResourceManageChan returns initialized ResourceManageChan object which contains a resource pool.
// The num is the resource limit.
func NewGoroutineControllerChan(num uint) *GoroutineControllerChan {
	mc := make(chan uint, num)
	return &GoroutineControllerChan{mc: mc, capnum: num}
}

// The GetOne apply for one resource.
// If resource pool is empty, current coroutine will be blocked.
func (this *GoroutineControllerChan) GetOne() {
	this.mc <- 1
}

// The FreeOne free resource and return it to resource pool.
func (this *GoroutineControllerChan) FreeOne() {
	<-this.mc
}

// The Has query for how many resource has been used.
func (this *GoroutineControllerChan) Has() uint {
	return uint(len(this.mc))
}

// The Left query for how many resource left in the pool.
func (this *GoroutineControllerChan) Left() uint {
	return this.capnum - uint(len(this.mc))
}
