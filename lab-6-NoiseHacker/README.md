[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/vemBksKM)
# GlobeSort: Distributed Sort

## Building

```bash
go build -o bin/globesort.exe cmd/globesort/main.go
```

## Usage

```bash
bin/globesort.exe <nodeID> <inputFilePath> <outputFilePath> <configFilePath>
```

Example:

```bash
bin/globesort.exe 0 inputs/input_0.dat outputs/sorted_0.dat config.yaml
```

Multi-node Example:

```bash
nohup bin/globesort.exe 0 inputs/input_0.dat outputs/sorted_0.dat config.yaml &
nohup bin/globesort.exe 1 inputs/input_1.dat outputs/sorted_1.dat config.yaml &
nohup bin/globesort.exe 2 inputs/input_2.dat outputs/sorted_2.dat config.yaml &
nohup bin/globesort.exe 3 inputs/input_3.dat outputs/sorted_3.dat config.yaml &
wait  # block until all nohup processes are finished
```

nohup will redirect all output to `nohup.out` file.

```bash
cat nohup.out
# or, to see output as it is generated
tail -f nohup.out
```

## Testing

You can use the same tools in Lab 1 to test your sort:

```bash
cp utils/<os-arch>/bin . -r  # Replace <os-arch> with your OS and architecture

utils/win-amd64/bin/gensort-amd64.exe "1 mb" 1mb.dat
utils/win-amd64/bin/showsort-amd64.exe  1mb.dat | sort > 1mb_sorted.txt
```

For multi-node results, you can concatenate the sorted files on each node to a single file:

```bash
# start by generating your input files
utils/win-amd64/bin/gensort-amd64.exe "1 mb" inputs/input_0.dat
utils/win-amd64/bin/gensort-amd64.exe "1 mb" inputs/input_1.dat
utils/win-amd64/bin/gensort-amd64.exe "1 mb" inputs/input_2.dat
utils/win-amd64/bin/gensort-amd64.exe "1 mb" inputs/input_3.dat

# sort the input files
nohup bin/globesort.exe 0 inputs/input_0.dat outputs/sorted_0.dat config.yaml &
nohup bin/globesort.exe 1 inputs/input_1.dat outputs/sorted_1.dat config.yaml &
nohup bin/globesort.exe 2 inputs/input_2.dat outputs/sorted_2.dat config.yaml &
nohup bin/globesort.exe 3 inputs/input_3.dat outputs/sorted_3.dat config.yaml &
wait  # block until all nohup processes are finished

# generate reference text file
cat inputs/input_*.dat > input_all.dat
utils/win-amd64/bin/showsort-amd64.exe input_all.dat | sort > ref.txt

# show your multinode results
cat outputs/sorted_*.dat > sorted_all.dat
utils/win-amd64/bin/showsort-amd64.exe sorted_all.dat > sorted_all.txt

# check if they are different
diff sorted_all.txt ref.txt
```
