package mediaprocessor

import (
	"image"
	_ "image/gif" // register decoders as a side effect
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"mime"
	"strings"

	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	maxPixels   = 100_000_000
	maxFileSize = 2 << 48
)

type MediaResult interface{}

type MediaUploadResult struct {
	VideoResult *VideoMetadata
	ImageResult *ImageMetadata
	ContentType string
}

func ProcessMediaUpload(ctx context.Context, file *multipart.FileHeader, cachePath string, targetFilename string) (*MediaUploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buf := make([]byte, 512)
	n, err := io.ReadFull(src, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	result := new(MediaUploadResult{})

	src.Seek(0, io.SeekStart)
	contentType := http.DetectContentType(buf[:n])
	slog.Debug("process", "content_type", contentType)
	switch {
	case isImageFile(contentType):
		result.ImageResult = new(ImageMetadata{})
		err = getImageMetadata(result.ImageResult, src)
		result.ImageResult.OriginalFilename = filepath.Base(file.Filename)
		if err != nil {
			return nil, err
		}
	case isVideoFile(contentType):
		result.VideoResult = new(VideoMetadata{})
		result.VideoResult.OriginalFilename = filepath.Base(file.Filename)
		// TODO
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported content type: %q", contentType)
	}

	result.ContentType = contentType
	fileExt := ".bin"
	possibleFileExtensions, err := mime.ExtensionsByType(contentType)
	if possibleFileExtensions != nil {
		fileExt = possibleFileExtensions[0]
	}

	targetPath := filepath.Join(cachePath, targetFilename+fileExt)

	dst, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := dst.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			os.Remove(targetPath) // no orphaned files on any failure path
		}
	}()

	src.Seek(0, io.SeekStart)
	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getImageMetadata(imageMetadata *ImageMetadata, src multipart.File) error {
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return err
	}
	cfg, _, err := image.DecodeConfig(src)
	if err != nil {
		return fmt.Errorf("invalid image: %w", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width*cfg.Height > maxPixels {
		return fmt.Errorf("unsupported image dimensions: %dx%d", cfg.Width, cfg.Height)
	}

	imageMetadata.Width = cfg.Width
	imageMetadata.Height = cfg.Height

	return nil
}

// very crude lol
func isImageFile(fileType string) bool {
	return strings.HasPrefix(fileType, "image/")
}
func isVideoFile(fileType string) bool {
	return strings.HasPrefix(fileType, "video/")
}
