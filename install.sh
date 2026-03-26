#!/usr/bin/env bash

set -euo pipefail

REPO="${RADCLI_INSTALL_REPOSITORY:-lloydhumphreys/radcli}"
VERSION="${RADCLI_VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-}"

usage() {
  cat <<'EOF'
Install radcli from GitHub Releases.

Usage:
  install.sh [--version <tag>] [--dir <install-dir>] [--repo <owner/repo>]

Environment:
  RADCLI_VERSION              Release tag to install, e.g. v0.1.0
  INSTALL_DIR                 Installation directory
  RADCLI_INSTALL_REPOSITORY   Alternate GitHub repository in owner/repo form

Examples:
  curl -fsSL https://raw.githubusercontent.com/lloydhumphreys/radcli/main/install.sh | bash
  curl -fsSL https://raw.githubusercontent.com/lloydhumphreys/radcli/main/install.sh | bash -s -- --version v0.1.0
  curl -fsSL https://raw.githubusercontent.com/lloydhumphreys/radcli/main/install.sh | INSTALL_DIR="$HOME/.local/bin" bash
EOF
}

while [ $# -gt 0 ]; do
  case "$1" in
    --version)
      VERSION="$2"
      shift 2
      ;;
    --dir)
      INSTALL_DIR="$2"
      shift 2
      ;;
    --repo)
      REPO="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}

need_cmd curl
need_cmd tar

tmpdir="$(mktemp -d)"
cleanup() {
  rm -rf "$tmpdir"
}
trap cleanup EXIT

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$os" in
  darwin|linux) ;;
  *)
    echo "unsupported operating system: $os" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)
    echo "unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

if [ -z "$INSTALL_DIR" ]; then
  if [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
  elif [ -w "/opt/homebrew/bin" ]; then
    INSTALL_DIR="/opt/homebrew/bin"
  else
    INSTALL_DIR="$HOME/.local/bin"
  fi
fi

mkdir -p "$INSTALL_DIR"

release_api="https://api.github.com/repos/$REPO/releases"
if [ "$VERSION" = "latest" ]; then
  release_api="$release_api/latest"
else
  case "$VERSION" in
    v*) ;;
    *) VERSION="v$VERSION" ;;
  esac
  release_api="$release_api/tags/$VERSION"
fi

release_json="$tmpdir/release.json"
curl -fsSL "$release_api" -o "$release_json"

tag="$(sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' "$release_json" | head -n 1)"
if [ -z "$tag" ]; then
  echo "failed to resolve release tag from $release_api" >&2
  exit 1
fi

version="${tag#v}"
archive="radcli_${version}_${os}_${arch}.tar.gz"
checksum_file="checksums.txt"
base_url="https://github.com/$REPO/releases/download/$tag"

archive_path="$tmpdir/$archive"
checksum_path="$tmpdir/$checksum_file"

curl -fsSL "$base_url/$archive" -o "$archive_path"
curl -fsSL "$base_url/$checksum_file" -o "$checksum_path"

expected_checksum="$(grep "  $archive\$" "$checksum_path" | awk '{print $1}')"
if [ -z "$expected_checksum" ]; then
  echo "failed to find checksum for $archive" >&2
  exit 1
fi

actual_checksum=""
if command -v shasum >/dev/null 2>&1; then
  actual_checksum="$(shasum -a 256 "$archive_path" | awk '{print $1}')"
elif command -v sha256sum >/dev/null 2>&1; then
  actual_checksum="$(sha256sum "$archive_path" | awk '{print $1}')"
elif command -v openssl >/dev/null 2>&1; then
  actual_checksum="$(openssl dgst -sha256 "$archive_path" | awk '{print $NF}')"
else
  echo "missing checksum tool: shasum, sha256sum, or openssl" >&2
  exit 1
fi

if [ "$actual_checksum" != "$expected_checksum" ]; then
  echo "checksum verification failed for $archive" >&2
  echo "expected: $expected_checksum" >&2
  echo "actual:   $actual_checksum" >&2
  exit 1
fi

tar -xzf "$archive_path" -C "$tmpdir"
binary_path="$tmpdir/rad"
if [ ! -f "$binary_path" ]; then
  echo "release archive did not contain rad" >&2
  exit 1
fi

install_path="$INSTALL_DIR/rad"
if [ -w "$INSTALL_DIR" ]; then
  install -m 0755 "$binary_path" "$install_path"
else
  sudo install -m 0755 "$binary_path" "$install_path"
fi

echo "Installed radcli $tag to $install_path"

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo "Add $INSTALL_DIR to your PATH to run 'rad' directly."
    ;;
esac
