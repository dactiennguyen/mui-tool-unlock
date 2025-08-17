//go:build terminal

package main

import (
	"flag"
	"fmt"

	interfaces "muitoolunlock/internal/interface"
	"muitoolunlock/internal/platform"
	"muitoolunlock/internal/types"
)

func main() {
	// CLI flags
	var (
		version    = flag.Bool("version", false, "Show version information")
		help       = flag.Bool("help", false, "Show help information")
		unlock     = flag.Bool("unlock", false, "Start unlock process")
		account    = flag.String("account", "", "Xiaomi account (email/phone/ID)")
		password   = flag.String("password", "", "Account password")
		deviceMode = flag.Bool("device", false, "Interactive device unlock mode")
	)

	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Printf("MUI Tool Unlock CLI v%s\n", types.AppVersion)
		fmt.Println("Built with Go - Xiaomi Device Unlocker")
		return
	}

	// Handle help flag
	if *help {
		printHelp()
		return
	}

	// Start CLI interface
	fmt.Println("🔓 MUI Tool Unlock - Xiaomi Device Unlocker")
	fmt.Println("============================================")
	fmt.Printf("[V%s] For issues: github.com/offici5l/MiUnlockTool\n", types.AppVersion)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Setup platform tools first
	fastbootPath := platform.Setup()
	if fastbootPath == "" {
		fmt.Println("❌ Failed to setup fastboot tools")
		return
	}

	if *unlock && *account != "" && *password != "" {
		// Direct unlock mode with parameters
		interfaces.ProcessDirectUnlock(*account, *password, fastbootPath)
	} else if *deviceMode {
		// Device interaction mode
		interfaces.RunDeviceMode(fastbootPath)
	} else {
		// Interactive mode
		interfaces.RunInteractiveUnlock(fastbootPath)
	}
}

func printHelp() {
	fmt.Println("MUI Tool Unlock - Xiaomi Device Unlocker")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  mui-tool-unlock-terminal [flags]")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --version                Show version information")
	fmt.Println("  --help                   Show this help message")
	fmt.Println("  --unlock                 Start unlock process")
	fmt.Println("  --account <account>      Xiaomi account (email/phone/ID)")
	fmt.Println("  --password <password>    Account password")
	fmt.Println("  --device                 Interactive device unlock mode")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  mui-tool-unlock-terminal")
	fmt.Println("  mui-tool-unlock-terminal --unlock --account user@mi.com --password mypass")
	fmt.Println("  mui-tool-unlock-terminal --device")
	fmt.Println("  mui-tool-unlock-terminal --version")
}
