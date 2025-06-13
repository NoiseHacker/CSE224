#!/usr/bin/env bash

# Usage: ./gen_inputs.sh <MODE> <N> <SIZE>
# MODE: split (default), multi, allzero
# N: log2(number_of_nodes)
# SIZE: size in MiB

MODE=$1
N=$2
SIZE=$3

NUM_NODES=$((1<<N))
GENSORT="utils/win-amd64/bin/gensort-amd64.exe"
SPLITTER="split_records.go"
INPUT_DIR="input"
RECORD_FILE="input_all.dat"

mkdir -p "$INPUT_DIR"
rm -f "$INPUT_DIR"/*.dat
rm -f ref.txt

case "$MODE" in
  multi)
    echo "Generating ${SIZE} MiB per-node input for ${NUM_NODES} nodes..."
    for ((i=0; i<NUM_NODES; i++)); do
      "$GENSORT" "${SIZE} mb" "$INPUT_DIR/input_$i.dat"
    done
    echo "Generating reference sorted output..."
    cat "$INPUT_DIR"/input_*.dat > "$RECORD_FILE"
    utils/win-amd64/bin/showsort-amd64.exe "$RECORD_FILE" | sort > ref.txt
    rm -f "$RECORD_FILE"
    ;;

  allzero)
    echo "Generating ${SIZE} MiB records for single input..."
    "$GENSORT" "${SIZE} mb" "$RECORD_FILE"
    echo "Assigning all records to node 0..."
    mv "$RECORD_FILE" "input/input_0.dat"
    for ((i=1; i<NUM_NODES; i++)); do
      :> "input/input_$i.dat"
    done
    echo "Generating reference sorted output..."
    cat "$INPUT_DIR"/input_*.dat > "$RECORD_FILE"
    utils/win-amd64/bin/showsort-amd64.exe "$RECORD_FILE" | sort > ref.txt
    rm -f "$RECORD_FILE"
    ;;

  split)
    echo "Generating ${SIZE} MiB records..."
    "$GENSORT" "${SIZE} mb" "$RECORD_FILE"
    echo "Splitting for ${NUM_NODES} nodes..."
    go run "$SPLITTER" --input "$RECORD_FILE" --outdir "input" --nodes "$NUM_NODES"
    echo "Generating reference sorted output..."
    cat "$INPUT_DIR"/input_*.dat > "$RECORD_FILE"
    utils/win-amd64/bin/showsort-amd64.exe "$RECORD_FILE" | sort > ref.txt
    rm -f "$RECORD_FILE"
    ;;

  *)
    echo "Unknown mode: $MODE"
    exit 1
    ;;
esac

echo "Done."