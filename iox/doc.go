// Package iox provides io.Reader and io.Writer helpers that the standard io
// package does not offer directly.
//
// HeadTailBuffer collects unbounded output, such as a child process's stdout
// and stderr, while keeping memory bounded: it retains the first and the last
// bytes written and drops the middle.
//
// Example:
//
//	buf := iox.NewHeadTailBuffer(4<<10, 16<<10, "\n... output truncated ...\n")
//	cmd.Stdout, cmd.Stderr = buf, buf
//	_ = cmd.Run()
//	log.Print(buf.String())
package iox
