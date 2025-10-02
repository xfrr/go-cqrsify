// Package apix provides small, framework-agnostic helpers for building
// consistent HTTP APIs in Go.
//
// It focuses on three stable standards:
//
//   - JSON:API (a pragmatic subset) for success responses,
//   - RFC 7807 (Problem Details) for errors, and
//   - RFC 7386 (JSON Merge Patch) for partial updates.
//
// Design notes:
//   - No external dependencies; only encoding/json and net/http.
//   - Minimal, idiomatic types that fit most REST APIs.
//   - Helpers return/encode plain structs; you can reuse them with any router.
//   - For errors, use Problem (RFC 7807) instead of JSON:API "errors" arrays.
package apix
