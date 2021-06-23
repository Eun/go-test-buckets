# go-test-buckets
Split your go tests into buckets, or exclude some packages/directories.

## Buckets
```golang
package main_test

import (
	"testing"

	"github.com/Eun/go-test-buckets"
)

func TestMain(m *testing.M) {
	buckets.Buckets(m)
	m.Run()
}

// run with BUCKET=0 TOTAL_BUCKETS=2 go test -count=1 -v ./...
// will run TestA and TestB

// run with BUCKET=1 TOTAL_BUCKETS=2 go test -count=1 -v ./...
// will run TestC

func TestA(t *testing.T) {
}

func TestB(t *testing.T) {
}

func TestC(t *testing.T) {
}
```

## Excluding Packages/Directories
1. Add to your package
   ```go
   func TestMain(m *testing.M) {
       buckets.Buckets(m)
       m.Run()
   }
   ```
2. Run `EXCLUDE_PACKAGES=package/path/to/exclude,package/path/to/exclude-2 go test -count=1 -v ./...`
3. Or `EXCLUDE_DIRECTORIES=/full/path/to/exclude go test -count=1 -v ./...`



> Because of `go test` package separation, you have to call `buckets.Buckets(m)` in every package you want to ignore.

# Why?
Speed up ci pipelines by parallelizing go tests without thinking about [t.Parallel](https://golang.org/pkg/testing/#T.Parallel).
And getting rid of weird piping `go test $(go list ./... | grep -v /ignore/)`