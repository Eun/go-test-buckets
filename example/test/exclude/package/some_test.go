package _package

import (
	"os"
	"testing"

	buckets "github.com/Eun/go-test-buckets"
)

func TestMain(m *testing.M) {
	buckets.Buckets(m)
	os.Exit(m.Run())
}

func TestSomething(t *testing.T) {
}
