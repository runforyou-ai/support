// Package filex provides file-system helpers that the standard os, io/fs and
// path/filepath packages do not offer directly: safe archive extraction,
// atomic file replacement, size-based rotating files, lexical containment
// checks and binary content sniffing.
//
// Example:
//
//	// Extract an untrusted archive; entries cannot escape dir.
//	err := filex.ExtractTarGz("bundle.tar.gz", dir)
//
//	// Replace a config file so readers never see a partial write.
//	err = filex.WriteAtomic("config.json", data, 0o644)
//
//	// Append logs, keeping config.log plus three 10 MiB backups.
//	log, err := filex.OpenRotating("config.log", 10<<20, 3)
package filex
