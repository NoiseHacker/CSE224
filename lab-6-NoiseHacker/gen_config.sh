#!/bin/bash

N=$1
NUM_NODES=$((1 << N))
PORT_BASE=8000
CONFIG_FILE=config.yaml

echo "nodes:" > $CONFIG_FILE
for ((i = 0; i < NUM_NODES; i++)); do
  echo "  - nodeID: $i" >> $CONFIG_FILE
  echo "    host: \"localhost\"" >> $CONFIG_FILE
  echo "    port: $((PORT_BASE + i))" >> $CONFIG_FILE
done

echo "Generated $CONFIG_FILE with $NUM_NODES nodes"
