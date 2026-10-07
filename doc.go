// Package support provides generic value helpers and the type constraints
// shared by its subpackages.
//
// Domain helpers live in subpackages:
//
//   - str: string manipulation and the fluent Stringable
//   - arr: slice helpers
//   - mapx: generic map helpers
//   - data: dot-notation access to nested map[string]any and []any values
//   - number: number formatting
//   - convert: conversion between basic types
//
// Example:
//
//	name := support.Default(input.Nickname, input.Name, "anonymous")
//	limit := support.DerefOr(req.Limit, 20)
package support
