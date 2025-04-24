#!/bin/bash
set +x  # Отключаем вывод команд

chmod +x test_cover.sh
echo "🔍 Running tests with coverage..."
COVERAGE=$(./test_cover.sh)
EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
  echo "❌ Unit tests failed!"
  exit 1
fi

echo "📊 Test coverage: $COVERAGE% (minimum required: $MIN_COVERAGE%)"
if [ $COVERAGE -lt $MIN_COVERAGE ]; then
  echo "❌ Coverage is below 30%"
  exit 1
fi