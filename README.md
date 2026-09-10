# Lookaway

Every 20 minutes, the screen locks so you look ~20 feet away for 20 seconds.

I built this because I sit on a laptop all day and my eyes get fried. Notifications were too easy to swipe away. I wanted something that actually takes the screen for 20 seconds.

The 20-20-20 rule: every 20 minutes, look 20 feet away for 20 seconds. AOA and AAO recommend it for digital eye strain. A 2023 study found it cut symptoms; it is not a cure.

[Download](https://github.com/adesokanayo/lookaway/releases/latest)

**Mac** — open the DMG, drag to Applications. If macOS blocks it: Control-click → Open.

**Windows** — unzip, run `Lookaway.exe`. Icon sits in the tray. SmartScreen: More info → Run anyway.

```bash
git clone https://github.com/adesokanayo/lookaway.git
cd lookaway
go test ./...
go run .
```

```bash
LOOKAWAY_WORK=30s LOOKAWAY_BREAK=20s go run .
./scripts/package-dmg.sh
./scripts/package-windows.sh
```

MIT. PRs welcome.
