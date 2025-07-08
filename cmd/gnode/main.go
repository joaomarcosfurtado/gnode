package main

import (
	"fmt"
	"os"

	"github.com/joaomarcosfurtado/gnode/internal/manager"
	"github.com/joaomarcosfurtado/gnode/pkg/config"
)

func printUsage() {
	fmt.Println("Usage: gnode <command> [args]")
	fmt.Println("\nCommands:")
	fmt.Println(" install <version>      Install some Node.js version")
	fmt.Println(" use <version>          Use some installed versions")
	fmt.Println(" list 					 List installed versions")
	fmt.Println(" list-remote 			 List versions availables to download")
	fmt.Println(" current				 Show current version")
	fmt.Println(" which  				 Show the executable path of Node.js")
	fmt.Println(" unistall <version>     Uninstall some Node.js version")
	fmt.Println(" help 					 Show this help")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Error loading configs: %v\n", err)
		os.Exit(1)
	}

	mgr, err := manager.NewManager(cfg)
	if err != nil {
		fmt.Printf("Error initializing manager: %v\n", err)
		os.Exit(1)
	}

	if err := mgr.Init(); err != nil {
		fmt.Printf("Initializing error")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Use: gnode install <version>")
			os.Exit(1)
		}

		if err := mgr.Install(os.Args[2]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "use":
		if len(os.Args) < 3 {
			fmt.Println("Use: gnode use <version>")
			os.Exit(1)
		}
		if err := mgr.Use(os.Args[2]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "list":
		if err := mgr.ListLocal(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "list-remote":
		if err := mgr.ListRemote(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "current":
		if err := mgr.ShowCurrent(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "which":
		if err := mgr.ShowWhich(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "uninstall":
		if len(os.Args) < 3 {
			fmt.Println("Use: gnode uninstall <version>")
			os.Exit(1)
		}

		if err := mgr.Uninstall(os.Args[2]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}
