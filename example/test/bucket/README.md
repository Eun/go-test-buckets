# Buckets
To run this example
```shell
# This will only run TestA and TestB (Bucket 0) not TestC (Bucket 1)
cd example/test/exclude/directory
BUCKET=0 TOTAL_BUCKETS=2 go test -count=1 -v ./...
```

