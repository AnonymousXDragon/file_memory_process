package file

import (
	"errors"
	// "fmt"
	// "log"
	"path/filepath"
	"strings"
	"sync"
)

type File struct {
	Name     string
	Content  []byte
	IsDir    bool
	Children map[string]*File
}

type FileSystem struct {
	Root *File
	mu   sync.RWMutex
}

func NewFileSystem() *FileSystem {
	return &FileSystem{
		Root: &File{
			Name:     "/",
			IsDir:    true,
			Children: make(map[string]*File),
		},
	}
}

func (fs *FileSystem) CreateFile(path string, content []byte) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dir, name := filepath.Split(path)
	parent, err := fs.findDir(dir)
	if err != nil {
		return err
	}

	parent.Children[name] = &File{Name: name, Content: content}
	return nil
}

func (fs *FileSystem) ReadFile(path string) ([]byte, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	f, err := fs.findFile(path)

	if err != nil {
		return nil, err
	}

	return f.Content, nil
}

func (fs *FileSystem) findFile(path string) (*File, error) {
	parts := strings.Split(filepath.Clean(path), string(filepath.Separator))
	current := fs.Root

	if string(filepath.Separator) == path {
		return current, nil
	}

	for _, part := range parts[1:] {
		if !current.IsDir {
			return nil, errors.New("not a directory")
		}

		file, ok := current.Children[part]
		if !ok {
			return nil, errors.New("file not found")
		}
		current = file
	}

	return current, nil
}

func (fs *FileSystem) findDir(path string) (*File, error) {
	dir, err := fs.findFile(path)

	if err != nil {
		return nil, err
	}

	if !dir.IsDir {
		return nil, errors.New("not a directory")
	}

	return dir, nil
}
