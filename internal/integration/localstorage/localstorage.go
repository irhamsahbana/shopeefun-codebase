package integration

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type LocalStorageContract interface {
	Save(base64String, path string) (fullpath string, err error)
}

var (
	ErrFileTypeNotSupported = errors.New("file type not supported")
	ErrDecodeBase64         = errors.New("failed to decode base64 string")
)

type localstorage struct {
	mu sync.Mutex
}

func NewLocalStorageIntegration() LocalStorageContract {
	return &localstorage{}
}

func (l *localstorage) Save(base64String, path string) (fullpath string, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Trim base64 prefix if present (e.g., "data:image/png;base64,")
	if idx := strings.Index(base64String, ","); idx != -1 {
		base64String = base64String[idx+1:]
	}

	// Decode base64 string to byte slice
	fileContent, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		log.Error().Err(err).Msg("localstorage: failed to decode base64 string")
		return "", fmt.Errorf("localstorage: %w", err)
	}

	// Get MIME type of the file
	mimeType := l.getMimeType(fileContent)
	if !l.isAcceptableMimeType(mimeType) {
		log.Error().Str("mimeType", mimeType).Msg("localstorage: file type not supported")
		return "", ErrFileTypeNotSupported
	}
	ext := l.extensionFromMimeType(mimeType)

	// Get file extension from MIME type
	filename := fmt.Sprintf("%s.%s", ulid.Make().String(), ext)

	// Save file to local storage
	fullpath = fmt.Sprintf("%s/%s", path, filename)
	if err := l.saveFile(fullpath, fileContent); err != nil {
		return "", err
	}

	return fullpath, nil
}

func (l *localstorage) saveFile(fullpath string, data []byte) error {
	path := strings.Split(fullpath, "/")         // Split path by "/"
	dir := strings.Join(path[:len(path)-1], "/") // Join path except the last element

	err := os.MkdirAll(dir, os.ModePerm) // Create directory if not exists
	if err != nil {
		log.Error().Err(err).Msg("localstorage: failed to create directory")
		return fmt.Errorf("localstorage: %w", err)
	}

	file, err := os.Create(fullpath) // Create file
	if err != nil {
		log.Error().Err(err).Msg("localstorage: failed to create file")
		return fmt.Errorf("localstorage: %w", err)
	}
	defer file.Close() // Close file after function ends

	if _, err := file.Write(data); err != nil { // Write data to file
		log.Error().Err(err).Msg("localstorage: failed to write data to file")
		return fmt.Errorf("localstorage: %w", err)
	}

	return nil
}

func (l *localstorage) getMimeType(data []byte) string {
	mimeType := http.DetectContentType(data)

	return mimeType
}

func (l *localstorage) extensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	default:
		return ""
	}
}

func (l *localstorage) isAcceptableMimeType(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}
