package main

import (
	"codebase-app/cmd"
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure/config"
	"os"
	"strings"

	"flag"

	"github.com/rs/zerolog/log"
)

func main() {
	os.Args = initialize()

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	seedCmd := flag.NewFlagSet("seed", flag.ExitOnError)
	wsCmd := flag.NewFlagSet("ws", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Info().Msg("No command provided, defaulting to 'server'")
		cmd.RunServer(serverCmd, os.Args[1:])
		os.Exit(0)
	}

	switch os.Args[1] {
	case "seed":
		cmd.RunSeed(seedCmd, os.Args[2:])
	case "server":
		cmd.RunServer(serverCmd, os.Args[2:])
	case "ws":
		cmd.RunWebsocket(wsCmd, os.Args[2:])
	default:
		log.Info().Msg("Invalid command provided, defaulting to 'server' with provided flags")
		if os.Args[1][0] == '-' { // check if the first argument is a flag
			cmd.RunServer(serverCmd, os.Args[1:])
			os.Exit(0)
		}

		cmd.RunServer(serverCmd, os.Args[2:]) // default to server if invalid command and flags are provided
	}
}

func initialize() (newArgs []string) {
	configPath := flag.String("config_path", "./", "path to config file")
	configFilename := flag.String("config_filename", ".env", "config file name")
	flag.Parse()

	var logCfg string
	if *configPath == "./" {
		logCfg = *configPath + *configFilename
	} else {
		logCfg = *configPath + "/" + *configFilename
	}

	log.Info().Msgf("Initializing configuration with config: %s", logCfg)

	config.Configuration(
		config.WithPath(*configPath),
		config.WithFilename(*configFilename),
	).Initialize()

	adapter.Adapters = &adapter.Adapter{}

	for _, arg := range os.Args {
		if strings.Contains(arg, "config_path") || strings.Contains(arg, "config_filename") {
			continue
		}

		newArgs = append(newArgs, arg)
	}

	return newArgs
}
