// Package support provides generic value helpers and the type constraints
// shared by its subpackages.
//
// Domain helpers live in subpackages:
//
//   - str: string manipulation and the fluent Stringable
//   - arr: slice helpers
//   - mapx: generic map helpers
//   - set: a generic set type
//   - data: dot-notation access to nested map[string]any and []any values
//   - number: number formatting
//   - convert: conversion between basic types
//   - random: cryptographically secure random bytes, strings and tokens
//   - validate: phone number and time zone validation
//   - filex: archive extraction, rotating and atomic file writes, path checks
//   - iox: io.Writer helpers
//
// Example:
//
//	limit := support.DerefOr(req.Limit, 20)
//	nickname := support.NilIfZero(input.Nickname)
//	err := support.Catch(func() error { return handle(req) })
//	var pe *support.PanicError
//	if errors.As(err, &pe) {
//		log.Printf("recovered: %v", pe.Value)
//	}
package support
