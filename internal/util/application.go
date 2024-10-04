package util

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/urfave/cli/v2"
)

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("No .env file found")
	} else {
		log.Info().Msg("Loaded .env file")
	}
}

func SetLogLevel(c *cli.Context) {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Logger.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.RFC3339,
	}).With().Caller().Logger()

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

var cliContext *cli.Context

func SetCliContext(c *cli.Context) {
	cliContext = c
}
func GetCliContext() *cli.Context {
	if cliContext == nil {
		cliContext = GetDummyCliContext()
	}
	return cliContext
}

func GetDummyCliContext() *cli.Context {
	app := &cli.App{
		Name:  "broke",
		Usage: "IAM broker",
		Flags: []cli.Flag{},
	}
	return cli.NewContext(app, nil, nil)
}
