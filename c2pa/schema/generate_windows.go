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

// Schema generation skipped on Windows pending resolution of upstream OpenSSL issue.
// The schema.go file is checked in and will be used as-is.
// TODO: Re-enable schema generation when c2pa-rs export_schema can build without OpenSSL on Windows.
//
// //go:generate cmd /C "cd ../../c2pa-rs/export_schema && cargo run --release --no-default-features"
// //go:generate powershell -NoProfile -Command "(Get-Content ../../c2pa-rs/export_schema/target/schema/Reader.schema.json) -replace '\"title\": \"Reader\"', '\"title\": \"ManifestStore\"' | Set-Content ManifestStore.schema.json"
// //go:generate npx --yes quicktype -s schema -l go --package schema --out schema.go ../../c2pa-rs/export_schema/target/schema/ManifestDefinition.schema.json ManifestStore.schema.json ../../c2pa-rs/export_schema/target/schema/Settings.schema.json
// //go:generate gofmt -w schema.go
// //go:generate cmd /C del /Q /F ManifestStore.schema.json
