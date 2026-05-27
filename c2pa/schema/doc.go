// Package schema contains Go types generated from the c2pa-rs JSON schemas
// for C2PA Content Authenticity / Content Integrity manifests (Content
// Credentials).
//
// Keywords: c2pa, Content Authenticity, Content Integrity, Content Credentials,
// manifest schema, c2pa-rs.
//
// The package exposes three top-level types matching the JSON shapes used by
// the underlying c2pa-rs library:
//
//   - ManifestDefinition: the manifest JSON accepted by Builder.from_json.
//     (Generated from c2pa-rs ManifestDefinition.schema.json.)
//   - ManifestStore: the manifest store JSON returned by Reader.json.
//     (Generated from c2pa-rs Reader.schema.json with the top-level title
//     renamed to ManifestStore.)
//   - Settings: the settings JSON consumed by Settings.from_toml /
//     LoadSettingsFromJson.
//
// To regenerate this package, run from the repository root:
//
//	go generate ./c2pa/schema/...
package schema
