// Lab 7: Implement a local filesystem video content service

package web

import (
	"fmt"
	"os"
	"path/filepath"
)

// FSVideoContentService implements VideoContentService using the local filesystem.
type FSVideoContentService struct{rootDir string}

// Uncomment the following line to ensure FSVideoContentService implements VideoContentService
var _ VideoContentService = (*FSVideoContentService)(nil)

// Constructor
func NewFSVideoContentService(root string) *FSVideoContentService {
	return &FSVideoContentService{rootDir: root}
}

// Write saves data under {root}/{videoId}/{filename}
func (fs *FSVideoContentService) Write(videoId string, filename string, data []byte) error {
	if videoId == "" {
    	panic("videoID cannot be empty!")
	}
	
	dirPath := filepath.Join(fs.rootDir, videoId)
	err := os.MkdirAll(dirPath, 0777)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	fullPath := filepath.Join(dirPath, filename)
	return os.WriteFile(fullPath, data, 0777)
}

// Read reads file from {root}/{videoId}/{filename}
func (fs *FSVideoContentService) Read(videoId string, filename string) ([]byte, error) {
	fullPath := filepath.Join(fs.rootDir, videoId, filename)
	return os.ReadFile(fullPath)
}