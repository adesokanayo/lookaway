# Lookaway

[![License: MIT](https://img.shields.io/badge/License-MIT-0F766E.svg)](LICENSE)

A free, **open-source** app for macOS and Windows that enforces the **20-20-20 rule** by blocking every screen and asking you to look into the distance.

## Install

Download a build from the [latest release](https://github.com/adesokanayo/lookaway/releases/latest).

### macOS

1. Download **Lookaway-1.1.0.dmg**.
2. Open the DMG and drag **Lookaway** into **Applications**.
3. Open Lookaway. An eye icon appears in the menu bar.

If macOS blocks the first launch, Control-click Lookaway → **Open** → **Open**. The Mac build is not Apple-notarized yet.

### Windows

1. Download **Lookaway-windows-amd64-1.1.0.zip**.
2. Unzip and double-click **Lookaway.exe**.
3. A tray icon appears near the clock. Right-click it for Start break, Pause, Skip, or Quit.

If SmartScreen says the app is unrecognized, choose **More info** → **Run anyway**. The Windows build is not Authenticode-signed yet.

Pin the exe to startup if you want it every time you log in: `shell:startup` in Explorer’s address bar, then drop a shortcut to Lookaway.exe.

## Why this is open source

Lookaway is released under the [MIT License](LICENSE). Anyone can use, study, change, and share it.

That matters for an app that **takes over your screen**:

- **Trust.** You can read the Go source and see that the overlay is a timer, not a keylogger, ad injector, or telemetry collector. Nothing phones home.
- **No paywall on a health habit.** The 20-20-20 rule should not be locked behind a subscription.
- **Fix and fork.** If you want a 15-minute interval, a different message, or Linux support, you can change it or send a pull request.
- **Longevity.** If the original maintainer stops, the code stays public.
- **Reuse.** The Mac overlay uses [DarwinKit](https://github.com/progrium/darwinkit). The Windows overlay uses the Win32 API (`CreateWindowEx`, `HWND_TOPMOST`, tray `Shell_NotifyIcon`) with no extra GUI toolkit.

## Why 20-20-20

The American Optometric Association and American Academy of Ophthalmology recommend this habit for digital eye strain:

1. Every **20 minutes** of screen time
2. Look at something about **20 feet** (6 meters) away
3. For at least **20 seconds**

Near focus keeps the ciliary muscles contracted. A brief distant gaze lets them relax. Screen use also drops blink rate from about 15 blinks a minute to 4–6, which dries the eyes. The overlay reminds you to blink fully during the break.

A 2023 *Contact Lens and Anterior Eye* study found the rule increased break frequency and reduced digital eye strain and dry-eye *symptoms*. It did not erase dry-eye signs on its own, so this app is a reminder, not a medical treatment.

## What the app does

- Lives in the macOS menu bar or the Windows system tray
- Counts down 20 minutes of work
- Then covers **all displays** with a topmost overlay
- Asks you to look 20 feet away while a 20-second timer runs
- Restarts automatically after the break
- Lets you start a break early, pause the timer, or skip if you must

On a Mac the overlay is an AppKit panel at screensaver window level. On Windows it is a topmost popup covering the virtual desktop (every monitor).

## Menu

- **Next break in…** — time until the overlay
- **Start break now** — block the screen immediately
- **Pause timer / Resume timer**
- **Skip current break** — also on the overlay (Esc on Windows)
- **Quit Lookaway**

## Build from source

You need Go 1.21+. On a Mac, also install the Xcode command line tools.

```bash
git clone https://github.com/adesokanayo/lookaway.git
cd lookaway
go test ./...
go run .
```

Shorter intervals while you test:

```bash
LOOKAWAY_WORK=30s LOOKAWAY_BREAK=20s go run .
```

Installers:

```bash
./scripts/package-dmg.sh          # Mac: dist/Lookaway-1.1.0.dmg
./scripts/package-windows.sh      # Windows: dist/Lookaway-windows-amd64-1.1.0.zip
```

The Windows zip can be built from macOS; it does not need CGO.

## Limits

The lock screen, Task Manager, Force Quit, and some system UI can still appear above the overlay. That is intentional so you are never locked out of the computer.

This is not a substitute for an eye exam, glasses, lighting changes, or less total screen time.

## Contributing

Issues and pull requests are welcome. Please keep changes focused: timer behavior, overlay safety (never trap the user), and packaging.

## License

[MIT](LICENSE)
