package filestore

import (
	"io"
	"os"
	"uuid"
)

type FileStorerObjectMetadata struct {
	ID        uuid.UUID
	Extension string
	Filetype  string
}

type FileStorer interface {
	GetFile(m *FileStorerObjectMetadata) (*io.ReadCloser, error)
	WriteFile(src io.ReadCloser) (*FileStorerObjectMetadata, error)
}

type LocalFileStorer struct {
	path string
}

// ensure LocalFileStorer implements FileStorer interface
var _ FileStorer = (*LocalFileStorer)(nil)

func NewLocalFileStorer(path string, perms os.FileMode) (FileStorer, error) {
	err := os.MkdirAll(path, perms)
	if err != nil {
		return nil, err
	}

	return new(LocalFileStorer{path: path}), nil
}

func (l LocalFileStorer) GetFile(o *FileStorerObjectMetadata) (*io.ReadCloser, error) {
	panic("not implemented yet")
}

func (l LocalFileStorer) WriteFile(src io.ReadCloser) (*FileStorerObjectMetadata, error) {
	panic("not implemented yet")
}
