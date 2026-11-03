package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/racktrace/freshloom/internal/config"
	"github.com/racktrace/freshloom/internal/step"
	"github.com/racktrace/freshloom/internal/system"
)

func main() {
	cfgPath := flag.String("config", "freshloom.yaml", "path to the config file")
	dryRun := flag.Bool("dry-run", false, "show what would change without changing anything")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	apt := &step.Apt{Packages: cfg.Packages.Apt, Exec: system.ShellExecutor{}}

	done, err := apt.Check()
	switch {
	case err != nil:
		fmt.Fprintf(os.Stderr, "[fail] %s: %v\n", apt.Name(), err)
		os.Exit(1)
	case done:
		fmt.Printf("[skip] %s\n", apt.Name())
	case *dryRun:
		fmt.Printf("[plan] %s\n", apt.Name())
	default:
		if err := apt.Apply(); err != nil {
			fmt.Fprintf(os.Stderr, "[fail] %s: %v\n", apt.Name(), err)
			os.Exit(1)
		}
		fmt.Printf("[ok]   %s\n", apt.Name())
	}
}