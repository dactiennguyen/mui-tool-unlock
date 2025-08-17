package device

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"muitoolunlock/internal/types"
)

// GetDeviceInfo retrieves device information using fastboot
func GetDeviceInfo(fastbootPath string) *types.DeviceInfo {
	// Try to get device info using fastboot commands
	deviceInfo := &types.DeviceInfo{}

	fmt.Print("⏳ Waiting for device")
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Println()

	// Get unlocked status
	fmt.Print("📋 Fetching 'unlocked' — please wait...")
	if output := RunFastbootCommand(fastbootPath, "getvar", "unlocked"); output != "" {
		deviceInfo.Unlocked = output
		fmt.Print("\r\033[K") // Clear line
	} else {
		fmt.Print("\r\033[K")
		return nil
	}

	// Get product info
	fmt.Print("📋 Fetching 'product' — please wait...")
	if output := RunFastbootCommand(fastbootPath, "getvar", "product"); output != "" {
		deviceInfo.Product = output
		fmt.Print("\r\033[K")
	}

	// Try to get token (determines SoC type)
	fmt.Print("📋 Fetching 'token' — please wait...")
	if token := RunFastbootCommand(fastbootPath, "oem", "get_token"); token != "" {
		deviceInfo.Token = token
		deviceInfo.SoC = "Mediatek"
	} else if token := RunFastbootCommand(fastbootPath, "getvar", "token"); token != "" {
		deviceInfo.Token = token
		deviceInfo.SoC = "Qualcomm"
	}
	fmt.Print("\r\033[K")

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
	fmt.Println("\n📱 Device Information:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🔓 Unlocked: %s\n", info.Unlocked)
	fmt.Printf("📱 Product: %s\n", info.Product)
	fmt.Printf("🔧 SoC: %s\n", info.SoC)
	if info.Token != "" {
		fmt.Printf("🔑 Token: %s...\n", info.Token[:min(len(info.Token), 20)])
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
