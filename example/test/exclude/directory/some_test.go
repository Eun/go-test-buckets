package _package

import (
	"testing"

	buckets "github.com/Eun/go-test-buckets"
)

func TestMain(m *testing.M) {
	buckets.Buckets(m)
	m.Run()
}

func TestSomething(t *testing.T) {
}
