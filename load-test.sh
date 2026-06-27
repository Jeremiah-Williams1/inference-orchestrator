#!/usr/bin/env bash
# load-test.sh — fires valid + invalid jobs at the inference API to visualize
# in Grafana (Loki logs + Prometheus metrics).
#
# Usage:
#   export INFERENCE_API_URL="http://192.168.49.2:32135"
#   chmod +x load-test.sh
#   ./load-test.sh

if [ -z "$INFERENCE_API_URL" ]; then
  echo "ERROR: INFERENCE_API_URL is not set."
  echo "  export INFERENCE_API_URL=\"http://192.168.49.2:32135\""
  exit 1
fi

ENDPOINT="/api/v1/jobs"

valid_payload() {
  cat <<EOF
{
  "type": "classification",
  "input": {
    "Time_spent_Alone": $((RANDOM % 12)),
    "Social_event_attendance": $((RANDOM % 30)),
    "Going_outside": $((RANDOM % 7)),
    "Friends_circle_size": $((RANDOM % 15)),
    "Post_frequency": $((RANDOM % 20)),
    "Stage_fear": "Yes"
  }
}
EOF
}

echo "Submitting valid jobs..."
for i in $(seq 1 8); do
  curl -s -o /dev/null -w "valid #$i -> %{http_code}\n" \
    -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
    -H "Content-Type: application/json" \
    -d "$(valid_payload)"
  sleep 0.3
done

echo "Submitting jobs with bad input (missing required fields)..."
for i in $(seq 1 4); do
  curl -s -o /dev/null -w "bad-input #$i -> %{http_code}\n" \
    -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
    -H "Content-Type: application/json" \
    -d '{"type": "classification", "input": {"Time_spent_Alone": "not-a-number"}}'
  sleep 0.3
done

echo "Submitting unknown job type..."
curl -s -o /dev/null -w "bad-type -> %{http_code}\n" \
  -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
  -H "Content-Type: application/json" \
  -d '{"type": "not-a-real-type", "input": {}}'

echo "Submitting malformed JSON..."
curl -s -o /dev/null -w "malformed -> %{http_code}\n" \
  -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
  -H "Content-Type: application/json" \
  -d '{not valid json'

echo "Done. Check Grafana Explore for the queue-depth spike + error logs."