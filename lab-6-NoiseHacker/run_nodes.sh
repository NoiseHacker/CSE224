#!/bin/bash

N=$1
NUM_NODES=$((1 << N))
CONFIG=config.yaml
INPUT_DIR=input
OUTPUT_DIR=output
LOG_DIR=log

mkdir -p $OUTPUT_DIR $LOG_DIR
rm $OUTPUT_DIR/*.dat
rm $LOG_DIR/*.log

BIN="go run cmd/globesort/main.go"
PIDS=()
EXIT_CODES=()

for ((i = 0; i < NUM_NODES; i++)); do
  echo "Starting node $i..."
  $BIN $i $INPUT_DIR/input_$i.dat $OUTPUT_DIR/output_$i.dat $CONFIG > $LOG_DIR/node_$i.log 2>&1 &
  PIDS[i]=$!
done

# wait for all the nodes to exit
for ((i = 0; i < NUM_NODES; i++)); do
  wait ${PIDS[i]}
  CODE=$?
  EXIT_CODES[i]=$CODE
  if [ $CODE -ne 0 ]; then
    echo "❌ Node $i (PID ${PIDS[i]}) exited with code $CODE"
  fi
done

# check logs
for ((i = 0; i < NUM_NODES; i++)); do
  if grep -Ei 'panic|fatal|failed' "$LOG_DIR/node_$i.log" > /dev/null; then
    echo "⚠️  Warning: Suspicious logs found in node_$i.log"
  fi
done

echo "✅ All nodes processed. Exit codes: ${EXIT_CODES[*]}"
