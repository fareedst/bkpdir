#!/usr/bin/env bash
set -euo pipefail

declare -A VALIDATION_RESULTS
declare -A DECISION_METRICS

# Initialize some values
VALIDATION_RESULTS[test1]=true
VALIDATION_RESULTS[test2]=false
DECISION_METRICS[test_metric]=100

echo "Testing associative arrays:"
echo "VALIDATION_RESULTS keys: ${!VALIDATION_RESULTS[@]}"
echo "DECISION_METRICS keys: ${!DECISION_METRICS[@]}"

for key in "${!VALIDATION_RESULTS[@]}"; do
    echo "  $key: ${VALIDATION_RESULTS[$key]}"
done 