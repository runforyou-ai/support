package filex_test

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/runforyou-ai/support/filex"
)

func ExampleEntryPath() {
	root := filepath.FromSlash("/srv/extract")
	path, err := filex.EntryPath(root, "docs/readme.txt")
	fmt.Println(filepath.ToSlash(path), err)
	_, err = filex.EntryPath(root, "../../etc/passwd")
	fmt.Println(err)
	// Output:
	// /srv/extract/docs/readme.txt <nil>
	// filex: archive entry "../../etc/passwd" escapes target directory
}

func ExampleExtractTarGz() {
	dir, err := os.MkdirTemp("", "filex-example")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	// Build a small archive.
	archive := filepath.Join(dir, "bundle.tar.gz")
	out, _ := os.Create(archive)
	gz := gzip.NewWriter(out)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "docs/hello.txt", Mode: 0o644, Size: 5, Typeflag: tar.TypeReg})
	_, _ = tw.Write([]byte("hello"))
	_ = tw.Close()
	_ = gz.Close()
	_ = out.Close()

	root := filepath.Join(dir, "out")
	fmt.Println(filex.ExtractTarGz(archive, root))
	data, _ := os.ReadFile(filepath.Join(root, "docs", "hello.txt"))
	fmt.Println(string(data))
	// Output:
	// <nil>
	// hello
}

func ExampleExtractZip() {
	dir, err := os.MkdirTemp("", "filex-example")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	// Build an archive with a malicious entry.
	archive := filepath.Join(dir, "evil.zip")
	out, _ := os.Create(archive)
	zw := zip.NewWriter(out)
	w, _ := zw.Create("../evil.txt")
	_, _ = w.Write([]byte("x"))
	_ = zw.Close()
	_ = out.Close()

	fmt.Println(filex.ExtractZip(archive, filepath.Join(dir, "out")))
	// Output: filex: archive entry "../evil.txt" escapes target directory
}

func ExampleOpenRotating() {
	dir, err := os.MkdirTemp("", "filex-example")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	path := filepath.Join(dir, "app.log")
	log, err := filex.OpenRotating(path, 10, 2)
	if err != nil {
		panic(err)
	}
	_, _ = fmt.Fprintln(log, "first line")  // 11 bytes: exceeds 10, rotates
	_, _ = fmt.Fprintln(log, "second line") // rotates again
	_, _ = fmt.Fprintln(log, "third")
	_ = log.Close()

	for _, name := range []string{"app.log", "app.log.1", "app.log.2"} {
		data, _ := os.ReadFile(filepath.Join(dir, name))
		fmt.Printf("%s: %q\n", name, data)
	}
	// Output:
	// app.log: "third\n"
	// app.log.1: "second line\n"
	// app.log.2: "first line\n"
}

func ExampleWriteAtomic() {
	dir, err := os.MkdirTemp("", "filex-example")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	path := filepath.Join(dir, "config.json")
	fmt.Println(filex.WriteAtomic(path, []byte(`{"v":1}`), 0o644))
	fmt.Println(filex.WriteAtomic(path, []byte(`{"v":2}`), 0o644))
	data, _ := os.ReadFile(path)
	fmt.Println(string(data))
	// Output:
	// <nil>
	// <nil>
	// {"v":2}
}

func ExampleWithin() {
	root := filepath.FromSlash("/srv/data")
	fmt.Println(filex.Within(root, filepath.FromSlash("/srv/data/a/b.txt")))
	fmt.Println(filex.Within(root, filepath.FromSlash("/srv/data/../secret")))
	fmt.Println(filex.Within(root, filepath.FromSlash("/srv/database")))
	// Output:
	// true
	// false
	// false
}

func ExampleIsBinary() {
	fmt.Println(filex.IsBinary([]byte("plain text\n")))
	fmt.Println(filex.IsBinary([]byte{0x89, 'P', 'N', 'G', 0x00}))
	// Output:
	// false
	// true
}
