# Exclude Directories
It is possible to exclude specific directories.
However this has 2 caveats:
1. `buckets.Buckets()` must be called inside the TestMain() function (also for every nested directory!).
2. The environment variable `EXCLUDE_DIRECTORIES` must be set to an absolute path.  
   (go test always navigates to the test folder)
   
To run this example
```shell
cd example/test/exclude/directory
EXCLUDE_DIRECTORIES=$(pwd)/should_be_ignored go test -count=1 -v ./...
```

You can specify multiple directories by joining them with an `;`:
```
EXCLUDE_DIRECTORIES=package1;package2
```