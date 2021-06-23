package bucket

import (
	"os"
	"testing"

	buckets "github.com/Eun/go-test-buckets"
)

// run with BUCKET=0 TOTAL_BUCKETS=2 go test -count=1 -v ./...
// will run TestA and TestB

// run with BUCKET=1 TOTAL_BUCKETS=2 go test -count=1 -v ./...
// will run TestC

func TestMain(m *testing.M) {
	buckets.Buckets(m)
	os.Exit(m.Run())
}

func TestA(t *testing.T) {
	t.Run("TestA1", func(t *testing.T) {})
	t.Run("TestA2", func(t *testing.T) {})
}

func TestB(t *testing.T) {
	t.Run("TestB1", func(t *testing.T) {})
}

func TestC(t *testing.T) {
	t.Run("TestC1", func(t *testing.T) {})
}
