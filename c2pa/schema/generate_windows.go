//go:build windows
// +build windows

// Schema generation directives for Windows.
//
// Run `go generate ./c2pa/schema/...` to regenerate schema.go.
//
// Notes:
//   - Forward slashes are used in paths because Go's //go:generate parser
//     treats backslashes inside "..." as escape sequences. Both cmd.exe and
//     PowerShell accept '/' as a path separator.
//   - PowerShell is used for the sed-equivalent string replacements; inner
//     double quotes are escaped with \" so go:generate passes a single arg.
//   - cmd /C is used for `copy` and `del` because those are shell builtins.

package schema

//go:generate cmd /C "cd ../../c2pa-rs/export_schema && cargo run --release --features c2pa/rust_native_crypto"
//go:generate cmd /C copy /Y ..\..\c2pa-rs\export_schema\target\schema\Settings.schema.json Settings.schema.json
//go:generate powershell -NoProfile -Command "(Get-Content ../../c2pa-rs/export_schema/target/schema/Builder.schema.json) -replace '\"title\": \"Builder\"', '\"title\": \"ManifestDefinition\"' | Set-Content ManifestDefinition.schema.json"
//go:generate powershell -NoProfile -Command "(Get-Content ../../c2pa-rs/export_schema/target/schema/Reader.schema.json) -replace '\"title\": \"Reader\"', '\"title\": \"ManifestStore\"' | Set-Content ManifestStore.schema.json"
//go:generate npx --yes quicktype -s schema -l go --package schema --out schema.go ManifestDefinition.schema.json ManifestStore.schema.json Settings.schema.json
//go:generate cmd /C del /Q /F ManifestDefinition.schema.json ManifestStore.schema.json Settings.schema.json
