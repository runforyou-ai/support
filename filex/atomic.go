package filex

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

// WriteAtomic writes data to path so that readers see either the old content
// or the new content, never a partial write: it writes a temporary file in the
// same directory, syncs it to disk, renames it over path and then syncs the
// parent directory so the rename survives a crash. Directory sync is skipped
// on Windows and where the file system does not support it. An existing
// regular file keeps its permission bits; a new file gets perm exactly, not
// reduced by the umask. A symbolic link at path is replaced by a regular file
// rather than followed, and any other existing non-regular file, such as a
// directory, is an error. The parent directory must exist. When an error
// occurs before the rename, path is left unchanged and the temporary file is
// removed; when syncing the directory fails, path already holds the new
// content and the error is returned.
func WriteAtomic(path string, data []byte, perm os.FileMode) error {
	mode := perm.Perm()
	info, err := os.Lstat(path)
	switch {
	case err == nil && info.Mode().IsRegular():
		mode = info.Mode().Perm()
	case err == nil && info.Mode()&os.ModeSymlink == 0:
		return fmt.Errorf("filex: write %s: not a regular file", path)
	case err != nil && !errors.Is(err, os.ErrNotExist):
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(temp.Name()) }()
	_, writeErr := temp.Write(data)
	if err := errors.Join(writeErr, temp.Sync(), temp.Close()); err != nil {
		return err
	}
	if err := os.Chmod(temp.Name(), mode); err != nil {
		return err
	}
	if err := os.Rename(temp.Name(), path); err != nil {
		return err
	}
	return syncDir(filepath.Dir(path))
}

// syncDir flushes the directory entry changes of dir to disk.
func syncDir(dir string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	file, err := os.Open(dir)
	if err != nil {
		return err
	}
	err = file.Sync()
	if errors.Is(err, errors.ErrUnsupported) || errors.Is(err, syscall.EINVAL) || errors.Is(err, syscall.ENOTSUP) {
		err = nil
	}
	return errors.Join(err, file.Close())
}
