package buckets

import (
	"math"
	"reflect"
	"testing"
	"unsafe"

	"strings"

	"fmt"

	"runtime"

	"path/filepath"

	"os"

	"strconv"

	"bou.ke/monkey"
)

var bucketIndex *int
var bucketCount *int
var directoriesToExclude *string
var packagesToExclude *string

var directoriesToExcludeList []string
var packagesToExcludeList []string

func init() {
	if v := os.Getenv("BUCKET"); v != "" {
		n, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			panic(fmt.Sprintf("unable to parse BUCKET %s: %v", v, err))
		}
		i := int(n)
		bucketIndex = &i
	}

	if v := os.Getenv("TOTAL_BUCKETS"); v != "" {
		n, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			panic(fmt.Sprintf("unable to parse BUCKET_COUNT %s: %v", v, err))
		}
		i := int(n)
		bucketCount = &i
	}

	if v := os.Getenv("EXCLUDE_DIRECTORIES"); v != "" {
		directoriesToExclude = &v
	}

	if v := os.Getenv("EXCLUDE_PACKAGES"); v != "" {
		directoriesToExclude = &v
	}

	if directoriesToExclude != nil {
		directoriesToExcludeList = strings.FieldsFunc(*directoriesToExclude, func(r rune) bool {
			return r == ',' || r == ';'
		})
	}
	if packagesToExclude != nil {
		packagesToExcludeList = strings.FieldsFunc(*packagesToExclude, func(r rune) bool {
			return r == ',' || r == ';'
		})
	}

	if !isBucketFeatureEnabled() && !isDirectoriesToExcludeEnabled() && !isPackagesToExcludeEnabled() {
		return
	}

	patchTestRun()
}

func isBucketFeatureEnabled() bool {
	if bucketCount == nil || bucketIndex == nil {
		return false
	}

	if *bucketCount <= 0 || *bucketIndex < 0 || *bucketIndex >= *bucketCount {
		return false
	}
	return true
}

func isDirectoriesToExcludeEnabled() bool {
	return len(directoriesToExcludeList) > 0
}

func isPackagesToExcludeEnabled() bool {
	return len(packagesToExcludeList) > 0
}

func getSourceFile(f func(*testing.T)) string {
	v := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	if v == nil {
		return ""
	}
	file, _ := v.FileLine(0)
	return file
}

func getPackageName(f func(*testing.T)) string {
	v := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	if v == nil {
		return ""
	}
	name := v.Name()

	// find the last slash
	lastSlash := strings.LastIndexFunc(name, func(r rune) bool {
		return r == '/'
	})
	if lastSlash <= -1 {
		lastSlash = 0
	}

	dot := strings.IndexRune(name[lastSlash:], '.')
	if dot < 0 {
		// no dot means no package
		return ""
	}
	dot += lastSlash
	return name[:dot]
}

func isFileInDir(file string, dirs ...string) bool {
dirLoop:
	for _, dir := range dirs {
		fileParts := strings.FieldsFunc(file, func(r rune) bool {
			return r == filepath.Separator
		})
		dirParts := strings.FieldsFunc(dir, func(r rune) bool {
			return r == filepath.Separator
		})

		if len(fileParts) < len(dirParts) {
			continue dirLoop
		}

		for i, part := range dirParts {
			if fileParts[i] != part {
				continue dirLoop
			}
		}
		return true
	}
	return false
}

func filterTests(tests *[]testing.InternalTest) {
	if isDirectoriesToExcludeEnabled() {
		for i := len(*tests) - 1; i >= 0; i-- {
			file := getSourceFile((*tests)[i].F)
			if file == "" {
				fmt.Printf("unable to find source of %s\n", (*tests)[i].Name)
				continue
			}
			if isFileInDir(file, directoriesToExcludeList...) {
				*tests = append((*tests)[:i], (*tests)[i+1:]...)
			}
		}
	}
	if isPackagesToExcludeEnabled() {
		for i := len(*tests) - 1; i >= 0; i-- {
			pkg := getPackageName((*tests)[i].F)
			if pkg == "" {
				fmt.Printf("unable to find package of %s\n", (*tests)[i].Name)
				continue
			}
			if isFileInDir(pkg, packagesToExcludeList...) {
				*tests = append((*tests)[:i], (*tests)[i+1:]...)
			}
		}
	}
	if isBucketFeatureEnabled() {
		perBucket := int(math.Ceil(float64(len(*tests)) / float64(*bucketCount)))

		from := *bucketIndex * perBucket
		to := from + perBucket

		if to > len(*tests) {
			to = len(*tests)
		}

		*tests = (*tests)[from:to]
	}
}

func patchTestRun() {
	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(&testing.M{}), "Run", func(m *testing.M) int {
		guard.Unpatch()
		defer guard.Restore()

		v := reflect.ValueOf(m).Elem()
		testsField := v.FieldByName("tests")
		ptr := unsafe.Pointer(testsField.UnsafeAddr())
		filterTests((*[]testing.InternalTest)(ptr))
		return m.Run()
	})
}
