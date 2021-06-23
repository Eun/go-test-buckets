#!/bin/bash
go get -u ./...
tests=$(go test -v -count=1 -coverprofile="coverage-$1.cov" -covermode=atomic -json $2 | jq -c -s '.[] | select(.Action == "run") | .Test')

IFS=';' read -ra SHOULD_RUN <<< "$TEST_SHOULD_RUN"
for i in "${SHOULD_RUN[@]}"; do
  if [ $(echo -ne "$tests" | jq -s "index(\"$i\")") == "null" ]; then
      echo "$i" should have been run
      exit 1
  fi
done

IFS=';' read -ra SHOULD_NOT_RUN <<< "$TEST_SHOULD_NOT_RUN"
for i in "${SHOULD_NOT_RUN[@]}"; do
    if [ $(echo -ne "$tests" | jq -s "index(\"$i\")") != "null" ]; then
        echo "$i" should have not been run
        exit 1
    fi
done

