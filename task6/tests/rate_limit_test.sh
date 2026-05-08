#!/bin/bash

API_URL="http://localhost:28080"
ENDPOINT="/api/ping"

echo "=========================================="
echo "Testing Rate Limiter"
echo "Config: 5 requests per second, burst 10"
echo "=========================================="
echo ""

echo "Phase 1: Sending 12 sequential requests..."
count=0
blocked=0
for i in {1..12}; do
    code=$(curl -s -o /dev/null -w "%{http_code}" "${API_URL}${ENDPOINT}")
    echo "Request $i: $code"
    if [ "$code" -eq 429 ]; then
        blocked=$((blocked + 1))
    fi
    count=$((count + 1))
done

echo ""
echo "Blocked requests: $blocked out of $count"

echo ""
echo "Phase 2: Waiting 2 seconds for tokens to replenish..."
sleep 2

echo "Sending 3 more requests..."
for i in {1..3}; do
    code=$(curl -s -o /dev/null -w "%{http_code}" "${API_URL}${ENDPOINT}")
    echo "Request after wait $i: $code"
done

echo ""
echo "=========================================="
echo "Test completed."