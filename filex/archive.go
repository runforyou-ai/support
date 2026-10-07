package filex

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// EntryPath returns the path at which the archive entry name is extracted
// under root. name uses forward slashes as in tar and zip headers; a leading
// slash is treated as relative to root. It returns an error when the entry
// would land outside root, as with "../evil" (zip slip).
func EntryPath(root, name string) (string, error) {
	path := filepath.Join(root, filepath.FromSlash(name))
	if filepath.VolumeName(filepath.FromSlash(name)) != "" || !Within(root, path) {
		return "", fmt.Errorf("filex: archive entry %q escapes target directory", name)
	}
	return path, nil
}

// ExtractTarGz extracts the gzip-compressed tar archive at archivePath into
// root, creating directories as needed. Only directories, regular files and
// symbolic links are extracted; other entry types, such as hard links and
// devices, are skipped. It returns an error when an entry would escape root,
// when an entry would be written through a symbolic link extracted earlier,
// when a symbolic link target is absolute or lexically outside root, or when
// any extracted link, possibly through a chain of links, resolves outside
// root. Entries extracted before an error are left in place. The total size
// and number of entries are not limited, so callers handling untrusted
// archives should bound the archive size beforehand.
func ExtractTarGz(archivePath, root string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	decompressed, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func() { _ = decompressed.Close() }()
	reader := tar.NewReader(decompressed)
	var links []string
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			// Reading the gzip stream to its end verifies its checksum and length.
			if _, err := io.Copy(io.Discard, decompressed); err != nil {
				return err
			}
			break
		}
		if err != nil {
			return err
		}
		path, err := writablePath(root, header.Name)
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(path, 0o755)
		case tar.TypeReg:
			err = writeFile(path, reader, header.FileInfo().Mode().Perm())
		case tar.TypeSymlink:
			// A symbolic link must have a relative target inside root.
			target := filepath.FromSlash(header.Linkname)
			if filepath.IsAbs(target) || filepath.VolumeName(target) != "" || strings.HasPrefix(target, string(filepath.Separator)) {
				return fmt.Errorf("filex: archive symlink %q is absolute", header.Name)
			}
			if !Within(root, filepath.Join(filepath.Dir(path), target)) {
				return fmt.Errorf("filex: archive symlink %q points outside target directory", header.Name)
			}
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			err = os.Symlink(target, path)
			links = append(links, path)
		}
		if err != nil {
			return err
		}
	}
	if len(links) == 0 {
		return nil
	}
	// Links may chain through other links, so check where each one resolves,
	// following links in dangling targets as far as they exist.
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	realRoot, ok := resolveLinks(absRoot)
	for _, link := range links {
		absLink, err := filepath.Abs(link)
		if err != nil {
			return err
		}
		resolved, linkOK := resolveLinks(absLink)
		if !ok || !linkOK || !Within(realRoot, resolved) {
			return fmt.Errorf("filex: archive symlink %q resolves outside target directory", link)
		}
	}
	return nil
}

// ExtractZip extracts the zip archive at archivePath into root, creating
// directories as needed. Every non-directory entry is written as a regular
// file that is at least readable and writable by its owner; symbolic links
// stored in the archive are written as regular files holding the link target.
// It returns an error when an entry would escape root or be written through a
// symbolic link that already exists under root. Entries extracted before an
// error are left in place.
// The total size and number of entries are not limited.
func ExtractZip(archivePath, root string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()
	for _, entry := range reader.File {
		path, err := writablePath(root, entry.Name)
		if err != nil {
			return err
		}
		if entry.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return err
			}
			continue
		}
		content, err := entry.Open()
		if err != nil {
			return err
		}
		err = writeFile(path, content, entry.Mode().Perm()|0o600)
		_ = content.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// maxLinkHops is the number of symbolic links resolveLinks follows before
// giving up.
const maxLinkHops = 255

// resolveLinks returns the absolute path abs with every symbolic link in it
// followed, including links reached through other links. Components that do
// not exist are joined lexically, so a dangling link resolves to the location
// it would create. It reports false on a link loop or an unreadable link.
func resolveLinks(abs string) (string, bool) {
	hops := maxLinkHops
	return resolveFrom("", abs, &hops)
}

// resolveFrom resolves name relative to the already resolved directory base,
// starting from its volume root when name is absolute; hops counts the links
// that may still be followed.
func resolveFrom(base, name string, hops *int) (string, bool) {
	if filepath.IsAbs(name) {
		volume := filepath.VolumeName(name)
		base, name = volume+string(filepath.Separator), name[len(volume):]
	}
	for part := range strings.FieldsFuncSeq(name, func(r rune) bool { return r == '/' || r == filepath.Separator }) {
		switch part {
		case ".":
			continue
		case "..":
			base = filepath.Dir(base)
			continue
		}
		next := filepath.Join(base, part)
		info, err := os.Lstat(next)
		if err != nil || info.Mode()&os.ModeSymlink == 0 {
			base = next
			continue
		}
		target, err := os.Readlink(next)
		if err != nil || *hops == 0 {
			return "", false
		}
		(*hops)--
		resolved, ok := resolveFrom(base, target, hops)
		if !ok {
			return "", false
		}
		base = resolved
	}
	return base, true
}

// writablePath returns the extraction path of the entry name, or an error when
// the path escapes root or passes through an existing symbolic link.
func writablePath(root, name string) (string, error) {
	path, err := EntryPath(root, name)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == "." {
		return path, err
	}
	current := root
	for part := range strings.SplitSeq(relative, string(filepath.Separator)) {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			break
		}
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("filex: archive entry %q passes through a symlink", name)
		}
	}
	return path, nil
}

// writeFile creates the parent directories of path and writes content to it.
func writeFile(path string, content io.Reader, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(file, content)
	return errors.Join(copyErr, file.Close())
}
