package validate

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
)

// zoneSources returns the zone directories and zip archives that
// time.LoadLocation searches before its embedded and GOROOT data, in order:
// $ZONEINFO, read once as the time package does, then the system directories
// of Unix platforms.
var zoneSources = sync.OnceValue(func() []string {
	var sources []string
	if env := os.Getenv("ZONEINFO"); env != "" {
		sources = append(sources, env)
	}
	switch runtime.GOOS {
	case "windows", "plan9", "android", "ios", "js", "wasip1":
	default:
		sources = append(sources, "/usr/share/zoneinfo", "/usr/share/lib/zoneinfo", "/usr/lib/locale/TZ", "/etc/zoneinfo")
	}
	return sources
})

// Timezone reports whether name is an IANA time zone name, such as
// "Asia/Shanghai" or "UTC", that time.LoadLocation can load. Names are
// case-sensitive on every platform: name must match the zone's name exactly
// as stored in the source time.LoadLocation reads it from, so
// "asia/shanghai" is rejected even where a case-insensitive file system lets
// time.LoadLocation open it. It returns false for "" and "Local", which
// time.LoadLocation maps to UTC and the system zone.
//
// Zone data is read from the system at run time, and $ZONEINFO is read once,
// as by the time package. Programs that must validate zones on systems
// without a zone database, such as minimal containers or Windows hosts
// without Go installed, should embed one by importing _ "time/tzdata" in
// their main package.
func Timezone(name string) bool {
	switch name {
	case "", "Local":
		return false
	case "UTC":
		return true
	}
	if _, err := time.LoadLocation(name); err != nil {
		return false
	}
	// The first source holding a parsable zone file for name is the one
	// time.LoadLocation used. Zip archives, like the embedded and GOROOT data
	// searched last, match names exactly.
	for _, source := range zoneSources() {
		if strings.HasSuffix(source, ".zip") {
			if zipHasZone(source, name) {
				return true
			}
			continue
		}
		data, err := os.ReadFile(filepath.Join(source, filepath.FromSlash(name)))
		if err != nil {
			continue
		}
		if _, err := time.LoadLocationFromTZData(name, data); err != nil {
			continue
		}
		dir := source
		for part := range strings.SplitSeq(name, "/") {
			entries, err := os.ReadDir(dir)
			if err != nil || !slices.ContainsFunc(entries, func(e os.DirEntry) bool { return e.Name() == part }) {
				return false
			}
			dir = filepath.Join(dir, part)
		}
		return true
	}
	return true
}

// zipHasZone reports whether the zip archive at path holds an entry named
// exactly name with parsable zone data.
func zipHasZone(path, name string) bool {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return false
	}
	defer func() { _ = reader.Close() }()
	for _, entry := range reader.File {
		if entry.Name != name {
			continue
		}
		file, err := entry.Open()
		if err != nil {
			return false
		}
		data, err := io.ReadAll(file)
		_ = file.Close()
		if err != nil {
			return false
		}
		_, err = time.LoadLocationFromTZData(name, data)
		return err == nil
	}
	return false
}
