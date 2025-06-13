// Lab 7: Implement a web server

package web

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type server struct {
	Addr string
	Port int

	metadataService VideoMetadataService
	contentService  VideoContentService

	mux *http.ServeMux
}

type VideoInfo struct {
	Id        string // original video ID
	EscapedId string // path-save video ID
	UploadTime  string // human-readable upload time
}

func NewServer(
	metadataService VideoMetadataService,
	contentService VideoContentService,
) *server {
	return &server{
		metadataService: metadataService,
		contentService:  contentService,
	}
}

func (s *server) Start(lis net.Listener) error {
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/upload", s.handleUpload)
	s.mux.HandleFunc("/videos/", s.handleVideo)
	s.mux.HandleFunc("/content/", s.handleVideoContent)
	s.mux.HandleFunc("/", s.handleIndex)

	return http.Serve(lis, s.mux)
}

// Handles the "/" endpoint
func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	videos, err := s.metadataService.List()
	if err != nil {
		http.Error(w, "Failed to list videos", http.StatusInternalServerError)
		return
	}

	var wrapped []VideoInfo
	for _, v := range videos {
		wrapped = append(wrapped, VideoInfo{
			Id:        v.Id,
			EscapedId: url.PathEscape(v.Id),
			UploadTime:  v.UploadedAt.Format("2006-01-02 15:04:05"),
		})
	}

	temp, err := template.New("index").Parse(indexHTML) // indexHTML 是 string
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	err = temp.Execute(w, wrapped)
	if err != nil {
		log.Println("Template rendering failed:", err)
	}
}

// Handles the "/upload" endpoint
func (s *server) handleUpload(w http.ResponseWriter, r *http.Request) {
	// Encode with max 100MB memory
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File field does not exist", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := header.Filename
	videoID := strings.TrimSuffix(filename, filepath.Ext(filename))

	existing, err := s.metadataService.Read(videoID)
	if err == nil && existing != nil {
		http.Error(w, "Video ID already exists", http.StatusConflict)
		return
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Failed to check existing video", http.StatusInternalServerError)
		return
	}

	tempDir, err := os.MkdirTemp("", "upload-*")
	if err != nil {
		http.Error(w, "Failed to create temp dir", http.StatusInternalServerError)
		return
	}

	dstPath := filepath.Join(tempDir, filename)
	dstFile, err := os.Create(dstPath)
	// fmt.Println("Full path" + dstPath)
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, file)
	if err != nil {
		http.Error(w, "Failed to save uploaded file", http.StatusInternalServerError)
		return
	}

	err = runFFmpeg(dstPath, tempDir)
	if err != nil {
		http.Error(w, "Video converting failed", http.StatusInternalServerError)
		return
	}

	err = StoreInContentService(s.contentService, videoID, tempDir)
	if err != nil {
		http.Error(w, "Failed to save video content", http.StatusInternalServerError)
		return
	}

	// data, err := os.ReadFile(dstPath)
	// if err != nil {
	// 	http.Error(w, "Failed to read temp file", http.StatusInternalServerError)
	// 	return
	// }

	// err = s.contentService.Write(videoID, filename, data)
	// if err != nil {
	// 	http.Error(w, "Failed to write content", http.StatusInternalServerError)
	// 	return
	// }

	err = s.metadataService.Create(videoID, time.Now())
	if err != nil {
		http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
		return
	}

	// Redirect to "/"
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handles the "/videos/:videoId" endpoint
func (s *server) handleVideo(w http.ResponseWriter, r *http.Request) {
	videoId := r.URL.Path[len("/videos/"):]
	log.Println("Video ID:", videoId)
	// Lookup metadata
	meta, err := s.metadataService.Read(videoId)
	if err != nil || meta == nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	info := VideoInfo{
		Id:        meta.Id,
		EscapedId: url.PathEscape(meta.Id),
		UploadTime:  meta.UploadedAt.Format("2006-01-02 15:04:05"),
	}

	temp, err := template.New("video").Parse(videoHTML)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	temp.Execute(w, info)	
}

// Handles the "/content/:videoId/:filename" endpoint
func (s *server) handleVideoContent(w http.ResponseWriter, r *http.Request) {
	// parse /content/<videoId>/<filename>
	videoId := r.URL.Path[len("/content/"):]
	parts := strings.Split(videoId, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid content path", http.StatusBadRequest)
		return
	}
	videoId = parts[0]
	filename := parts[1]
	// log.Println("Video ID:", videoId, "Filename:", filename)

	data, err := s.contentService.Read(videoId, filename)
	if err != nil {
		http.Error(w, "Failed to read content", http.StatusInternalServerError)
		return
	}

	if strings.HasSuffix(filename, ".mpd") {
		w.Header().Set("Content-Type", "application/dash+xml")
	} else if strings.HasSuffix(filename, ".m4s"){
		w.Header().Set("Content-Type", "video/mp4")
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Converts MP4 to MPEG-DASH format using ffmpeg
func runFFmpeg(videoPath, tempDir string) error {
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,                          // input file
		"-c:v", "libx264",                        // video codec
		"-c:a", "aac",                            // audio codec
		"-bf", "1",                               // max 1 b-frame
		"-keyint_min", "30",                     // minimum keyframe interval
		"-g", "30",                              // keyframe every 25 frames
		"-sc_threshold", "0",                     // scene change threshold
		"-b:v", "3000k",                          // video bitrate
		"-b:a", "128k",                           // audio bitrate
		"-f", "dash",                             // dash format
		"-use_timeline", "1",                     // use timeline
		"-use_template", "1",                     // use template
		"-init_seg_name", "init-$RepresentationID$.m4s",       // init segment naming
		"-media_seg_name", "chunk-$RepresentationID$-$Number%05d$.m4s", // media segment naming
		"-seg_duration", "2",                     // segment duration in seconds
		"-min_seg_duration", "2000000",				// minimum segment duration in us
		"manifest.mpd",                             // output file
	)
	cmd.Dir = tempDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func StoreInContentService(fs VideoContentService, videoID, dir string) error {
	fmt.Println("Inspecting tempDir:", dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to list dir: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()

		if strings.HasSuffix(name, ".mp4") {
			continue
		}

		path := filepath.Join(dir, name)
		// Read back the file content for writing to contentService
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}

		// Save to content service
		err = fs.Write(videoID, name, data)
		if err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	// time.Sleep(5 * time.Second)
	return nil
}
