package interfaces

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"muitoolunlock/internal/auth"
	"muitoolunlock/internal/device"
	"muitoolunlock/internal/storage"
	"muitoolunlock/internal/unlock"
)

// RunInteractiveUnlock runs the interactive unlock process
func RunInteractiveUnlock(fastbootPath string) {
	fmt.Println("\n🔐 Interactive Xiaomi Device Unlock")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Load existing data
	data := storage.LoadUnlockData()

	// Get account info
	if data.User == "" {
		fmt.Print("📧 Xiaomi Account (ID/Email/Phone): ")
		reader := bufio.NewReader(os.Stdin)
		account, _ := reader.ReadString('\n')
		data.User = strings.TrimSpace(account)
		storage.SaveUnlockData(data)
	}

	if data.Password == "" {
		fmt.Print("🔒 Enter password: ")
		reader := bufio.NewReader(os.Stdin)
		password, _ := reader.ReadString('\n')
		data.Password = strings.TrimSpace(password)
		storage.SaveUnlockData(data)
	}

	// Get web browser ID if not exists (similar to Python wb_id flow)
	if data.WbID == "" {
		fmt.Println("\n🌐 Web Authentication Required")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("📢 Notice: If logged in with any account in your browser,")
		fmt.Println("   please log out before continuing.")
		fmt.Print("\n🌐 Press Enter to open confirmation page...")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')

		// Get device ID from web authentication
		deviceID := auth.GetWebBrowserID()
		if deviceID == "" {
			fmt.Println("❌ Web authentication failed")
			return
		}
		data.WbID = deviceID
		storage.SaveUnlockData(data)
	}

	// Authenticate with Xiaomi
	fmt.Println("\n⏳ Authenticating with Xiaomi servers...")
	authData, err := auth.AuthenticateXiaomi(data.User, data.Password, data.WbID)
	if err != nil {
		fmt.Printf("❌ Authentication failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Authentication successful! Account ID: %s\n", authData.UserID)

	// Save login success
	if data.Login != "ok" {
		data.Login = "ok"
		data.UID = authData.UserID
		storage.SaveUnlockData(data)
		fmt.Println("💾 Login saved.")
	}

	// Get device info
	fmt.Println("\n📱 Fetching device information...")
	fmt.Println("⚠️  Ensure you're in Bootloader mode (fastboot mode)")

	deviceInfo := device.GetDeviceInfo(fastbootPath)
	if deviceInfo == nil {
		fmt.Println("❌ Failed to get device info. Please ensure device is in fastboot mode.")
		return
	}

	device.DisplayDeviceInfo(deviceInfo)

	// Confirm unlock
	fmt.Print("\n⚠️  Are you sure you want to unlock this device? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("❌ Unlock cancelled")
		return
	}

	// Perform real unlock with API
	unlock.PerformUnlock(deviceInfo, authData, fastbootPath)
}

// ProcessDirectUnlock processes direct unlock with account/password parameters
func ProcessDirectUnlock(account, password, fastbootPath string) {
	fmt.Printf("🔐 Processing unlock for account: %s\n", account)
	fmt.Println("❌ Direct unlock requires web authentication. Please use interactive mode.")
	fmt.Println("💡 Run: mui-tool-unlock-terminal (without flags)")
}

// RunDeviceMode runs device information mode
func RunDeviceMode(fastbootPath string) {
	fmt.Println("\n📱 Device Information Mode")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")

	deviceInfo := device.GetDeviceInfo(fastbootPath)
	if deviceInfo == nil {
		fmt.Println("❌ No device found. Please ensure device is connected and in fastboot mode.")
		return
	}

	device.DisplayDeviceInfo(deviceInfo)
}
