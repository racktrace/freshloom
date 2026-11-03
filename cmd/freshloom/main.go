package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/racktrace/freshloom/internal/config"
)

func main() {
	cfgPath := flag.String("config", "freshloom.yaml", "path to the config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", cfg)
}