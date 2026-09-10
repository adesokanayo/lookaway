//go:build darwin

package main

import (
	"github.com/progrium/darwinkit/macos"
	"github.com/progrium/darwinkit/macos/appkit"
)

func main() {
	macos.RunApp(func(app appkit.Application, delegate *appkit.ApplicationDelegate) {
		startApp(app)
	})
}
