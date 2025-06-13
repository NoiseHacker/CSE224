// Lab 8: Implement a network video content service (server)

package storage

// Implement a network video content service (server)
import (
	"context"
	"os"
	"path/filepath"
	"tritontube/internal/proto"
)

type StorageHandler struct {
	baseDir string
	proto.UnimplementedStorageServiceServer
}

func NewStorageHandler(baseDir string) *StorageHandler {
	return &StorageHandler{baseDir: baseDir}
}

func (h *StorageHandler) WriteFile(ctx context.Context, req *proto.WriteFileRequest) (*proto.WriteFileResponse, error) {
	path := filepath.Join(h.baseDir, req.Key)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &proto.WriteFileResponse{Success: false}, err
	}
	err := os.WriteFile(path, req.Data, 0644)
	return &proto.WriteFileResponse{Success: err == nil}, err
}

func (h *StorageHandler) ReadFile(ctx context.Context, req *proto.ReadFileRequest) (*proto.ReadFileResponse, error) {
	path := filepath.Join(h.baseDir, req.Key)
	data, err := os.ReadFile(path)
	return &proto.ReadFileResponse{Data: data}, err
}

func (h *StorageHandler) DeleteFile(ctx context.Context, req *proto.DeleteFileRequest) (*proto.DeleteFileResponse, error) {
	path := filepath.Join(h.baseDir, req.Key)
	err := os.Remove(path)
	return &proto.DeleteFileResponse{Success: err == nil}, err
}

func (h *StorageHandler) ListKeys(ctx context.Context, request *proto.ListKeysRequest) (*proto.ListKeysResponse, error) {
	var keys []string

	entries, err := os.ReadDir(h.baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		videoID := entry.Name()
		videoDir := filepath.Join(h.baseDir, videoID)

		files, err := os.ReadDir(videoDir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			filename := file.Name()
			key := videoID + "/" + filename
			keys = append(keys, key)
		}
	}

	return &proto.ListKeysResponse{Keys: keys}, nil
}