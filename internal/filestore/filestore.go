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
	Init() error
	GetFile(m *FileStorerObjectMetadata) (*io.ReadCloser, error)
	WriteFile(src io.ReadCloser) (*FileStorerObjectMetadata, error)
}

type LocalFileStorer struct {
	path string
}

// ensure LocalFileStorer implements FileStorer interface
var _ FileStorer = (*LocalFileStorer)(nil)

func New(path string) *LocalFileStorer {
	return new(LocalFileStorer{path: path})
}

func (l *LocalFileStorer) Init() error {
	return os.MkdirAll(l.path, 0o777)
}

func (l *LocalFileStorer) GetFile(o *FileStorerObjectMetadata) (*io.ReadCloser, error) {
	panic("not implemented yet")
}

func (l *LocalFileStorer) WriteFile(src io.ReadCloser) (*FileStorerObjectMetadata, error) {
	panic("not implemented yet")
}
