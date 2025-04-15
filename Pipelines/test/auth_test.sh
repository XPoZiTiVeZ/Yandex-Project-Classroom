#!/bin/bash

MICROSERVICE_DIRNAME="Auth"

cd $MICROSERVICE_DIRNAME
PATHS_FOR_TESTS=$(go list ./...)

OUTPUT=$(go test -v $PATHS_FOR_TESTS)
EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
    echo "❌ Tests failed!"
    FAILED_TESTS=$(echo "$OUTPUT" | grep -E "=== (PASS|FAIL)|^--- PASS:|^--- FAIL:" | grep "^--- FAIL:")
    echo "$FAILED_TESTS"
    exit 1
else
    echo "✅ All tests passed!"  
fi

cd ..