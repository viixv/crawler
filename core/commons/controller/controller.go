package controller

type GoroutineController interface {
	GetOne()
	FreeOne()
	Has() uint
	Left() uint
}
