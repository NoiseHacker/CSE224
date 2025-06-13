# test.sh
#!/usr/bin/env bash

# Usage: ./test.sh <N> <SIZE> [split|multi|allzero]
N=$1
SIZE=$2
MODE=${3:-split}

set -euo pipefail

echo "🔧 Generating config..."
./gen_config.sh "$N"

echo "📦 Generating and splitting input..."
./gen_inputs.sh "$MODE" "$N" "$SIZE"

echo "🚀 Running all nodes..."
./run_nodes.sh "$N"

echo "🧩 Merging and transforming output..."
mkdir -p temp
ls output/output_*.dat | sort -V | xargs cat > temp/all_output.dat
utils/win-amd64/bin/showsort-amd64.exe temp/all_output.dat > temp/test.txt

echo "🔍 Comparing output with reference..."
if diff -q ref.txt temp/test.txt; then
  echo "✅ SUCCESS: Output matches reference."
  rm -rf input output log temp ref.txt
  exit 0
else
  echo "❌ MISMATCH DETECTED!"
  exit 1
fi