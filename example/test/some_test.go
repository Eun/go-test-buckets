package test

import (
	"testing"

	_ "github.com/Eun/go-test-buckets"
)

// run with go test -bucket=0 -total-buckets=2 -v
// will run TestA and TestB

// run with go test -bucket=1 -total-buckets=2 -v
// will run TestC

func TestA(t *testing.T) {
	t.Run("TestA1", func(t *testing.T) {})
	t.Run("TestA2", func(t *testing.T) {})
}

func TestB(t *testing.T) {
	t.Run("TestB1", func(t *testing.T) {})
}

func TestC(t *testing.T) {
}
