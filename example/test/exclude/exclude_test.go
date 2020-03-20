package exclude

import (
	"testing"

	_ "github.com/Eun/go-test-buckets"
)

// run with EXCLUDE_PACKAGES=github.com/Eun/go-test-buckets/example/test/exclude go test -count=1 -v ./...
func TestExclude(t *testing.T) {
}
