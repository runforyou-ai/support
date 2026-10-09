package validate

import (
	"encoding/binary"
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

// zipHasZone reports whether the zip archive at path holds parsable zone data
// for name as time.LoadLocation reads it: the end of central directory
// record must close the file, and the first central directory entry named
// exactly name must be stored uncompressed.
func zipHasZone(path, name string) bool {
	data, err := os.ReadFile(path)
	if err != nil || len(data) < 22 {
		return false
	}
	le := binary.LittleEndian
	tail := data[len(data)-22:]
	if le.Uint32(tail) != 0x06054b50 {
		return false
	}
	count, size, offset := int(le.Uint16(tail[10:])), int(le.Uint32(tail[12:])), int(le.Uint32(tail[16:]))
	if offset > len(data) || size > len(data)-offset {
		return false
	}
	dir := data[offset : offset+size]
	for range count {
		if len(dir) < 46 || le.Uint32(dir) != 0x02014b50 {
			return false
		}
		method, length := le.Uint16(dir[10:]), int(le.Uint32(dir[24:]))
		nameLen, extraLen, commentLen := int(le.Uint16(dir[28:])), int(le.Uint16(dir[30:])), int(le.Uint16(dir[32:]))
		local := int(le.Uint32(dir[42:]))
		if len(dir) < 46+nameLen+extraLen+commentLen {
			return false
		}
		entryName := string(dir[46 : 46+nameLen])
		dir = dir[46+nameLen+extraLen+commentLen:]
		if entryName != name {
			continue
		}
		if method != 0 || local > len(data) || len(data)-local < 30+nameLen {
			return false
		}
		header := data[local:]
		if le.Uint32(header) != 0x04034b50 || le.Uint16(header[8:]) != method ||
			int(le.Uint16(header[26:])) != nameLen || string(header[30:30+nameLen]) != name {
			return false
		}
		start := local + 30 + nameLen + int(le.Uint16(header[28:]))
		if start > len(data) || length > len(data)-start {
			return false
		}
		_, err := time.LoadLocationFromTZData(name, data[start:start+length])
		return err == nil
	}
	return false
}
