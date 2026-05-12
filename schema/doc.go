// Package schema contains Go types generated from the c2pa-rs JSON schemas.
//
// The package exposes three top-level types matching the JSON shapes used by
// the underlying c2pa-rs library:
//
//   - ManifestDefinition: the manifest JSON accepted by Builder.from_json.
//     (Generated from c2pa-rs Builder.schema.json with the top-level title
//     renamed to ManifestDefinition.)
//   - ManifestStore: the manifest store JSON returned by Reader.json.
//     (Generated from c2pa-rs Reader.schema.json with the top-level title
//     renamed to ManifestStore.)
//   - Settings: the settings JSON consumed by Settings.from_toml /
//     LoadSettingsFromJson.
//
// To regenerate this package, run from the repository root:
//
//	go generate ./schema/...
package schema

//go:generate sh -c "cd ../c2pa-rs/export_schema && cargo run --release --features c2pa/rust_native_crypto"
//go:generate sh -c "sed 's/\"title\": \"Builder\"/\"title\": \"ManifestDefinition\"/' ../c2pa-rs/export_schema/target/schema/Builder.schema.json > ManifestDefinition.schema.json"
//go:generate sh -c "sed 's/\"title\": \"Reader\"/\"title\": \"ManifestStore\"/' ../c2pa-rs/export_schema/target/schema/Reader.schema.json > ManifestStore.schema.json"
//go:generate cp ../c2pa-rs/export_schema/target/schema/Settings.schema.json Settings.schema.json
//go:generate npx --yes quicktype -s schema -l go --package schema --out schema.go ManifestDefinition.schema.json ManifestStore.schema.json Settings.schema.json
//go:generate sh -c "rm -f ManifestDefinition.schema.json ManifestStore.schema.json Settings.schema.json"
