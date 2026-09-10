#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
version="${VERSION:-1.1.0}"
dist="$root/dist"
stage="$dist/dmg"
app="$stage/Lookaway.app"

rm -rf "$stage"
mkdir -p "$app/Contents/MacOS" "$app/Contents/Resources"

(
  cd "$root"
  GOOS=darwin CGO_ENABLED=1 go build -ldflags="-s -w" -o "$app/Contents/MacOS/Lookaway" .
)

cp "$root/packaging/Info.plist" "$app/Contents/Info.plist"
ln -s /Applications "$stage/Applications"

if command -v codesign >/dev/null; then
  codesign --force --deep --sign - "$app"
fi

hdiutil create \
  -volname "Lookaway" \
  -srcfolder "$stage" \
  -ov \
  -format UDZO \
  "$dist/Lookaway-${version}.dmg"

echo "Created $dist/Lookaway-${version}.dmg"
