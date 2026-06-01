// Package c2pa provides Go bindings for the c2pa-rs Content Authenticity SDK.
//
// Keywords: c2pa, c2pa-go, c2pa-rs, Content Authenticity, Content Integrity, Content Credentials, content provenance.
//
// The package exposes the C2PA reader, builder, signer, settings and context
// APIs through cgo so Go programs can read, create and sign C2PA manifests
// (also known as Content Credentials) on any asset format that c2pa-rs
// supports (JPEG, PNG, TIFF, MP4, and others). Use it to add Content
// Authenticity / Content Integrity metadata to media assets and to verify
// existing C2PA manifests.
//
// All operations are scoped to a [Context]. A context optionally carries
// [Settings] (trust list, verification policy, builder defaults, etc.). Create
// a context either with [NewContext] for defaults, or via
// [NewContextBuilder] when you need to attach settings:
//
//	ctx, err := c2pa.NewContext()
//	if err != nil { ... }
//	defer ctx.Close()
//
// # Reading manifests
//
// Use [NewReader] with [Reader.WithFile] to inspect manifests embedded in or
// accompanying an asset:
//
//	r, err := c2pa.NewReader(ctx)
//	if err != nil { ... }
//	err = r.WithFile("signed.jpg")
//	if err != nil { ... }
//	defer r.Close()
//	fmt.Println(r.Json())
//
// # Signing assets
//
// Implement the [Signer] interface (or wrap your own signing callback) and
// drive [Builder.Sign]:
//
//	b, err := c2pa.BuilderFromJson(ctx, manifestJson)
//	if err != nil { ... }
//	defer b.Close()
//	manifest, err := b.Sign("in.jpg", "out.jpg", signer)
//
// The [Builder] type mirrors the upstream API, including SetIntent,
// SetRemoteUrl, AddAction, AddResource, AddIngredient, and archive helpers.
//
// # Linking
//
// This package links against the c2pa_c shared/static library built from the
// c2pa-rs c2pa_c_ffi crate. You can either download a prebuilt library
// (no Rust toolchain required) using the fetchlib helper, or rebuild the
// C FFI from source via go generate. See the "Building" section of the
// project README for the full recipe:
//
// https://github.com/duggaraju/c2pa-go#building
package c2pa
