#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
version="${VERSION:-1.1.1}"
dist="$root/dist"
mkdir -p "$dist"

(
  cd "$root"
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H windowsgui" -o "$dist/Lookaway.exe" .
)

cat > "$dist/README-Windows.txt" << EOF
Lookaway ${version}

Run Lookaway.exe. Tray icon by the clock. Right-click for Break now, Pause, Skip, Quit.

SmartScreen: More info -> Run anyway.

https://github.com/adesokanayo/lookaway
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
