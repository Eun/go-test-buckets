package buckets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuckets(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to get working dir: %+v", err)
	}
	t.Run("Buckets", func(t *testing.T) {
		t.Run("Bucket0", func(t *testing.T) {
			subTest(t,
				[]string{
					"BUCKET=0",
					"TOTAL_BUCKETS=2",
				},
				"bucket0",
				"./example/test/bucket/...",
				[]string{
					"TestA",
					"TestA/TestA1",
					"TestA/TestA2",
					"TestB",
					"TestB/TestB1",
				},
				[]string{
					"TestC",
					"TestC/TestC1",
				},
			)
		})
		t.Run("Bucket1", func(t *testing.T) {
			subTest(t,
				[]string{
					"BUCKET=1",
					"TOTAL_BUCKETS=2",
				},
				"bucket1",
				"./example/test/bucket/...",
				[]string{
					"TestC",
					"TestC/TestC1",
				},
				[]string{
					"TestA",
					"TestA/TestA1",
					"TestA/TestA2",
					"TestB",
					"TestB/TestB1",
				},
			)
		})
	})
	t.Run("Exclude", func(t *testing.T) {
		t.Run("Package", func(t *testing.T) {
			subTest(t,
				[]string{
					"EXCLUDE_PACKAGES=github.com/Eun/go-test-buckets/example/test/exclude/package/should_be_ignored",
				},
				"exclude-package",
				"./example/test/exclude/package/...",
				[]string{
					"TestSomething",
				},
				[]string{
					"TestSomethingWillBeExcluded",
				},
			)
		})
		t.Run("Directory", func(t *testing.T) {
			subTest(t,
				[]string{
					"EXCLUDE_DIRECTORIES=" + filepath.Join(wd, "example", "test", "exclude", "directory", "should_be_ignored"),
				},
				"exclude-directory",
				"./example/test/exclude/directory/...",
				[]string{
					"TestSomething",
				},
				[]string{
					"TestSomethingWillBeExcluded",
				},
			)
		})
	})
}

func subTest(t *testing.T, env []string, coverProfile, packageName string, testShouldRun, testShouldNotRun []string) {
	var buf bytes.Buffer
	coverProfile = os.Getenv("COVERAGE_PREFIX") + coverProfile
	cmd := exec.Command("go", "test", "-v", fmt.Sprintf("-coverprofile=%s.cov", coverProfile), "-covermode=atomic", "-json", packageName)
	cmd.Stdout = &buf
	cmd.Stderr = ioutil.Discard
	cmd.Env = append(os.Environ(), env...)

	if err := cmd.Run(); err != nil {
		t.Fatal(err, buf.String())
	}

	type test struct {
		Action string `json:"Action"`
		Test   string `json:"Test"`
		Output string `json:"Output"`
	}

	var allTests []test

	dec := json.NewDecoder(strings.NewReader(buf.String()))
	for {
		var testLine test

		err := dec.Decode(&testLine)
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		allTests = append(allTests, testLine)
	}

	testRan := func(tests []test, testName string) bool {
		for _, o := range tests {
			if o.Test == testName && o.Action == "run" {
				return true
			}
		}
		return false
	}

	testOutput := func(tests []test) string {
		var sb strings.Builder
		for _, o := range tests {
			if o.Action == "output" {
				sb.WriteString(o.Output)
			}
		}
		return sb.String()
	}

	for _, s := range testShouldRun {
		if !testRan(allTests, s) {
			t.Fatalf("Test `%s' did not run\nEnv:\n%s\nOutput:\n%s\n", s, strings.Join(cmd.Env, "\n"), testOutput(allTests))
		}
	}

	for _, s := range testShouldNotRun {
		if testRan(allTests, s) {
			t.Fatalf("Test `%s' should not run\nEnv:\n%s\nOutput:\n%s\n", s, strings.Join(cmd.Env, "\n"), testOutput(allTests))
		}
	}
}

func Test_filterTestsBuckets(t *testing.T) {
	testCases := []struct {
		name        string
		testCount   int
		bucketCount int
		buckets     []int
	}{
		{
			"1/1",
			1,
			1,
			[]int{1},
		},
		{
			"1/2",
			1,
			2,
			[]int{1, 0},
		},
		{
			"1/3",
			1,
			3,
			[]int{1, 0, 0},
		},
		{
			"3/2",
			3,
			2,
			[]int{2, 1},
		},
		{
			"4/2",
			4,
			2,
			[]int{2, 2},
		},
		{
			"5/2",
			5,
			2,
			[]int{3, 2},
		},
		{
			"11/3",
			11,
			3,
			[]int{4, 4, 3},
		},
		{
			"12/3",
			12,
			3,
			[]int{4, 4, 4},
		},
		{
			"13/3",
			13,
			3,
			[]int{5, 5, 3},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.bucketCount; i++ {
				tests := make([]testing.InternalTest, tt.testCount)
				for j := 0; j < tt.testCount; j++ {
					tests[j].Name = fmt.Sprintf("Test%d", j)
				}
				filterTests(&tests, i, tt.bucketCount, nil, nil)

				if tt.buckets[i] != len(tests) {
					t.Fatalf("expected %d tests in bucket %d but got %d", tt.buckets[i], i, len(tests))
				}
			}
		})
	}
}
