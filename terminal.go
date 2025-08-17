//go:build terminal

package main

import (
	"flag"
	"fmt"

	"muitoolunlock/internal/colors"
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
		fmt.Printf("%s v%s\n", colors.BoldText("MUI Tool Unlock CLI"), colors.BoldText(types.AppVersion))
		fmt.Println(colors.DimText("Built with Go - Xiaomi Device Unlocker"))
		return
	}

	// Handle help flag
	if *help {
		printHelp()
		return
	}

	// Start CLI interface
	fmt.Println(colors.Rainbow("🔓 MUI Tool Unlock - Xiaomi Device Unlocker"))
	fmt.Println(colors.Gradient("============================================"))
	fmt.Printf("%s%s%s %s\n",
		colors.DimText("[V"), colors.BoldText(types.AppVersion), colors.DimText("] For issues:"),
		colors.UnderlineText("github.com/offici5l/MiUnlockTool"))
	fmt.Println(colors.Section("System Initialization"))

	// Setup platform tools first
	fastbootPath := platform.Setup()
	if fastbootPath == "" {
		fmt.Println(colors.Error("Failed to setup fastboot tools"))
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
	fmt.Println(colors.Header("MUI Tool Unlock - Xiaomi Device Unlocker"))
	fmt.Println()
	fmt.Println(colors.BoldText("Usage:"))
	fmt.Printf("  %s [flags]\n", colors.UnderlineText("mui-tool-unlock-terminal"))
	fmt.Println()
	fmt.Println(colors.BoldText("Flags:"))
	fmt.Printf("  %s                %s\n", colors.Info("--version"), colors.DimText("Show version information"))
	fmt.Printf("  %s                   %s\n", colors.Info("--help"), colors.DimText("Show this help message"))
	fmt.Printf("  %s                 %s\n", colors.Info("--unlock"), colors.DimText("Start unlock process"))
	fmt.Printf("  %s      %s\n", colors.Info("--account <account>"), colors.DimText("Xiaomi account (email/phone/ID)"))
	fmt.Printf("  %s    %s\n", colors.Info("--password <password>"), colors.DimText("Account password"))
	fmt.Printf("  %s                 %s\n", colors.Info("--device"), colors.DimText("Interactive device unlock mode"))
	fmt.Println()
	fmt.Println(colors.BoldText("Examples:"))
	fmt.Printf("  %s\n", colors.Success("mui-tool-unlock-terminal"))
	fmt.Printf("  %s\n", colors.Success("mui-tool-unlock-terminal --unlock --account user@mi.com --password mypass"))
	fmt.Printf("  %s\n", colors.Success("mui-tool-unlock-terminal --device"))
	fmt.Printf("  %s\n", colors.Success("mui-tool-unlock-terminal --version"))
}
