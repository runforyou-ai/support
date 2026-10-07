package filex

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// WriteAtomic writes data to path so that readers see either the old content
// or the new content, never a partial write: it writes a temporary file in the
// same directory, syncs it to disk, renames it over path and syncs the
// directory where the platform supports it. An existing
// regular file keeps its permission bits; a new file gets perm exactly, not
// reduced by the umask. A symbolic link at path is replaced by a regular file
// rather than followed, and any other existing non-regular file, such as a
// directory, is an error. The parent directory must exist. On failure path is
// left unchanged and the temporary file is removed.
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
	// Syncing the directory persists the rename; platforms without directory sync ignore it.
	if dir, err := os.Open(filepath.Dir(path)); err == nil {
		_ = dir.Sync()
		_ = dir.Close()
	}
	return nil
}
