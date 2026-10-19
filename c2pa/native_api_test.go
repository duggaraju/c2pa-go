//go:build linux && cgo
// +build linux,cgo

package c2pa

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Check actual declarations and calls, not the C examples in header comments.
func TestNonDeprecatedNativeAPICoverage(t *testing.T) {
	native, err := parser.ParseFile(token.NewFileSet(), "native.go", nil, 0)
	require.NoError(t, err)
	calls := make(map[string]bool)
	ast.Inspect(native, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if pkg, ok := selector.X.(*ast.Ident); ok && pkg.Name == "C" {
			calls[selector.Sel.Name] = true
		}
		return true
	})

	comments := regexp.MustCompile(`(?s)/\*.*?\*/|//[^\n]*`)
	function := regexp.MustCompile(`\b(c2pa_\w+)\s*\(`)
	helper, err := os.ReadFile("c2pa_helper.h")
	require.NoError(t, err)
	for _, match := range function.FindAllStringSubmatch(comments.ReplaceAllString(string(helper), ""), -1) {
		calls[match[1]] = true
	}

	// Audit every installed build flavor; either may be used by local/CI builds.
	headers, err := filepath.Glob(filepath.Join("..", "c2pa-rs", "target", "*", "c2pa.h"))
	require.NoError(t, err)
	require.NotEmpty(t, headers, "generate or fetch the native library before testing API coverage")
	declaration := regexp.MustCompile(`C2PA_API\s+extern\s+([^;]+);`)
	for _, path := range headers {
		t.Run(filepath.Base(filepath.Dir(path)), func(t *testing.T) {
			header, err := os.ReadFile(path)
			require.NoError(t, err)
			declarations := declaration.FindAllStringSubmatch(comments.ReplaceAllString(string(header), ""), -1)
			require.NotEmpty(t, declarations, "no C API declarations found")
			active := 0
			for _, decl := range declarations {
				if strings.Contains(decl[1], "C2PA_DEPRECATED") {
					continue
				}
				match := function.FindStringSubmatch(decl[1])
				require.Len(t, match, 2, "unrecognized C API declaration: %s", decl[1])
				active++
				if !calls[match[1]] {
					t.Errorf("non-deprecated API %s has no Go binding", match[1])
				}
			}
			require.Positive(t, active, "no non-deprecated C APIs found")
			t.Logf("checked %d non-deprecated C APIs", active)
		})
	}
}
