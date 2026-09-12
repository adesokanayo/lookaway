#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
version="${VERSION:-1.1.3}"
dist="$root/dist"
stage="$dist/stage"
app="$stage/Lookaway.app"
rw="$dist/Lookaway.rw.dmg"
final="$dist/Lookaway-${version}.dmg"

"$root/scripts/make-icns.sh"

rm -rf "$stage" "$rw" "$final"
mkdir -p "$app/Contents/MacOS" "$app/Contents/Resources"

(
  cd "$root"
  GOOS=darwin CGO_ENABLED=1 go build -ldflags="-s -w" -o "$app/Contents/MacOS/Lookaway" .
)

cp "$root/packaging/Info.plist" "$app/Contents/Info.plist"
cp "$root/packaging/AppIcon.icns" "$app/Contents/Resources/AppIcon.icns"
echo -n 'APPL????' > "$app/Contents/PkgInfo"

cat > "$stage/Install Lookaway.command" << 'EOF'
#!/bin/bash
set -e
cd "$(dirname "$0")"
killall Lookaway 2>/dev/null || true
sleep 1
xattr -cr "Lookaway.app" || true
rm -rf /Applications/Lookaway.app
ditto "Lookaway.app" /Applications/Lookaway.app
xattr -cr /Applications/Lookaway.app || true
open /Applications/Lookaway.app
echo "Installed. Look for the eye in the menu bar."
EOF
chmod +x "$stage/Install Lookaway.command"

cat > "$stage/README.txt" << 'EOF'
Lookaway

1. Double-click "Install Lookaway.command"
   or drag Lookaway.app to Applications.

2. If macOS says it cannot verify the app:
   System Settings → Privacy & Security → Open Anyway
   or Control-click the app → Open.

Unsigned builds show that warning. It is not a blank disk.
EOF

ln -s /Applications "$stage/Applications"

if command -v codesign >/dev/null; then
  codesign --force --deep --sign - "$app"
fi

# HFS+ so Finder lists the files. APFS srcfolder images often open as a blank window.
size=$(du -sm "$stage" | awk '{print $1 + 20}')
hdiutil create -ov -fs HFS+ -size "${size}m" -volname Lookaway -srcfolder "$stage" "$rw" >/dev/null

if mount=$(hdiutil attach -readwrite -noverify -nobrowse "$rw" 2>/dev/null | awk '/\/Volumes\//{print $NF}'); then
  osascript << EOF || true
tell application "Finder"
  tell disk "Lookaway"
    open
    set current view of container window to icon view
    set toolbar visible of container window to false
    set statusbar visible of container window to false
    set bounds of container window to {200, 120, 860, 540}
    set opts to icon view options of container window
    set arrangement of opts to not arranged
    set icon size of opts to 128
    set position of item "Lookaway.app" of container window to {150, 200}
    set position of item "Applications" of container window to {420, 200}
    set position of item "Install Lookaway.command" of container window to {150, 380}
    set position of item "README.txt" of container window to {420, 380}
    update without registering applications
    delay 1
    close
  end tell
end tell
EOF
  cp "$root/packaging/AppIcon.icns" "$mount/.VolumeIcon.icns"
  if command -v SetFile >/dev/null; then
    SetFile -c icnC "$mount/.VolumeIcon.icns" || true
    SetFile -a C "$mount" || true
  fi
  sync
  hdiutil detach "$mount" -quiet || true
fi

hdiutil convert "$rw" -format UDZO -imagekey zlib-level=9 -o "$final" >/dev/null
rm -f "$rw"

set_file_icon() {
  local icns="$1" dest="$2"
  local script
  script="$(mktemp -t lookaway-seticon).swift"
  cat > "$script" << 'SWIFT'
import AppKit
guard CommandLine.arguments.count >= 3,
      let img = NSImage(contentsOfFile: CommandLine.arguments[1]) else {
  exit(1)
}
NSWorkspace.shared.setIcon(img, forFile: CommandLine.arguments[2], options: [])
SWIFT
  swift "$script" "$icns" "$dest" || true
  rm -f "$script"
}
set_file_icon "$root/packaging/AppIcon.icns" "$final"

ditto -c -k --keepParent "$app" "$dist/Lookaway-${version}-mac.zip"
echo "Created $final"
echo "Created $dist/Lookaway-${version}-mac.zip"
