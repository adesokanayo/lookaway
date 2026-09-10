# Lookaway

[![License: MIT](https://img.shields.io/badge/License-MIT-0F766E.svg)](LICENSE)

A free, **open-source** macOS menu bar app that enforces the **20-20-20 rule** by blocking every screen and asking you to look into the distance.

## Install (DMG)

1. Download **Lookaway-1.0.0.dmg** from the [latest release](https://github.com/adesokanayo/lookaway/releases/latest).
2. Open the DMG and drag **Lookaway** into **Applications**.
3. Open Lookaway from Applications (or Spotlight). An eye icon appears in the menu bar.

The first launch may be blocked because the app is not yet Apple-notarized. In Finder, Control-click Lookaway → **Open** → **Open**. After that, macOS will remember your choice.

## Why this is open source

Lookaway is released under the [MIT License](LICENSE). Anyone can use, study, change, and share it.

That matters for an app that **takes over your screen**:

- **Trust.** You can read the Go source and see that the overlay is a timer, not a keylogger, ad injector, or telemetry collector. Nothing phones home.
- **No paywall on a health habit.** The 20-20-20 rule should not be locked behind a subscription.
- **Fix and fork.** If you want a 15-minute interval, a different message, or Linux support, you can change it or send a pull request.
- **Longevity.** If the original maintainer stops, the code stays public. The habit does not die with an App Store listing.
- **Reuse.** Other developers can learn how to build a native Mac overlay in Go with [DarwinKit](https://github.com/progrium/darwinkit).

Commercial screen-break tools often stay closed. That hides what they do while they have the highest window level on your machine. Open source flips that: the most privileged UI is also the most inspectable.

## Why 20-20-20

The American Optometric Association and American Academy of Ophthalmology recommend this habit for digital eye strain:

1. Every **20 minutes** of screen time
2. Look at something about **20 feet** (6 meters) away
3. For at least **20 seconds**

Near focus keeps the ciliary muscles contracted. A brief distant gaze lets them relax. Screen use also drops blink rate from about 15 blinks a minute to 4–6, which dries the eyes. The overlay reminds you to blink fully during the break.

A 2023 *Contact Lens and Anterior Eye* study found the rule increased break frequency and reduced digital eye strain and dry-eye *symptoms*. It did not erase dry-eye signs on its own, so this app is a reminder, not a medical treatment.

## What the app does

- Lives in the menu bar (no Dock icon)
- Counts down 20 minutes of work
- Then covers **all displays**, including the menu bar, with a high-level overlay
- Asks you to look 20 feet away while a 20-second timer runs
- Restarts automatically after the break
- Lets you start a break early, pause the timer, or skip if you must

The overlay uses native AppKit (`NSPanel` at `NSScreenSaverWindowLevel`) via DarwinKit, so it can sit above other apps — including fullscreen ones.

## Menu

- **Next break in…** — time until the overlay
- **Start break now** — block the screen immediately
- **Pause timer / Resume timer**
- **Skip current break** — also available on the overlay
- **Quit Lookaway**

## Build from source

You need Go 1.21+ and the Xcode command line tools.

```bash
git clone https://github.com/adesokanayo/lookaway.git
cd lookaway
go test ./...
go run .
```

Use **Start break now** to try the full-screen block immediately.

Shorter intervals while you test:

```bash
LOOKAWAY_WORK=30s LOOKAWAY_BREAK=20s go run .
```

Build the installer on a Mac:

```bash
./scripts/package-dmg.sh
# output: dist/Lookaway-1.0.0.dmg
```

## Limits

macOS can still show the lock screen, Force Quit (`Option-Command-Escape`), and some system UI above almost anything. That is intentional so you are never locked out of the computer.

This is not a substitute for an eye exam, glasses, lighting changes, or less total screen time.

## Contributing

Issues and pull requests are welcome. Please keep changes focused: timer behavior, overlay safety (never trap the user), and packaging.

## License

[MIT](LICENSE)
