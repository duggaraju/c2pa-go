//go:build linux && cgo
// +build linux,cgo

package c2pa

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	schema "github.com/duggaraju/c2pa-go/c2pa/schema"
	"github.com/stretchr/testify/assert"
)

func assertManifestStoreHasActiveManifest(t *testing.T, store *schema.ManifestStore) {
	t.Helper()
	if assert.NotNil(t, store) && assert.NotNil(t, store.ActiveManifest) {
		assert.NotEmpty(t, *store.ActiveManifest)
		assert.NotEmpty(t, store.Manifests)
		manifest, ok := store.Manifests[*store.ActiveManifest]
		if assert.True(t, ok) {
			_ = manifest
		}
	}
}

func assertCrJsonHasManifests(t *testing.T, raw string) map[string]any {
	t.Helper()
	var cr map[string]any
	if assert.NoError(t, json.Unmarshal([]byte(raw), &cr)) && assert.NotNil(t, cr) {
		manifests, ok := cr["manifests"]
		if assert.True(t, ok) {
			manifestList, ok := manifests.([]any)
			if assert.True(t, ok) {
				assert.NotEmpty(t, manifestList)
			}
		}
	}
	return cr
}

func TestReaderNilContext(t *testing.T) {
	_, err := NewReader(nil)
	assert.Error(t, err)
}

func TestC2paError_AfterReaderFailure(t *testing.T) {
	// Trigger a failure that happens inside the C2PA library (so c2pa_error is populated),
	// not a Go-side error like "os.Open" failing.
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jpg")
	err := os.WriteFile(path, nil, 0o644)
	assert.NoError(t, err)

	ctx, ctxErr := NewContext()
	assert.NoError(t, ctxErr)
	defer ctx.Close()

	r, err := NewReader(ctx)
	assert.NoError(t, err)
	defer r.Close()

	readerErr := r.WithFile(path)
	assert.Error(t, readerErr)

	lastErr := c2paError()
	assert.NotEmpty(t, lastErr)
}

func TestReaderWithFile_NotFound(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	r, err := NewReader(ctx)
	assert.NoError(t, err)
	defer r.Close()

	err = r.WithFile("/nonexistent/file/path.jpg")
	assert.Error(t, err)
}

func TestReaderWithFile_Valid(t *testing.T) {
	ctx, ctxErr := NewContext()
	assert.NoError(t, ctxErr)
	defer ctx.Close()

	r, err := NewReader(ctx)
	assert.NoError(t, err)
	defer r.Close()

	err = r.WithFile(fixturePath("sdk", "tests", "fixtures", "C.jpg"))
	assert.NoError(t, err)
	assert.NotEmpty(t, r.Json())
}

func TestNewDefaultReader_Valid(t *testing.T) {
	r, err := NewDefaultReader()
	assert.NoError(t, err)
	defer r.Close()

	err = r.WithFile(fixturePath("sdk", "tests", "fixtures", "C.jpg"))
	assert.NoError(t, err)
	assert.NotEmpty(t, r.Json())
	assert.NotEmpty(t, r.CrJson())
	assert.NotEmpty(t, r.DetailedJson())
}

func TestReaderMetadataAndManifestParsing(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	reader, err := NewReader(ctx)
	assert.NoError(t, err)
	defer reader.Close()

	asset := fixturePath("sdk", "tests", "fixtures", "C.jpg")
	assert.NoError(t, reader.WithFile(asset))
	assert.Empty(t, reader.RemoteUrl())
	assert.True(t, reader.IsEmbedded())
	assert.NotEmpty(t, reader.Json())
	assert.NotEmpty(t, reader.DetailedJson())
	assert.NotEmpty(t, reader.CrJson())

	manifest, err := reader.Manifest()
	assert.NoError(t, err)
	assertManifestStoreHasActiveManifest(t, manifest)

	detailed, err := reader.DetailedManifest()
	assert.NoError(t, err)
	assertManifestStoreHasActiveManifest(t, detailed)
	if assert.NotNil(t, detailed) {
		assert.NotEmpty(t, detailed.Manifests)
	}

	assertCrJsonHasManifests(t, reader.CrJson())
}

func TestReaderManifestMethod(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	reader, err := NewReader(ctx)
	assert.NoError(t, err)
	defer reader.Close()

	err = reader.WithFile(fixturePath("sdk", "tests", "fixtures", "C.jpg"))
	assert.NoError(t, err)

	manifest, err := reader.Manifest()
	assert.NoError(t, err)
	assertManifestStoreHasActiveManifest(t, manifest)
	if manifest != nil && manifest.ActiveManifest != nil {
		active := manifest.Manifests[*manifest.ActiveManifest]
		assert.NotNil(t, active.Format)
		if active.Format != nil {
			assert.NotEmpty(t, *active.Format)
		}
	}
}

func TestReaderDetailedManifestMethod(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	reader, err := NewReader(ctx)
	assert.NoError(t, err)
	defer reader.Close()

	err = reader.WithFile(fixturePath("sdk", "tests", "fixtures", "C.jpg"))
	assert.NoError(t, err)

	detailed, err := reader.DetailedManifest()
	assert.NoError(t, err)
	assertManifestStoreHasActiveManifest(t, detailed)
	if detailed != nil && detailed.ActiveManifest != nil {
		_, ok := detailed.Manifests[*detailed.ActiveManifest]
		assert.True(t, ok)
	}
}

func TestReaderCrJsonMethod(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	reader, err := NewReader(ctx)
	assert.NoError(t, err)
	defer reader.Close()

	err = reader.WithFile(fixturePath("sdk", "tests", "fixtures", "C.jpg"))
	assert.NoError(t, err)

	cr := assertCrJsonHasManifests(t, reader.CrJson())
	if manifests, ok := cr["manifests"].([]any); ok && len(manifests) > 0 {
		first, ok := manifests[0].(map[string]any)
		if assert.True(t, ok) {
			claim, ok := first["claim"].(map[string]any)
			if assert.True(t, ok) {
				assert.NotEmpty(t, claim["dc:title"])
			}
		}
	}
}

func TestReaderWithManifestData(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	reader, err := NewReader(ctx)
	assert.NoError(t, err)
	defer reader.Close()

	asset, err := os.Open(fixturePath("sdk", "tests", "fixtures", "cloud.jpg"))
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, asset.Close())
	}()

	manifestData, err := os.ReadFile(fixturePath("sdk", "tests", "fixtures", "cloud_manifest.c2pa"))
	assert.NoError(t, err)

	err = reader.WithManifestData("image/jpeg", asset, manifestData)
	assert.NoError(t, err)
	assert.NotEmpty(t, reader.Json())

	manifest, err := reader.Manifest()
	assert.NoError(t, err)
	assert.NotNil(t, manifest)
}

func TestReaderWithFragment(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	reader, err := NewReader(ctx)
	assert.NoError(t, err)
	defer reader.Close()

	asset, err := os.Open(fixturePath("sdk", "tests", "fixtures", "dashinit.mp4"))
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, asset.Close())
	}()

	fragment, err := os.Open(fixturePath("sdk", "tests", "fixtures", "dash1.m4s"))
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, fragment.Close())
	}()

	err = reader.WithFragment("video/mp4", asset, fragment)
	assert.NoError(t, err)
	assert.NotEmpty(t, reader.Json())

	manifest, err := reader.Manifest()
	assert.NoError(t, err)
	assert.NotNil(t, manifest)
}

func TestReaderResourceToFile(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	reader, err := NewReader(ctx)
	assert.NoError(t, err)
	defer reader.Close()

	asset := fixturePath("sdk", "tests", "fixtures", "C.jpg")
	err = reader.WithFile(asset)
	assert.NoError(t, err)

	output := filepath.Join(t.TempDir(), "thumbnail.jpg")
	const uri = "self#jumbf=c2pa.assertions/c2pa.thumbnail.claim.jpeg"

	n, err := reader.ResourceToFile(uri, output)
	assert.NoError(t, err)
	assert.Greater(t, n, int64(0))

	data, err := os.ReadFile(output)
	assert.NoError(t, err)
	assert.Len(t, data, int(n))
}

func TestReaderResourceToFileUnknownURI(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	reader, err := NewReader(ctx)
	assert.NoError(t, err)
	defer reader.Close()

	asset := fixturePath("sdk", "tests", "fixtures", "C.jpg")
	err = reader.WithFile(asset)
	assert.NoError(t, err)

	_, err = reader.ResourceToFile("self#jumbf=c2pa.assertions/does-not-exist", filepath.Join(t.TempDir(), "missing.bin"))
	assert.Error(t, err)
}
