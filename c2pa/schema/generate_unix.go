//go:build unix
// +build unix

// Schema generation directives for Unix-like systems (Linux, macOS, BSD).
//
// Run `go generate ./c2pa/schema/...` to regenerate schema.go.

package schema

//go:generate bash -c "cd ../../c2pa-rs/export_schema && cargo run --release --features c2pa/rust_native_crypto"
//go:generate sh -c "sed 's/\"title\": \"Reader\"/\"title\": \"ManifestStore\"/' ../../c2pa-rs/export_schema/target/schema/Reader.schema.json > ManifestStore.schema.json"
//go:generate npx --yes quicktype -s schema -l go --package schema --out schema.go ../../c2pa-rs/export_schema/target/schema/ManifestDefinition.schema.json ManifestStore.schema.json ../../c2pa-rs/export_schema/target/schema/Settings.schema.json
//go:generate gofmt -w schema.go
//go:generate sh -c "rm -f ManifestStore.schema.json"
