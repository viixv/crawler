package controller

type GoroutineControllerChan struct {
	capNum     uint
	goroutines chan uint
}

func NewGoroutineControllerChan(num uint) *GoroutineControllerChan {
	return &GoroutineControllerChan{goroutines: make(chan uint, num), capNum: num}
}

func (this *GoroutineControllerChan) GetOne() {
	this.goroutines <- 1
}

func (this *GoroutineControllerChan) FreeOne() {
	<-this.goroutines
}

func (this *GoroutineControllerChan) Has() uint {
	return uint(len(this.goroutines))
}

func (this *GoroutineControllerChan) Left() uint {
	return this.capNum - uint(len(this.goroutines))
}
