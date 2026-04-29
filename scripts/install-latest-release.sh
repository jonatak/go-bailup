#!/usr/bin/env bash
# Install the latest published bailup binary from GitHub Releases.

set -euo pipefail

repo="${REPO:-jonatak/go-bailup}"
arch="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')"
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
release_name="bailup-$os-$arch"
bin_home="${XDG_BIN_HOME:-$HOME/.local/bin}"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

if ! command -v python3 >/dev/null 2>&1; then
  echo "install-latest-release.sh: python3 is required (JSON parsing)" >&2
  exit 1
fi

mkdir -p "$tmp_dir"

api_url="https://api.github.com/repos/${repo}/releases/latest"
json="$(curl -fsSL -H "Accept: application/vnd.github+json" "$api_url")"

url="$(printf '%s' "$json" | python3 -c '
import json, sys
name = sys.argv[1]
data = json.load(sys.stdin)
for a in data.get("assets", []):
    if a.get("name") == name:
        print(a["browser_download_url"])
        sys.exit(0)
print("no release asset named:", repr(name), file=sys.stderr)
sys.exit(1)
' "$release_name")"

curl -fsSL -L -o "$tmp_dir/$release_name" "$url"

mkdir -p "$bin_home"
install -m755 "$tmp_dir/$release_name" "$bin_home/bailup"
