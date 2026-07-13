#!/usr/bin/env bash
# load-test.sh — fires valid + invalid jobs at the inference API to visualize
# in Grafana (Loki logs + Prometheus metrics).
#
# Usage:
#   export INFERENCE_API_URL="http://<minikube-ip>:<nodeport>"
#   chmod +x load-test.sh
#   ./load-test.sh

if [ -z "$INFERENCE_API_URL" ]; then
  echo "ERROR: INFERENCE_API_URL is not set."
  echo "  export INFERENCE_API_URL=\"http://192.168.49.2:<nodeport>\""
  exit 1
fi

ENDPOINT="/api/v1/jobs"

# ---------------------------------------------------------------------------
# Classification payloads
# ---------------------------------------------------------------------------
classification_payload() {
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

# ---------------------------------------------------------------------------
# Regression payloads -- Auto MPG features
# Cylinders: 4/6/8, Displacement: 70-455, Horsepower: 50-230,
# Weight: 1600-5000, Acceleration: 8-25, Model Year: 70-82
# Origin: USA / Europe / Japan
# ---------------------------------------------------------------------------
ORIGINS=("USA" "Europe" "Japan")
CYLINDERS=(4 6 8)

regression_payload() {
  origin=${ORIGINS[$((RANDOM % 3))]}
  cylinders=${CYLINDERS[$((RANDOM % 3))]}
  cat <<EOF
{
  "type": "regression",
  "input": {
    "Cylinders": ${cylinders},
    "Displacement": $((RANDOM % 386 + 70)),
    "Horsepower": $((RANDOM % 181 + 50)),
    "Weight": $((RANDOM % 3401 + 1600)),
    "Acceleration": $((RANDOM % 18 + 8)),
    "Model Year": $((RANDOM % 13 + 70)),
    "Origin": "${origin}"
  }
}
EOF
}

# ---------------------------------------------------------------------------
# Valid jobs -- both workers
# ---------------------------------------------------------------------------
echo "--- Submitting valid classification jobs ---"
for i in $(seq 1 8); do
  curl -s -o /dev/null -w "classification valid #$i -> %{http_code}\n" \
    -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
    -H "Content-Type: application/json" \
    -d "$(classification_payload)"
  sleep 0.3
done

echo "--- Submitting valid regression jobs ---"
for i in $(seq 1 8); do
  curl -s -o /dev/null -w "regression valid #$i -> %{http_code}\n" \
    -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
    -H "Content-Type: application/json" \
    -d "$(regression_payload)"
  sleep 0.3
done

# ---------------------------------------------------------------------------
# Bad input -- missing required fields
# ---------------------------------------------------------------------------
echo "--- Submitting classification jobs with bad input ---"
for i in $(seq 1 3); do
  curl -s -o /dev/null -w "classification bad-input #$i -> %{http_code}\n" \
    -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
    -H "Content-Type: application/json" \
    -d '{"type": "classification", "input": {"Time_spent_Alone": "not-a-number"}}'
  sleep 0.3
done

echo "--- Submitting regression jobs with bad input ---"
for i in $(seq 1 3); do
  curl -s -o /dev/null -w "regression bad-input #$i -> %{http_code}\n" \
    -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
    -H "Content-Type: application/json" \
    -d '{"type": "regression", "input": {"Cylinders": "not-a-number"}}'
  sleep 0.3
done

# ---------------------------------------------------------------------------
# Edge cases
# ---------------------------------------------------------------------------
echo "--- Submitting unknown job type ---"
curl -s -o /dev/null -w "bad-type -> %{http_code}\n" \
  -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
  -H "Content-Type: application/json" \
  -d '{"type": "not-a-real-type", "input": {}}'

echo "--- Submitting regression job with invalid Origin ---"
curl -s -o /dev/null -w "bad-origin -> %{http_code}\n" \
  -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
  -H "Content-Type: application/json" \
  -d '{"type": "regression", "input": {"Cylinders": 4, "Displacement": 120.0, "Horsepower": 80.0, "Weight": 2200.0, "Acceleration": 16.0, "Model Year": 78, "Origin": "Mars"}}'

echo "--- Submitting malformed JSON ---"
curl -s -o /dev/null -w "malformed -> %{http_code}\n" \
  -X POST "${INFERENCE_API_URL}${ENDPOINT}" \
  -H "Content-Type: application/json" \
  -d '{not valid json'

echo ""
echo "Done. Open Grafana Explore:"
echo "  - Loki: filter by {app='classification-worker'} and {app='regression-worker'}"
echo "  - Prometheus: check inference_jobs_completed_total, inference_duration_seconds"
echo "  - Dashboard: queue depth spike should be visible during the run"
