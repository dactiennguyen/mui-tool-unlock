package device

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"muitoolunlock/internal/colors"
	"muitoolunlock/internal/types"
)

// GetDeviceInfo retrieves device information using fastboot
func GetDeviceInfo(fastbootPath string) *types.DeviceInfo {
	// Try to get device info using fastboot commands
	deviceInfo := &types.DeviceInfo{}

	fmt.Print(colors.Progress("Waiting for device"))
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(colors.DimText("."))
	}
	fmt.Println()

	// Get unlocked status
	fmt.Print(colors.Info("Fetching 'unlocked' — please wait..."))
	if output := RunFastbootCommand(fastbootPath, "getvar", "unlocked"); output != "" {
		deviceInfo.Unlocked = output
		fmt.Print("\r\033[K") // Clear line
		fmt.Println(colors.Success("Retrieved unlock status"))
	} else {
		fmt.Print("\r\033[K")
		fmt.Println(colors.Error("Failed to get device info"))
		return nil
	}

	// Get product info
	fmt.Print(colors.Info("Fetching 'product' — please wait..."))
	if output := RunFastbootCommand(fastbootPath, "getvar", "product"); output != "" {
		deviceInfo.Product = output
		fmt.Print("\r\033[K")
		fmt.Println(colors.Success("Retrieved product info"))
	}

	// Try to get token (determines SoC type)
	fmt.Print(colors.Info("Fetching 'token' — please wait..."))
	if token := RunFastbootCommand(fastbootPath, "oem", "get_token"); token != "" {
		deviceInfo.Token = token
		deviceInfo.SoC = "Mediatek"
		fmt.Print("\r\033[K")
		fmt.Println(colors.Success("Retrieved Mediatek token"))
	} else if token := RunFastbootCommand(fastbootPath, "getvar", "token"); token != "" {
		deviceInfo.Token = token
		deviceInfo.SoC = "Qualcomm"
		fmt.Print("\r\033[K")
		fmt.Println(colors.Success("Retrieved Qualcomm token"))
	} else {
		fmt.Print("\r\033[K")
		fmt.Println(colors.Warning("Token not available"))
	}

	return deviceInfo
}

// RunFastbootCommand executes fastboot command and returns output
func RunFastbootCommand(cmd string, args ...string) string {
	// Try to execute fastboot command
	execCmd := exec.Command(cmd, args...)
	output, err := execCmd.CombinedOutput()

	if err != nil {
		return ""
	}

	// Parse fastboot output to extract variable value
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		if len(args) > 1 && strings.Contains(line, args[1]+":") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
		// For token commands, look for token data
		if strings.Contains(line, "token") && len(line) > 10 {
			return strings.TrimSpace(line)
		}
	}

	// For basic device detection
	if strings.Contains(outputStr, "fastboot") || strings.Contains(outputStr, "waiting") {
		return "detected"
	}

	return ""
}

// DisplayDeviceInfo prints device information to console
func DisplayDeviceInfo(info *types.DeviceInfo) {
	fmt.Println(colors.Header("📱 Device Information"))

	// Unlocked status with conditional coloring
	unlockStatus := info.Unlocked
	if unlockStatus == "yes" || unlockStatus == "true" {
		fmt.Printf("%s %s\n", colors.Unlock("Unlocked:"), colors.Success(unlockStatus))
	} else {
		fmt.Printf("%s %s\n", colors.Lock("Locked:"), colors.Error(unlockStatus))
	}

	fmt.Printf("%s %s\n", colors.Device("Product:"), colors.BoldText(info.Product))
	fmt.Printf("%s %s\n", colors.Tool("SoC:"), colors.BoldText(info.SoC))

	if info.Token != "" {
		tokenDisplay := info.Token
		if len(tokenDisplay) > 20 {
			tokenDisplay = tokenDisplay[:20] + "..."
		}
		fmt.Printf("%s %s\n", colors.Key("Token:"), colors.DimText(tokenDisplay))
	} else {
		fmt.Printf("%s %s\n", colors.Key("Token:"), colors.Warning("Not available"))
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
