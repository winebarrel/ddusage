package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/ddusage"
)

var version string

type options struct {
	ddusage.ClientOptions
	ddusage.PrintUsageSummaryOptions
}

func init() {
	log.SetFlags(0)
}

func main() {
	var cli struct {
		options
		Version kong.VersionFlag
	}

	kong.Parse(
		&cli,
		kong.Vars{"version": version},
	)

	client := ddusage.NewClient(&cli.ClientOptions)
	err := client.PrintUsageSummary(os.Stdout, &cli.PrintUsageSummaryOptions)

	if err != nil {
		log.Fatal(err)
	}
}
