#!/bin/bash

set -euo pipefail

SRC=$(realpath "$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)")
APP=gausssql
PLATFORM=linux
ARCH=${ARCH:-amd64}
VER=${VER:-0.21.4-og.1}
GO_BIN=${GO_BIN:-go}

usage() {
  cat <<'EOF'
Usage: ./build-opengauss.sh [-a amd64|arm64] [-v version] [-g /path/to/go]

Build the openGauss-focused fork artifacts.
EOF
}

while getopts "a:v:g:h" opt; do
  case "$opt" in
    a) ARCH=$OPTARG ;;
    v) VER=$OPTARG ;;
    g) GO_BIN=$OPTARG ;;
    h)
      usage
      exit 0
      ;;
    *)
      usage
      exit 1
      ;;
  esac
done

case "$ARCH" in
  amd64|arm64) ;;
  *)
    echo "ERROR: unsupported arch '$ARCH' (expected amd64 or arm64)"
    exit 1
    ;;
esac

if [[ "$GO_BIN" == "go" && -x "$SRC/.local/go/bin/go" ]]; then
  GO_BIN="$SRC/.local/go/bin/go"
fi

BUILD_DIR=$SRC/build/$PLATFORM/$ARCH/$VER
BIN=$BUILD_DIR/$APP
OUT=$BUILD_DIR/$APP-$VER-$PLATFORM-$ARCH.tar.bz2
TAGS="postgres no_base no_chart"
LDFLAGS="-s -w -X github.com/xo/usql/text.CommandName=$APP -X github.com/xo/usql/text.CommandVersion=$VER"

mkdir -p "$BUILD_DIR"

echo "APP:         $APP/$VER ($PLATFORM/$ARCH)"
echo "BUILD TAGS:  $TAGS"
echo "LDFLAGS:     $LDFLAGS"
echo "BUILDING:    $BIN"

pushd "$SRC" >/dev/null

GOCACHE=${GOCACHE:-$SRC/.local/gocache}
mkdir -p "$GOCACHE"

CGO_ENABLED=0 \
GOCACHE="$GOCACHE" \
GOOS=$PLATFORM \
GOARCH=$ARCH \
"$GO_BIN" build \
  -mod=mod \
  -trimpath \
  -ldflags="$LDFLAGS" \
  -tags="$TAGS" \
  -o "$BIN" .

chmod +x "$BIN"
file "$BIN"
"$BIN" --version

cp "$SRC/LICENSE" "$BUILD_DIR/LICENSE"
tar -C "$BUILD_DIR" -cjf "$OUT" "$APP" LICENSE
sha256sum "$BIN" "$OUT"

popd >/dev/null
