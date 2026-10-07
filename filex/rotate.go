package filex

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// RotatingFile is an append-only io.WriteCloser that rotates its file by size:
// once a write takes the file past the size limit, the file is renamed to
// <path>.1, existing backups shift to <path>.2, <path>.3 and so on, backups
// beyond the configured count are removed, and a new empty file is started.
// It is safe for concurrent use. Create it with OpenRotating.
type RotatingFile struct {
	path    string
	maxSize int64
	backups int
	mu      sync.Mutex
	file    *os.File
	size    int64
}

// OpenRotating opens path for appending, creating it with mode 0644 and its
// parent directories with mode 0755 when missing. The file is rotated after
// any write that takes it past maxSize bytes, and at most backups rotated
// files are kept; with backups zero or negative the file is simply truncated
// on rotation. A single write is never split, so a file may exceed maxSize by
// up to one write.
func OpenRotating(path string, maxSize int64, backups int) (*RotatingFile, error) {
	f := &RotatingFile{path: path, maxSize: maxSize, backups: max(backups, 0)}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("filex: create directory: %w", err)
	}
	if err := f.open(); err != nil {
		return nil, err
	}
	return f, nil
}

// Write appends data to the current file and rotates the file when it now
// exceeds the size limit. It returns os.ErrClosed after Close. A rotation
// failure is reported together with the number of bytes already written.
func (f *RotatingFile) Write(data []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.file == nil {
		return 0, os.ErrClosed
	}
	written, err := f.file.Write(data)
	f.size += int64(written)
	if err != nil {
		return written, err
	}
	if f.size > f.maxSize {
		if err := f.rotate(); err != nil {
			return written, err
		}
	}
	return written, nil
}

// Close closes the current file. Closing an already closed RotatingFile
// returns nil.
func (f *RotatingFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.file == nil {
		return nil
	}
	err := f.file.Close()
	f.file = nil
	return err
}

// open opens the current file for appending and records its size.
func (f *RotatingFile) open() error {
	file, err := os.OpenFile(f.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("filex: open rotating file: %w", err)
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("filex: open rotating file: %w", err)
	}
	f.file, f.size = file, info.Size()
	return nil
}

// rotate closes the current file, shifts the backups, removes the one beyond
// the retention count and starts a new current file.
func (f *RotatingFile) rotate() error {
	if err := f.file.Close(); err != nil {
		f.file = nil
		return fmt.Errorf("filex: close rotating file: %w", err)
	}
	f.file = nil
	if f.backups > 0 {
		if err := os.Remove(fmt.Sprintf("%s.%d", f.path, f.backups)); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("filex: rotate file: %w", err)
		}
		for index := f.backups - 1; index >= 1; index-- {
			err := os.Rename(fmt.Sprintf("%s.%d", f.path, index), fmt.Sprintf("%s.%d", f.path, index+1))
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("filex: rotate file: %w", err)
			}
		}
		if err := os.Rename(f.path, f.path+".1"); err != nil {
			return fmt.Errorf("filex: rotate file: %w", err)
		}
	} else if err := os.Remove(f.path); err != nil {
		return fmt.Errorf("filex: rotate file: %w", err)
	}
	return f.open()
}
