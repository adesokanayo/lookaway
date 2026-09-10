#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
src="$root/packaging/icon.png"
setdir="$root/packaging/AppIcon.iconset"
icns="$root/packaging/AppIcon.icns"

rm -rf "$setdir"
mkdir -p "$setdir"

# Source may be JPEG saved as .png; normalize first.
work="$(mktemp -t lookaway-icon).png"
sips -s format png "$src" --out "$work" >/dev/null

for size in 16 32 64 128 256 512 1024; do
  sips -z "$size" "$size" "$work" --out "$setdir/icon_${size}x${size}.png" >/dev/null
done
cp "$setdir/icon_32x32.png" "$setdir/icon_16x16@2x.png"
cp "$setdir/icon_64x64.png" "$setdir/icon_32x32@2x.png"
cp "$setdir/icon_256x256.png" "$setdir/icon_128x128@2x.png"
cp "$setdir/icon_512x512.png" "$setdir/icon_256x256@2x.png"
cp "$setdir/icon_1024x1024.png" "$setdir/icon_512x512@2x.png"
rm -f "$setdir/icon_64x64.png"

iconutil -c icns "$setdir" -o "$icns"
rm -f "$work"
echo "Wrote $icns"
