package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/joaomarcosfurtado/gnode/internal/manager"
	"github.com/joaomarcosfurtado/gnode/pkg/config"
)

func printUsage() {
	fmt.Println("Usage: gnode <command> [args]")
	fmt.Println("\nCommands:")
	fmt.Println(" install <version>     Install some Node.js version")
	fmt.Println(" use <version>         Use some installed version")
	fmt.Println(" list                  List installed versions")
	fmt.Println(" list-remote           List versions available to download")
	fmt.Println(" current               Show current version")
	fmt.Println(" which                 Show the executable path of Node.js")
	fmt.Println(" uninstall <version>   Uninstall some Node.js version")
	fmt.Println(" status                Show gnode status")
	fmt.Println(" help                  Show this help")

	if runtime.GOOS == "windows" {
		fmt.Println("\nWindows specific:")
		fmt.Println(" - First 'gnode use' will configure PATH automatically")
		fmt.Println(" - After that, it works just like nvm-windows")
	}
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
		fmt.Printf("Error initializing directories: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gnode install <version>")
			os.Exit(1)
		}
		if err := mgr.Install(os.Args[2]); err != nil {
			fmt.Printf("Error installing: %v\n", err)
			os.Exit(1)
		}
	case "use":
		if len(os.Args) < 3 {
			fmt.Println("Use: gnode use <version>")
			os.Exit(1)
		}
		printEnv := false
		if len(os.Args) > 3 && os.Args[3] == "--print-env" {
			printEnv = true
		}
		if err := mgr.Use(os.Args[2], printEnv); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "list":
		if len(os.Args) > 2 {
			fmt.Println("Command 'list' does not accept arguments")
			os.Exit(1)
		}
		if err := mgr.ListLocal(); err != nil {
			fmt.Printf("Error listing: %v\n", err)
			os.Exit(1)
		}
	case "list-remote":
		if len(os.Args) > 2 {
			fmt.Println("Command 'list-remote' does not accept arguments")
			os.Exit(1)
		}
		if err := mgr.ListRemote(); err != nil {
			fmt.Printf("Error listing remote: %v\n", err)
			os.Exit(1)
		}
	case "current":
		if len(os.Args) > 2 {
			fmt.Println("Command 'current' does not accept arguments")
			os.Exit(1)
		}
		if err := mgr.ShowCurrent(); err != nil {
			fmt.Printf("Error showing current version: %v\n", err)
			os.Exit(1)
		}
	case "which":
		if len(os.Args) > 2 {
			fmt.Println("Command 'which' does not accept arguments")
			os.Exit(1)
		}
		if err := mgr.ShowWhich(); err != nil {
			fmt.Printf("Error showing path: %v\n", err)
			os.Exit(1)
		}
	case "uninstall":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gnode uninstall <version>")
			os.Exit(1)
		}
		if err := mgr.Uninstall(os.Args[2]); err != nil {
			fmt.Printf("Error uninstalling: %v\n", err)
			os.Exit(1)
		}
	case "status":
		if err := mgr.Status(); err != nil {
			fmt.Printf("Error checking status: %v\n", err)
			os.Exit(1)
		}
	case "init":
		if err := mgr.Init(); err != nil {
			fmt.Printf("Error setting up environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("gnode environment initialized successfully!")
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}
