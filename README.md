# go-test-buckets
Split your go tests into buckets.

```golang
package main_test

import (
	"testing"

	_ "github.com/Eun/go-test-buckets"
)

// run with go test -bucket=0 -total-buckets=2 -v
// will run TestA and TestB

// run with go test -bucket=1 -total-buckets=2 -v
// will run TestC

func TestA(t *testing.T) {
}

func TestB(t *testing.T) {
}

func TestC(t *testing.T) {
}
```

> Note that uses some nasty memory patching to make this possible. So use with care.

# Why?
Speed up ci pipelines by parallelizing go tests without thinking about [t.Parallel](https://golang.org/pkg/testing/#T.Parallel).