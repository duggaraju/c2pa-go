// Schema regeneration. Always runs on `go generate`, regardless of build tags.
//
// Run `go generate ./c2pa/...` to refresh the generated schema package from
// the c2pa-rs JSON schemas.

package c2pa

//go:generate go generate github.com/duggaraju/c2pa-go/c2pa/schema
