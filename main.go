// Package buckets provides a way to separate go tests into buckets.
package buckets

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

// Buckets must be called to get the test bucket feature working.
// It will modify the tests present in the testing.M struct.
func Buckets(m *testing.M) {
	var bucketIndex *int
	var bucketCount *int

	var directoriesToExclude []string
	var packagesToExclude []string

	if v := os.Getenv("BUCKET"); v != "" {
		//nolint: gomnd // use 64 bits for parsing
		n, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			panic(fmt.Sprintf("unable to parse BUCKET %s: %v", v, err))
		}
		i := int(n)
		bucketIndex = &i
	}

	if v := os.Getenv("TOTAL_BUCKETS"); v != "" {
		//nolint: gomnd // use 64 bits for parsing
		n, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			panic(fmt.Sprintf("unable to parse BUCKET_COUNT %s: %v", v, err))
		}
		i := int(n)
		bucketCount = &i
	}

	if v := os.Getenv("EXCLUDE_DIRECTORIES"); v != "" {
		directoriesToExclude = strings.FieldsFunc(v, func(r rune) bool {
			return r == ',' || r == ';'
		})
		for i := range directoriesToExclude {
			directoriesToExclude[i] = filepath.ToSlash(directoriesToExclude[i])
		}
	}

	if v := os.Getenv("EXCLUDE_PACKAGES"); v != "" {
		packagesToExclude = strings.FieldsFunc(v, func(r rune) bool {
			return r == ',' || r == ';'
		})
	}

	if bucketIndex == nil && bucketCount == nil && len(directoriesToExclude) == 0 && len(packagesToExclude) == 0 {
		return
	}

	v := reflect.ValueOf(m).Elem()
	testsField := v.FieldByName("tests")
	//nolint: gosec // allow the usage of unsafe so we can get the test slice.
	ptr := unsafe.Pointer(testsField.UnsafeAddr())
	filterTests((*[]testing.InternalTest)(ptr), *bucketIndex, *bucketCount, directoriesToExclude, packagesToExclude)
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
			return r == '/'
		})
		dirParts := strings.FieldsFunc(dir, func(r rune) bool {
			return r == '/'
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

func filterTests(tests *[]testing.InternalTest, bucketIndex, bucketCount int, directoriesToExcludeList []string, packagesToExcludeList []string) {
	if len(directoriesToExcludeList) > 0 {
		for i := len(*tests) - 1; i >= 0; i-- {
			file := getSourceFile((*tests)[i].F)
			if file == "" {
				fmt.Printf("unable to find source of %s\n", (*tests)[i].Name)
				continue
			}

			if isFileInDir(filepath.ToSlash(file), directoriesToExcludeList...) {
				*tests = append((*tests)[:i], (*tests)[i+1:]...)
			}
		}
	}
	if len(packagesToExcludeList) > 0 {
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

	if bucketCount > 0 && bucketIndex >= 0 && bucketIndex < bucketCount {
		perBucket := int(math.Ceil(float64(len(*tests)) / float64(bucketCount)))

		from := bucketIndex * perBucket
		to := from + perBucket

		if to > len(*tests)-1 {
			to = len(*tests)
		}

		*tests = (*tests)[from:to]
	}
}
