package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	StarterFoo    string
	StarterFooKey = "STARTER_FOO"

	StarterBar    int
	StarterBarKey = "STARTER_BAR"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})

	godotenv.Load()
	log.Info().Msg("env loaded")

	flag.StringVar(&StarterFoo, "foo", lookupEnvOrString(StarterFooKey, "abc"), "foo(string)")
	flag.IntVar(&StarterBar, "bar", lookupEnvOrInt(StarterBarKey, 1), "bar(int)")

	flag.Parse()
	log.Info().Msg("flag parsed")

}

// lookupEnvOrString lookup environment value with key or fallback with default value
func lookupEnvOrString(key string, val string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return val
	}
	return v
}

// lookupEnvOrInt lookup environment value with key or fallback with default value
func lookupEnvOrInt(key string, val int) int {
	strv, ok := os.LookupEnv(key)
	if !ok {
		return val
	}
	v, err := strconv.Atoi(strv)
	if err != nil {
		log.Fatal().Msgf("lookupEnvOrInt[%s]: %v", key, err)
	}
	return v
}
