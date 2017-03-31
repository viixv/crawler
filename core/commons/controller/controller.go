// Package resource_manage implements a resource management.
package controller

// ResourceManage is an interface that who want implement an management object can realize these functions.
type GoroutineController interface {
	GetOne()
	FreeOne()
	Has() uint
	Left() uint
}
