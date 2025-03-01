// Package app provides a wrapper around the Arche ECS world
// that helps with rapid prototyping and app/model development.
package app

// init initializes the package.
// It adjusts time resolution for Windows.
// See this issue: https://github.com/golang/go/issues/44343
func init() {
	initTimer()
}
