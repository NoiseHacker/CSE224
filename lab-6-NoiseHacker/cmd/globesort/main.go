package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/bits"
	"net"
	"os"
	"slices"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"

	pb "globesort/sortlog"
)

type Record struct {
	key   []byte
	value []byte
}

type NodeConfig struct {
	NodeID int    `yaml:"nodeID"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

type GlobeSortConfig struct {
	Nodes []NodeConfig `yaml:"nodes"`
}

type Server struct {
	pb.UnimplementedGlobeSortServer
	ping    chan int
	channel chan *pb.Record // channel to collect ready records
	cache   []*pb.Record    // records of each node's portion
	done    chan int        // EXIT after all records are processed
}

func (serve *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	pairID, err := strconv.Atoi(req.NodeId)
	if err != nil {
		log.Printf("Invalid serverID: %v", err)
	}
	serve.ping <- pairID

	return &pb.PingResponse{Message: "pong"}, nil
}

func (serve *Server) SendRecord(ctx context.Context, rec *pb.Record) (*pb.Ack, error) {
	serve.channel <- rec
	return &pb.Ack{Message: "received"}, nil
}

func (serve *Server) Close(ctx context.Context, _ *pb.Empty) (*pb.Ack, error) {
	select {
	case serve.done <- 1:
	default:
		fmt.Println("Nothing ready")
	}
	return &pb.Ack{Message: "closed"}, nil
}

// Decode .yaml configuration file
func readConfig(configFilePath string) (*GlobeSortConfig, error) {
	config := GlobeSortConfig{}
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// All below inspired from Lab 1:
// ****************************** START ******************************

// Write a big-endian uint32 to a byte slice of length at least 4.
func ReadBigEndianUint32(buffer []byte) uint32 {
	if len(buffer) < 4 {
		panic("buffer too short to read uint32")
	}
	return binary.BigEndian.Uint32(buffer)
}

// Write a big-endian uint32 to a byte slice of length at least 4.
func WriteBigEndianUint32(buffer []byte, num uint32) {
	if len(buffer) < 4 {
		panic("buffer too short to write uint32")
	}
	binary.BigEndian.PutUint32(buffer, num)
}

// Split every Record struct.
func SplitRecords(data []byte) []*pb.Record {
	var records []*pb.Record
	offset := 0

	for offset < len(data) {
		// Split length
		lengths := ReadBigEndianUint32(data[offset : offset+4])
		// validate the length
		if lengths < 10 {
			fmt.Printf("invalid length")
		}

		offset += 4
		//validate the completeness of record
		endOfRecord := offset + int(lengths)
		if endOfRecord > len(data) {
			fmt.Printf("incomplete record")
		}

		// Split key
		key := data[offset : offset+10]
		offset += 10

		// Split value
		value := data[offset:endOfRecord]

		records = append(records, &pb.Record{Key: key, Value: value})
		offset = endOfRecord
	}
	return records
}

// Sort Record by Custom Function
func CompareRecords(a, b Record) int {
	return bytes.Compare(a.key[:], b.key[:])
}

// ****************************** END ******************************

// Calculate the length of ID bits
func idLength(n_nodes uint) int {
	return bits.Len(n_nodes) - 1
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)

	if len(os.Args) != 5 {
		fmt.Println("Usage:", os.Args[0], "<nodeID> <inputFilePath> <outputFilePath> <configFilePath>")
		os.Exit(1)
	}

	serverId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Server ID must be an integer! Got '%s'", os.Args[1])
	}

	inputFilePath := os.Args[2]
	outputFilePath := os.Args[3]
	configFilePath := os.Args[4]

	log.Printf("serverID: %d", serverId)
	log.Printf("inputFilePath: %s", inputFilePath)
	log.Printf("outputFilePath: %s", outputFilePath)
	log.Printf("configFilePath: %s", configFilePath)

	config, err := readConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	log.Printf("Configured nodes: %+v", config.Nodes)

	var node NodeConfig
	for _, n := range config.Nodes {
		if n.NodeID == serverId {
			node = n
			break
		}
	}
	addr := node.Host + ":" + strconv.Itoa(node.Port)

	listener, err := net.Listen("tcp", addr)
	// log.Printf("Node %d listening on %s", node.NodeID, addr)
	if err != nil {
		log.Fatalf("Failed to listen to client: %v", err)
	}
	defer listener.Close()

	log.Printf("Server started on %s, serving files from %s", addr, inputFilePath)

	// ****************************** Server Side ******************************
	grpcServer := grpc.NewServer()
	serve := &Server{
		ping:    make(chan int, len(config.Nodes)-1),
		channel: make(chan *pb.Record, 1024),
		cache:   make([]*pb.Record, 0),
		done:    make(chan int, len(config.Nodes)-1),
	}

	done := make(chan int)
	go func() {
		for rec := range serve.channel {
			serve.cache = append(serve.cache, rec)
		}
		close(done) // signal that consumption is complete, specialize for 1 NODE case
	}()

	pb.RegisterGlobeSortServer(grpcServer, serve)
	go func() {
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("gRPC Server Exit: %v", err)
		}
	}()
	time.Sleep(time.Second)

	// ****************************** Client Side ******************************
	clients := make(map[int]pb.GlobeSortClient) // Map Clients to NodeID
	for _, pair := range config.Nodes {         //Connect to other n-1 Nodes
		if pair.NodeID == serverId {
			continue
		}
		dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Second) // 2s timeout
		defer dialCancel()
		pairAddr := pair.Host + ":" + strconv.Itoa(pair.Port)
		// conn, err := grpc.Dial(pairAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn, err := grpc.DialContext(
			dialCtx,
			pairAddr,
			grpc.WithBlock(), // wait for connection
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to dial node %d: %v", pair.NodeID, err)
		}
		client := pb.NewGlobeSortClient(conn)
		clients[pair.NodeID] = client

		pingCtx, pingCancel := context.WithTimeout(context.Background(), time.Second)
		defer pingCancel()

		// Ping
		_, err = client.Ping(pingCtx, &pb.PingRequest{
			NodeId: strconv.Itoa(serverId),
		})
		if err != nil {
			log.Fatalf("Ping %s Failed: %v", pairAddr, err)
		}

	}
	// Wait until all Ping complete
	for i := 0; i < len(config.Nodes)-1; i++ {
		<-serve.ping
	}

	// Read its own designated input data
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		fmt.Printf("Read file error: %v\n", err)
	}
	records := SplitRecords(data)
	// Send Records
	for _, rec := range records {
		length := idLength(uint(len(config.Nodes)))
		target := int(rec.Key[0] >> (8 - length)) // right shift to obtain the value of first n bits
		if target == serverId {                   // Solve Concurrency
			serve.channel <- rec
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			_, err := clients[target].SendRecord(ctx, rec)
			if err != nil {
				log.Fatalf("SendRecord to node %d failed: %v", target, err)
			}
			cancel()
		}
	}
	// Close all the clients after all records processed
	for nodeID, client := range clients {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := client.Close(ctx, &pb.Empty{})
		if err != nil {
			log.Fatalf("Close to node %d failed: %v", nodeID, err)
		}
		cancel()
	}

	// Wait until all Closed
	for i := 0; i < len(config.Nodes)-1; i++ {
		<-serve.done
		fmt.Printf("Received close signal #%d\n", i+1)
	}
	close(serve.channel)
	<-done

	// ****************************** Sort and Output the File ******************************

	record := make([]Record, len(serve.cache))
	for i, r := range serve.cache {
		// copy key and value into my Record
		record[i] = Record{
			key:   r.Key,
			value: r.Value,
		}
	}

	//sort records
	slices.SortFunc(record, CompareRecords)

	output, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Create file error %v\n", err)
	}
	defer output.Close()

	for _, rec := range record {
		// write length-prefixed key+value (as in Lab1)
		buffer := make([]byte, 4)
		WriteBigEndianUint32(buffer, uint32(len(rec.key)+len(rec.value)))

		// write length to output file
		_, err := output.Write(buffer)
		if err != nil {
			fmt.Printf("Write Length error %v\n", err)
		}

		// write key to output file
		_, err = output.Write(rec.key)
		if err != nil {
			fmt.Printf("Write Prefix-Key error %v\n", err)
		}

		// write value to output file
		_, err = output.Write(rec.value)
		if err != nil {
			fmt.Printf("Write Value error %v\n", err)
		}
	}

	// ****************************** THE END ******************************
	grpcServer.GracefulStop()

}
