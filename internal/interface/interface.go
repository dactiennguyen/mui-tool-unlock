package interfaces

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"muitoolunlock/internal/auth"
	"muitoolunlock/internal/colors"
	"muitoolunlock/internal/device"
	"muitoolunlock/internal/storage"
	"muitoolunlock/internal/unlock"
)

// RunInteractiveUnlock runs the interactive unlock process
func RunInteractiveUnlock(fastbootPath string) {
	fmt.Println(colors.Header("🔐 Interactive Xiaomi Device Unlock"))

	// Load existing data
	data := storage.LoadUnlockData()

	// Get account info
	if data.User == "" {
		fmt.Print(colors.Email("Xiaomi Account (ID/Email/Phone): "))
		reader := bufio.NewReader(os.Stdin)
		account, _ := reader.ReadString('\n')
		data.User = strings.TrimSpace(account)
		storage.SaveUnlockData(data)
		fmt.Println(colors.Save("Account saved"))
	}

	if data.Password == "" {
		fmt.Print(colors.Lock("Enter password: "))
		reader := bufio.NewReader(os.Stdin)
		password, _ := reader.ReadString('\n')
		data.Password = strings.TrimSpace(password)
		storage.SaveUnlockData(data)
		fmt.Println(colors.Save("Password saved"))
	}

	// Get web browser ID if not exists (similar to Python wb_id flow)
	if data.WbID == "" {
		fmt.Println(colors.Section("🌐 Web Authentication Required"))
		fmt.Println(colors.Notice("If logged in with any account in your browser,"))
		fmt.Println(colors.Notice("please log out before continuing."))
		fmt.Println(colors.Rocket("Opening Xiaomi authentication page automatically..."))

		// Get device ID from web authentication (auto-open browser)
		deviceID := auth.GetWebBrowserID()
		if deviceID == "" {
			fmt.Println(colors.Error("Web authentication failed"))
			return
		}
		data.WbID = deviceID
		storage.SaveUnlockData(data)
	}

	// Authenticate with Xiaomi
	fmt.Println(colors.Progress("Authenticating with Xiaomi servers..."))
	authData, err := auth.AuthenticateXiaomi(data.User, data.Password, data.WbID)
	if err != nil {
		fmt.Println(colors.Error(fmt.Sprintf("Authentication failed: %v", err)))
		return
	}

	fmt.Printf("%s %s\n", colors.Success("Authentication successful! Account ID:"), colors.BoldText(authData.UserID))

	// Save login success
	if data.Login != "ok" {
		data.Login = "ok"
		data.UID = authData.UserID
		storage.SaveUnlockData(data)
		fmt.Println(colors.Save("Login saved."))
	}

	// Get device info
	fmt.Println(colors.Section("📱 Device Information"))
	fmt.Println(colors.Warning("Ensure you're in Bootloader mode (fastboot mode)"))

	deviceInfo := device.GetDeviceInfo(fastbootPath)
	if deviceInfo == nil {
		fmt.Println(colors.Error("Failed to get device info. Please ensure device is in fastboot mode."))
		return
	}

	device.DisplayDeviceInfo(deviceInfo)

	// Confirm unlock
	fmt.Print(colors.Warning("Are you sure you want to unlock this device? (y/N): "))
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println(colors.Error("Unlock cancelled"))
		return
	}

	// Perform real unlock with API
	unlock.PerformUnlock(deviceInfo, authData, fastbootPath)
}

// ProcessDirectUnlock processes direct unlock with account/password parameters
func ProcessDirectUnlock(account, password, fastbootPath string) {
	fmt.Printf("%s %s\n", colors.Progress("Processing unlock for account:"), colors.BoldText(account))
	fmt.Println(colors.Error("Direct unlock requires web authentication. Please use interactive mode."))
	fmt.Println(colors.Info("💡 Run: mui-tool-unlock-terminal (without flags)"))
}

// RunDeviceMode runs device information mode
func RunDeviceMode(fastbootPath string) {
	fmt.Println(colors.Header("📱 Device Information Mode"))

	deviceInfo := device.GetDeviceInfo(fastbootPath)
	if deviceInfo == nil {
		fmt.Println(colors.Error("No device found. Please ensure device is connected and in fastboot mode."))
		return
	}

	device.DisplayDeviceInfo(deviceInfo)
}
