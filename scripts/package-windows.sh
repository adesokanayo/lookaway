#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
version="${VERSION:-1.1.0}"
dist="$root/dist"
mkdir -p "$dist"

(
  cd "$root"
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H windowsgui" -o "$dist/Lookaway.exe" .
)

cat > "$dist/README-Windows.txt" << EOF
Lookaway ${version} for Windows
================================

1. Copy Lookaway.exe wherever you like (Desktop or a folder).
2. Double-click Lookaway.exe. A tray icon appears near the clock.
3. Right-click the tray icon to start a break, pause, skip, or quit.

If Windows SmartScreen warns that the app is unrecognized, choose
More info -> Run anyway. The build is not Authenticode-signed yet.

Every 20 minutes Lookaway covers every monitor and asks you to look
about 20 feet away for 20 seconds.

Source: https://github.com/adesokanayo/lookaway
License: MIT
EOF

(
  cd "$dist"
  rm -f "Lookaway-windows-amd64-${version}.zip"
  if command -v zip >/dev/null; then
    zip -q "Lookaway-windows-amd64-${version}.zip" Lookaway.exe README-Windows.txt
  else
    python3 - << PY
import zipfile
z = zipfile.ZipFile("Lookaway-windows-amd64-${version}.zip", "w", zipfile.ZIP_DEFLATED)
z.write("Lookaway.exe")
z.write("README-Windows.txt")
z.close()
PY
  fi
)

echo "Created $dist/Lookaway-windows-amd64-${version}.zip"
