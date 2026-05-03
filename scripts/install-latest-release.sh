#!/usr/bin/env bash
# Install a published baillconnect-to-mqtt binary from GitHub Releases.

set -euo pipefail

repo="${REPO:-jonatak/baillconnect-to-mqtt}"
version="${VERSION:-${BUILD_VERSION:-latest}}"
os="${BUILD_OS:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
bin_home="${XDG_BIN_HOME:-$HOME/.local/bin}"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

raw_arch="${BUILD_ARCH:-$(uname -m)}"
case "$raw_arch" in
  x86_64 | amd64) arch="amd64" ;;
  aarch64 | arm64) arch="arm64" ;;
  armv7l | armv7) arch="armv7" ;;
  *)
    echo "install-latest-release.sh: unsupported architecture: $raw_arch" >&2
    exit 1
    ;;
esac

case "$os" in
  linux | darwin) ;;
  *)
    echo "install-latest-release.sh: unsupported operating system: $os" >&2
    exit 1
    ;;
esac

if [ "$os" = "darwin" ] && [ "$arch" = "armv7" ]; then
  echo "install-latest-release.sh: darwin armv7 builds are not published" >&2
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  checksum_cmd=(sha256sum -c -)
elif command -v shasum >/dev/null 2>&1; then
  checksum_cmd=(shasum -a 256 -c -)
else
  echo "install-latest-release.sh: sha256sum or shasum is required" >&2
  exit 1
fi

release_name="baillconnect-to-mqtt-$os-$arch"
case "$version" in
  latest)
    download_base="https://github.com/${repo}/releases/latest/download"
    ;;
  v*)
    download_base="https://github.com/${repo}/releases/download/${version}"
    ;;
  *)
    download_base="https://github.com/${repo}/releases/download/v${version}"
    ;;
esac

curl -fsSL -L -o "$tmp_dir/checksums.txt" "$download_base/checksums.txt"
curl -fsSL -L -o "$tmp_dir/$release_name" "$download_base/$release_name"

expected="$(grep " ${release_name}$" "$tmp_dir/checksums.txt" | cut -d ' ' -f 1)"
if [ -z "$expected" ]; then
  echo "install-latest-release.sh: checksum not found for $release_name" >&2
  exit 1
fi

printf '%s  %s\n' "$expected" "$tmp_dir/$release_name" | "${checksum_cmd[@]}"

mkdir -p "$bin_home"
install -m755 "$tmp_dir/$release_name" "$bin_home/baillconnect-to-mqtt"
