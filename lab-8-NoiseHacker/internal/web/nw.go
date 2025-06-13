// Lab 8: Implement a network video content service (client using consistent hashing)

package web

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"net"
	"slices"
	"tritontube/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NetworkVideoContentService implements VideoContentService using a network of nodes.
type NetworkVideoContentService struct{
	adminAddr string
	storageAddrs []string
	nodeHashes []uint64
	nodeHashToAddr map[uint64]string
	clients map[string]proto.StorageServiceClient
	proto.UnimplementedVideoContentAdminServiceServer
}

// Uncomment the following line to ensure NetworkVideoContentService implements VideoContentService
var _ VideoContentService = (*NetworkVideoContentService)(nil)


// ******************** 1. NEW network content service ********************
// Initializes the network content service with a hash ring and gRPC clients in web/main.go
func NewNetworkVideoContentService(adminAddr string, storageAddrs []string) (*NetworkVideoContentService, error) {
	nodeHashes := make([]uint64, 0, len(storageAddrs))
	nodeHashToAddr := make(map[uint64]string)
	clients := make(map[string]proto.StorageServiceClient)

	for _, nodeAddr := range storageAddrs {
		hash := hashStringToUint64(nodeAddr)
		nodeHashes = append(nodeHashes, hash)
		nodeHashToAddr[hash] = nodeAddr

		connection, err := grpc.Dial(nodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to storage node %s: %w", nodeAddr, err)
		}
		clients[nodeAddr] = proto.NewStorageServiceClient(connection)
	}

	slices.Sort(nodeHashes)

	service := &NetworkVideoContentService{
		adminAddr:      adminAddr,
		storageAddrs:   storageAddrs,
		nodeHashes:     nodeHashes,
		nodeHashToAddr: nodeHashToAddr,
		clients:        clients,
	}
	return service, nil
}

// ******************** 2. Read and Write ******************************
// Retrieves a content file from the correct storage node using consistent hashing
func (n *NetworkVideoContentService) Read(videoID string, filename string) ([]byte, error) {
	key := videoID + "/" + filename
	nodeAddr := n.getNodeForKey(key)
	client := n.getStorageClient(nodeAddr)

	req := &proto.ReadFileRequest{Key: key}
	response, _ := client.ReadFile(context.Background(), req)
	return response.Data, nil
}

// Stores a content file to the correct storage node using consistent hashing
func (n *NetworkVideoContentService) Write(videoID string, filename string, data []byte) error {
	key := videoID + "/" + filename
	nodeAddr := n.getNodeForKey(key)
	client := n.getStorageClient(nodeAddr)

	req := &proto.WriteFileRequest{
		Key:  key,
		Data: data,
	}
	response, _ := client.WriteFile(context.Background(), req)
	if !response.Success {
		return fmt.Errorf("storage node %s failed to write key %s", nodeAddr, key)
	}
	fmt.Println("Writing to", nodeAddr, "key =", key)
	return nil
}


// ******************** 3. Starts the gRPC server ********************
// Starts the gRPC server to handle admin CLI requests
func (n *NetworkVideoContentService) StartAdminGRPCServer() {
	listener, err := net.Listen("tcp", n.adminAddr)
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}
	grpcServer := grpc.NewServer()
	proto.RegisterVideoContentAdminServiceServer(grpcServer, n)
	grpcServer.Serve(listener)
}


// ********** 4. Implement Node Operations Specified in admin.proto **********
// Returns the list of storage nodes in the hash ring in sorted order
func (n *NetworkVideoContentService) ListNodes(ctx context.Context, req *proto.ListNodesRequest) (*proto.ListNodesResponse, error) {
	var nodeAddrs []string
	for _, nodeHash := range n.nodeHashes {
		nodeAddrs = append(nodeAddrs, n.nodeHashToAddr[nodeHash])
	}
	response := &proto.ListNodesResponse{Nodes: nodeAddrs}
	return response, nil
}
// Adds a new node to the cluster and migrates affected files to the new node
func (n *NetworkVideoContentService) AddNode(ctx context.Context, req *proto.AddNodeRequest) (*proto.AddNodeResponse, error) {
	nodeAddr := req.NodeAddress
	hash := hashStringToUint64(nodeAddr)
	n.nodeHashes = append(n.nodeHashes, hash)
	n.nodeHashToAddr[hash] = nodeAddr
	slices.Sort(n.nodeHashes)

	conn, _ := grpc.Dial(nodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	n.clients[nodeAddr] = proto.NewStorageServiceClient(conn)

	migratedCount := 0
	for _, h := range n.nodeHashes {
		addr := n.nodeHashToAddr[h]
		if addr == nodeAddr {
			continue
		}
		keys, err := n.getAllKeysFromNode(addr)
		if err != nil {
			continue
		}
		for _, key := range keys {
			target := n.getNodeForKey(key)
			if target == nodeAddr {
				err := n.migrateFile(key, addr, nodeAddr)
				if err == nil {
					migratedCount++
				}
			}
		}
	}

	response := &proto.AddNodeResponse{MigratedFileCount: int32(migratedCount)}
	return response, nil
}
// Removes a node from the cluster and migrates its files to remaining nodes
func (n *NetworkVideoContentService) RemoveNode(ctx context.Context, req *proto.RemoveNodeRequest) (*proto.RemoveNodeResponse, error) {
	nodeAddr := req.NodeAddress
	nodeHash := hashStringToUint64(nodeAddr)

	delete(n.nodeHashToAddr, nodeHash)
	for i, h := range n.nodeHashes {
		if h == nodeHash {
			n.nodeHashes = append(n.nodeHashes[:i], n.nodeHashes[i+1:]...)
			break
		}
	}

	keys, _ := n.getAllKeysFromNode(nodeAddr)

	migratedCount := 0
	for _, key := range keys {
		newNodeAddr := n.getNodeForKey(key)
		if newNodeAddr != nodeAddr {
			err := n.migrateFile(key, nodeAddr, newNodeAddr)
			if err == nil {
				migratedCount++
			}
		}
	}
	delete(n.clients, nodeAddr)

	response := &proto.RemoveNodeResponse{MigratedFileCount: int32(migratedCount)}
	return response, nil
}

// ******************** 5. Consistent Hashing Logic ********************
func hashStringToUint64(s string) uint64 {
	sum := sha256.Sum256([]byte(s))
	return binary.BigEndian.Uint64(sum[:8])
}
// Returns the storage node responsible for the given key
func (n *NetworkVideoContentService) getNodeForKey(key string) string {
	keyHash := hashStringToUint64(key)
	for _, nodeHash := range n.nodeHashes {
		if keyHash <= nodeHash {
			return n.nodeHashToAddr[nodeHash]
		}
	}
	// if larger than all nodes, map to the smallest node since its a ring
	return n.nodeHashToAddr[n.nodeHashes[0]]
}
// Returns the gRPC client for a given node address
func (n *NetworkVideoContentService) getStorageClient(nodeAddr string) proto.StorageServiceClient {
	client, ok := n.clients[nodeAddr]
	if !ok {
		panic("client not found for nodeAddr: " + nodeAddr)
	}
	return client
}
// Copies a file from one node to another and optionally deletes the original
func (n *NetworkVideoContentService) migrateFile(key string, fromAddr string, toAddr string) error {
	fromClient := n.getStorageClient(fromAddr)
	readRequest := &proto.ReadFileRequest{Key: key}
	readResponse, err := fromClient.ReadFile(context.Background(), readRequest)
	if err != nil {
		return err
	}

	toClient := n.getStorageClient(toAddr)
	writeRequest := &proto.WriteFileRequest{Key: key, Data: readResponse.Data}
	writeResponse, err := toClient.WriteFile(context.Background(), writeRequest)
	if err != nil || !writeResponse.Success {
		return fmt.Errorf("failed to write to node %s", toAddr)
	}

	deleteRequest := &proto.DeleteFileRequest{Key: key}
	_, err = fromClient.DeleteFile(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}
// Retrieves all content keys stored at a given node (used during migration)
func (n *NetworkVideoContentService) getAllKeysFromNode(nodeAddr string) ([]string, error) {
	client := n.getStorageClient(nodeAddr)
	request := &proto.ListKeysRequest{}
	response, err := client.ListKeys(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response.Keys, nil
}