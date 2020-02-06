# go-test-buckets
Split your go tests into buckets, or exclude some packages/directories.

## Buckets
```golang
package main_test

import (
	"testing"

	// import to get bucket and exclude feature
	_ "github.com/Eun/go-test-buckets"
)

// run with go test -v ./... -bucket=0 -total-buckets=2
// will run TestA and TestB

// run with go test -v ./... -bucket=1 -total-buckets=2
// will run TestC

func TestA(t *testing.T) {
}

func TestB(t *testing.T) {
}

func TestC(t *testing.T) {
}
```

## Excluding Packages/Directories
1. Add to your imports `_ "github.com/Eun/go-test-buckets"`
2. Run `go test -v ./... -exclude-packages=package/path/to/exclude,package/path/to/exclude-2`
3. Or `go test -v ./... -exclude-directories=/full/path/to/exclude`



> Because of `go test` package separation, you have to import `github.com/Eun/go-test-buckets` in each package you want to ignore.

> Note that uses some nasty memory patching to make this possible. So use with care.

# Why?
Speed up ci pipelines by parallelizing go tests without thinking about [t.Parallel](https://golang.org/pkg/testing/#T.Parallel).
And getting rid of weird piping `go test $(go list ./... | grep -v /ignore/)`