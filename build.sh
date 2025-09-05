#!/usr/bin/env bash
set -euo pipefail

# Simple helper to build the go submodule.
# Usage: ./build.sh [debug|release]
# Default: debug

BUILD_TYPE=${1:-debug}
if [ "$BUILD_TYPE" = "release" ]; then
  CARGO_ARGS="--release"
  TARGET_DIR="target/release"
  GO_ARGS="-ldfalgs '-s -w'"
else
  GO_ARGS="-gcflags '-N -l'"
  CARGO_ARGS=""
  TARGET_DIR="target/debug"
fi

ROOT_DIR="$(git rev-parse --show-toplevel)"
SUBMODULE_DIR="$ROOT_DIR/c2pa-rs/c2pa_c_ffi"

if [ ! -d "$SUBMODULE_DIR" ]; then
  echo "Warning: submodule directory not found at $SUBMODULE_DIR"
  echo "Add the submodule first, e.g.: git submodule update --init"
  exit 1
fi

pushd "$SUBMODULE_DIR" >/dev/null
echo "Building c2pa-rs (cargo build $CARGO_ARGS)"
cargo build $CARGO_ARGS --features rust_native_crypto

popd >/dev/null

echo "Done. Built artifacts are in: $SUBMODULE_DIR/$TARGET_DIR"

go build ./lib $GO_ARGS
go build ./example $GO_ARGS
