// Package interfaces contains some common interface of crawler project.
package interfaces

// The Task represents interface that contains environment variables.
// It inherits by Spider.
type Task interface {
	TaskName() string
}
