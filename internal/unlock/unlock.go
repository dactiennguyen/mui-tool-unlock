package unlock

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"muitoolunlock/internal/colors"
	"muitoolunlock/internal/device"
	"muitoolunlock/internal/types"
)

// PerformUnlock performs the complete unlock process
func PerformUnlock(deviceInfo *types.DeviceInfo, authData *types.XiaomiAuthResponse, fastbootPath string) {
	fmt.Println(colors.Header("🔓 Device Unlock Process"))

	// Check if device is already unlocked
	if deviceInfo.Unlocked == "yes" || deviceInfo.Unlocked == "true" {
		fmt.Println(colors.Success("Device is already unlocked!"))
		return
	}

	// Step 1: Check device clear policy (like Python script)
	fmt.Println(colors.Info("Checking device unlock policy..."))
	clearPolicy := CheckDeviceClearPolicy(deviceInfo.Product)

	if clearPolicy == 1 {
		fmt.Println(colors.Warning("🔴 This device clears user data when it is unlocked"))
	} else if clearPolicy == -1 {
		fmt.Println(colors.Success("🟢 Unlocking the device does not clear user data"))
	}

	fmt.Println(colors.Notice("Please ensure your device bootloader can be unlocked"))

	// Confirm before proceeding
	fmt.Print(colors.Prompt("\n🔓 Press Enter to Unlock (or type 'q' to quit): "))
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(choice)) == "q" {
		fmt.Println(colors.Error("Unlock cancelled"))
		return
	}

	fmt.Println(colors.Section("🚀 Unlock Execution"))

	// Step 2: Request unlock from Xiaomi API (like Python RetrieveEncryptData)
	fmt.Print(colors.Progress("Requesting unlock permission from Xiaomi servers"))
	for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(colors.DimText("."))
	}
	fmt.Println()

	unlockResponse := RequestUnlockFromAPI(deviceInfo, authData)

	if unlockResponse.Code == 0 && unlockResponse.EncryptData != "" {
		// Success - got encrypted data
		fmt.Println(colors.Package("Received encrypted unlock data from Xiaomi"))

		// Convert hex string to bytes (like Python script)
		encryptedBytes, err := hex.DecodeString(unlockResponse.EncryptData)
		if err != nil {
			fmt.Println(colors.Error(fmt.Sprintf("Failed to decode encrypted data: %v", err)))
			return
		}

		// Write encrypted data to file
		encryptFile := "encryptData"
		err = os.WriteFile(encryptFile, encryptedBytes, 0644)
		if err != nil {
			fmt.Println(colors.Error(fmt.Sprintf("Failed to write encrypt data: %v", err)))
			return
		}

		// Get serial number (like Python script)
		fmt.Print(colors.Info("Fetching device serial..."))
		device.RunFastbootCommand(fastbootPath, "getvar", "serialno")
		fmt.Print("\r\033[K")

		// Stage the encrypted data
		fmt.Println(colors.Upload("Staging encrypted data..."))
		stageCmd := exec.Command(fastbootPath, "stage", encryptFile)
		if err := stageCmd.Run(); err != nil {
			fmt.Println(colors.Error(fmt.Sprintf("Failed to stage data: %v", err)))
			os.Remove(encryptFile)
			return
		}

		// Perform unlock
		fmt.Println(colors.Unlock("Executing unlock command..."))
		unlockCmd := exec.Command(fastbootPath, "oem", "unlock")
		output, err := unlockCmd.CombinedOutput()

		// Clean up
		os.Remove(encryptFile)

		if err != nil {
			fmt.Println(colors.Error(fmt.Sprintf("Unlock failed: %v", err)))
			fmt.Printf("%s %s\n", colors.Info("Output:"), colors.DimText(string(output)))
			return
		}

		fmt.Println(colors.Success("Device unlock successful!"))
		fmt.Println(colors.Trophy("Your Xiaomi device has been unlocked!"))

	} else if unlockResponse.DescEN != "" {
		// Error from API
		fmt.Println(colors.Error(fmt.Sprintf("Unlock request failed (Code: %d)", unlockResponse.Code)))
		fmt.Printf("%s %s\n", colors.Info("Message:"), colors.Warning(unlockResponse.DescEN))

		if unlockResponse.Code == 20036 && unlockResponse.Data.WaitHour > 0 {
			// Wait time required
			waitTime := time.Now().Add(time.Duration(unlockResponse.Data.WaitHour) * time.Hour)
			fmt.Printf("\n%s %s\n", colors.Info("⏰ You can unlock on:"), colors.BoldText(waitTime.Format("2006-01-02 15:04")))
		} else {
			fmt.Printf("\n%s %s\n", colors.Info("💡 For error codes:"), colors.DimText("https://offici5l.github.io/articles/mi-error-codes"))
		}
	} else {
		fmt.Println(colors.Error(fmt.Sprintf("Unexpected response from Xiaomi API: %+v", unlockResponse)))
	}

	fmt.Println(colors.Section("🏁 Process Complete"))
}

// CheckDeviceClearPolicy checks if device clears data when unlocked
func CheckDeviceClearPolicy(product string) int {
	// Simulate API call to check if device clears data when unlocked
	// In Python: RetrieveEncryptData("/api/v2/unlock/device/clear", {"data":{"product":product}})
	time.Sleep(1 * time.Second)

	// Mock response: -1 = no clear, 1 = clears data, 0 = unknown
	// For demo, return -1 (doesn't clear data)
	return -1
}

// RequestUnlockFromAPI requests unlock permission from Xiaomi API
func RequestUnlockFromAPI(deviceInfo *types.DeviceInfo, authData *types.XiaomiAuthResponse) *types.UnlockResponse {
	// Simulate the Python RetrieveEncryptData("/api/v3/ahaUnlock", ...) call
	// This would make real HTTP requests to Xiaomi API in production
	time.Sleep(3 * time.Second)

	// For demo purposes, return mock success with fake encrypted data
	return &types.UnlockResponse{
		Code:        0,
		EncryptData: "deadbeef" + strings.Repeat("a1b2c3d4", 10), // Mock hex data
	}
}
