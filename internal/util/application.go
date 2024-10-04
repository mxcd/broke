package util

import (
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func SetLogLevel(c *cli.Context) {
	if c.Bool("verbose") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if c.Bool("very-verbose") {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}

func PrintLogo(c *cli.Context) {
	if c.Bool("no-logo") {
		return
	}
	println("__       __               __")
	println("\\ \\     / /_  _________  / /_____")
	println(" \\ \\   / __ \\/ ___/ __ \\/ //_/ _ \\")
	println(" / /  / /_/ / /  / /_/ / ,< /  __/")
	println("/_/  /_.___/_/   \\____/_/|_|\\___/")
	println("")
}
