# Exclude Packages
It is possible to exclude packages directories.
However this has 1 caveats:
1. `buckets.Buckets()` must be called inside the TestMain() function (also for every nested package!).
2. The environment variable `EXCLUDE_DIRECTORIES` must be set to an absolute path.  
   (go test always navigates to the test folder)
   
To run this example
```shell
cd example/test/exclude/directory
EXCLUDE_PACKAGES=github.com/Eun/go-test-buckets/example/test/exclude/package go test -count=1 -v ./...
```

You can specify multiple packages by joining them with an `;`:
```
EXCLUDE_PACKAGES=package1;package2
```