package buckets

import (
	"flag"
	"math"
	"reflect"
	"testing"
	"unsafe"

	"bou.ke/monkey"
)

var bucketIndex = flag.Int("bucket", 0, "bucket index to use")
var bucketCount = flag.Int("total-buckets", 0, "total bucket count")

func init() {
	flag.Parse()

	if bucketCount == nil || bucketIndex == nil {
		return
	}

	if *bucketCount <= 0 || *bucketIndex < 0 || *bucketIndex >= *bucketCount {
		return
	}

	patchTestRun()
}

func patchTestRun() {
	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(&testing.M{}), "Run", func(m *testing.M) int {
		guard.Unpatch()
		defer guard.Restore()

		v := reflect.ValueOf(m).Elem()
		testsField := v.FieldByName("tests")
		ptr := unsafe.Pointer(testsField.UnsafeAddr())
		tests := (*[]testing.InternalTest)(ptr)

		perBucket := int(math.Ceil(float64(len(*tests)) / float64(*bucketCount)))

		from := *bucketIndex * perBucket
		to := from + perBucket

		if to > len(*tests) {
			to = len(*tests)
		}

		*tests = (*tests)[from:to]

		return m.Run()
	})
}
