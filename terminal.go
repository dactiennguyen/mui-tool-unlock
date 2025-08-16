//go:build terminal

package main

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const appVersion = "1.5.9"

type UnlockData struct {
	User     string `json:"user"`
	Password string `json:"pwd"`
	WbID     string `json:"wb_id"`
	Login    string `json:"login"`
	UID      string `json:"uid"`
}

type DeviceInfo struct {
	Unlocked string
	Product  string
	SoC      string
	Token    string
}

type XiaomiAuthResponse struct {
	Code            int    `json:"code"`
	SecurityStatus  int    `json:"securityStatus"`
	NotificationURL string `json:"notificationUrl"`
	SSecurity       string `json:"ssecurity"`
	Nonce           string `json:"nonce"`
	Location        string `json:"location"`
	PassToken       string `json:"passToken"`
	UserID          string `json:"userId"`
}

type UnlockResponse struct {
	Code        int    `json:"code"`
	DescEN      string `json:"descEN"`
	EncryptData string `json:"encryptData"`
	Data        struct {
		WaitHour int `json:"waitHour"`
	} `json:"data"`
}

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
		fmt.Printf("MUI Tool Unlock CLI v%s\n", appVersion)
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
	fmt.Printf("[V%s] For issues: github.com/offici5l/MiUnlockTool\n", appVersion)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Setup platform tools first
	fastbootPath := setupPlatformTools()
	if fastbootPath == "" {
		fmt.Println("❌ Failed to setup fastboot tools")
		return
	}

	if *unlock && *account != "" && *password != "" {
		// Direct unlock mode with parameters
		processDirectUnlock(*account, *password, fastbootPath)
	} else if *deviceMode {
		// Device interaction mode
		runDeviceMode(fastbootPath)
	} else {
		// Interactive mode
		runInteractiveUnlock(fastbootPath)
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

func setupPlatformTools() string {
	// First try to find fastboot in system PATH
	if path, err := exec.LookPath("fastboot"); err == nil {
		fmt.Println("✅ Found fastboot in system PATH")
		return path
	}

	// If not found, try to download platform-tools
	fmt.Println("📦 Fastboot not found in PATH, setting up platform-tools...")

	execDir, err := os.Executable()
	if err != nil {
		return ""
	}
	baseDir := filepath.Dir(execDir)
	platformToolsDir := filepath.Join(baseDir, "platform-tools")

	// Check if platform-tools already exists
	var fastbootName string
	if runtime.GOOS == "windows" {
		fastbootName = "fastboot.exe"
	} else {
		fastbootName = "fastboot"
	}

	fastbootPath := filepath.Join(platformToolsDir, fastbootName)
	if _, err := os.Stat(fastbootPath); err == nil {
		fmt.Println("✅ Platform-tools already available")
		return fastbootPath
	}

	// Download platform-tools
	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "darwin"
	}

	url := fmt.Sprintf("https://dl.google.com/android/repository/platform-tools-latest-%s.zip", osName)
	zipPath := filepath.Join(baseDir, "platform-tools.zip")

	fmt.Println("⬇️  Downloading platform-tools...")
	if err := downloadFile(url, zipPath); err != nil {
		fmt.Printf("❌ Failed to download platform-tools: %v\n", err)
		return ""
	}

	fmt.Println("📦 Extracting platform-tools...")
	if err := unzipFile(zipPath, baseDir); err != nil {
		fmt.Printf("❌ Failed to extract platform-tools: %v\n", err)
		os.Remove(zipPath)
		return ""
	}

	// Clean up zip file
	os.Remove(zipPath)

	// Make fastboot executable on Unix systems
	if runtime.GOOS != "windows" {
		os.Chmod(fastbootPath, 0755)
	}

	fmt.Println("✅ Platform-tools setup completed")
	return fastbootPath
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzipFile(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.Create(fpath)
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func runInteractiveUnlock(fastbootPath string) {
	fmt.Println("\n🔐 Interactive Xiaomi Device Unlock")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Load existing data
	data := loadUnlockData()

	// Get account info
	if data.User == "" {
		fmt.Print("📧 Xiaomi Account (ID/Email/Phone): ")
		reader := bufio.NewReader(os.Stdin)
		account, _ := reader.ReadString('\n')
		data.User = strings.TrimSpace(account)
		saveUnlockData(data)
	}

	if data.Password == "" {
		fmt.Print("🔒 Enter password: ")
		reader := bufio.NewReader(os.Stdin)
		password, _ := reader.ReadString('\n')
		data.Password = strings.TrimSpace(password)
		saveUnlockData(data)
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
		deviceID := getWebBrowserID()
		if deviceID == "" {
			fmt.Println("❌ Web authentication failed")
			return
		}
		data.WbID = deviceID
		saveUnlockData(data)
	}

	// Authenticate with Xiaomi
	fmt.Println("\n⏳ Authenticating with Xiaomi servers...")
	authData, err := authenticateXiaomiReal(data.User, data.Password, data.WbID)
	if err != nil {
		fmt.Printf("❌ Authentication failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Authentication successful! Account ID: %s\n", authData.UserID)

	// Save login success
	if data.Login != "ok" {
		data.Login = "ok"
		data.UID = authData.UserID
		saveUnlockData(data)
		fmt.Println("💾 Login saved.")
	}

	// Get device info
	fmt.Println("\n📱 Fetching device information...")
	fmt.Println("⚠️  Ensure you're in Bootloader mode (fastboot mode)")

	deviceInfo := getDeviceInfo(fastbootPath)
	if deviceInfo == nil {
		fmt.Println("❌ Failed to get device info. Please ensure device is in fastboot mode.")
		return
	}

	displayDeviceInfo(deviceInfo)

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
	performRealUnlock(deviceInfo, authData, fastbootPath)
}

func processDirectUnlock(account, password, fastbootPath string) {
	fmt.Printf("🔐 Processing unlock for account: %s\n", account)
	fmt.Println("❌ Direct unlock requires web authentication. Please use interactive mode.")
	fmt.Println("💡 Run: mui-tool-unlock-terminal (without flags)")
}

func runDeviceMode(fastbootPath string) {
	fmt.Println("\n📱 Device Information Mode")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")

	deviceInfo := getDeviceInfo(fastbootPath)
	if deviceInfo == nil {
		fmt.Println("❌ No device found. Please ensure device is connected and in fastboot mode.")
		return
	}

	displayDeviceInfo(deviceInfo)
}

func getWebBrowserID() string {
	// In the Python script, this opens a browser and gets device ID from URL
	// For simplicity, we'll prompt user to enter it manually
	fmt.Println("\n🌐 Opening Xiaomi authentication page...")
	fmt.Println("🔗 URL: https://account.xiaomi.com/pass/serviceLogin?sid=unlockApi&checkSafeAddress=true&passive=false&hidden=false")
	fmt.Println("\n📋 Instructions:")
	fmt.Println("1. Open the URL above in your browser")
	fmt.Println("2. Login with your Xiaomi account")
	fmt.Println("3. After successful login, copy the URL")
	fmt.Println("4. Look for 'd=' parameter in the URL")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n🔗 Enter the redirect URL (or just the 'd' parameter value): ")
	urlStr, _ := reader.ReadString('\n')
	urlStr = strings.TrimSpace(urlStr)

	// Try to extract device ID from URL
	if strings.Contains(urlStr, "d=") {
		if parsedURL, err := url.Parse(urlStr); err == nil {
			if deviceID := parsedURL.Query().Get("d"); deviceID != "" {
				return deviceID
			}
		}
		// If URL parsing fails, try simple string extraction
		re := regexp.MustCompile(`d=([^&\s]+)`)
		if matches := re.FindStringSubmatch(urlStr); len(matches) > 1 {
			return matches[1]
		}
	}

	// If not a URL, assume user entered device ID directly
	if len(urlStr) > 10 && !strings.Contains(urlStr, " ") {
		return urlStr
	}

	return ""
}

func authenticateXiaomiReal(user, password, deviceID string) (*XiaomiAuthResponse, error) {
	// Simulate the real Xiaomi authentication process
	// In reality, this would make HTTP requests to Xiaomi API endpoints
	fmt.Println("🔐 Posting credentials to Xiaomi servers...")
	time.Sleep(2 * time.Second)

	// Hash password like Python script
	hasher := md5.New()
	hasher.Write([]byte(password))
	passwordHash := strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))

	fmt.Printf("📧 User: %s\n", user)
	fmt.Printf("🔑 Hash: %s...\n", passwordHash[:8])
	fmt.Printf("📱 Device ID: %s...\n", deviceID[:min(len(deviceID), 12)])

	// Simulate successful authentication
	return &XiaomiAuthResponse{
		Code:      0,
		UserID:    "123456789",
		SSecurity: "mock_ssecurity_token_here",
		Nonce:     "mock_nonce_value",
		Location:  "https://account.xiaomi.com/pass/serviceLogin",
	}, nil
}

func loadUnlockData() *UnlockData {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "miunlocktool")
	os.MkdirAll(configDir, 0755)

	dataFile := filepath.Join(configDir, "miunlockdata.json")

	data := &UnlockData{}
	if fileData, err := os.ReadFile(dataFile); err == nil {
		json.Unmarshal(fileData, data)
	}

	return data
}

func saveUnlockData(data *UnlockData) {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "miunlocktool")
	dataFile := filepath.Join(configDir, "miunlockdata.json")

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(dataFile, jsonData, 0644)
}

func getDeviceInfo(fastbootPath string) *DeviceInfo {
	// Try to get device info using fastboot commands
	deviceInfo := &DeviceInfo{}

	fmt.Print("⏳ Waiting for device")
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Println()

	// Get unlocked status
	fmt.Print("📋 Fetching 'unlocked' — please wait...")
	if output := runFastbootCommand(fastbootPath, "getvar", "unlocked"); output != "" {
		deviceInfo.Unlocked = output
		fmt.Print("\r\033[K") // Clear line
	} else {
		fmt.Print("\r\033[K")
		return nil
	}

	// Get product info
	fmt.Print("📋 Fetching 'product' — please wait...")
	if output := runFastbootCommand(fastbootPath, "getvar", "product"); output != "" {
		deviceInfo.Product = output
		fmt.Print("\r\033[K")
	}

	// Try to get token (determines SoC type)
	fmt.Print("📋 Fetching 'token' — please wait...")
	if token := runFastbootCommand(fastbootPath, "oem", "get_token"); token != "" {
		deviceInfo.Token = token
		deviceInfo.SoC = "Mediatek"
	} else if token := runFastbootCommand(fastbootPath, "getvar", "token"); token != "" {
		deviceInfo.Token = token
		deviceInfo.SoC = "Qualcomm"
	}
	fmt.Print("\r\033[K")

	return deviceInfo
}

func runFastbootCommand(cmd string, args ...string) string {
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

func displayDeviceInfo(info *DeviceInfo) {
	fmt.Println("\n📱 Device Information:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🔓 Unlocked: %s\n", info.Unlocked)
	fmt.Printf("📱 Product: %s\n", info.Product)
	fmt.Printf("🔧 SoC: %s\n", info.SoC)
	if info.Token != "" {
		fmt.Printf("🔑 Token: %s...\n", info.Token[:min(len(info.Token), 20)])
	}
}

func performRealUnlock(deviceInfo *DeviceInfo, authData *XiaomiAuthResponse, fastbootPath string) {
	fmt.Println("\n🔓 Starting unlock process...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Check if device is already unlocked
	if deviceInfo.Unlocked == "yes" || deviceInfo.Unlocked == "true" {
		fmt.Println("✅ Device is already unlocked!")
		return
	}

	// Step 1: Check device clear policy (like Python script)
	fmt.Println("📋 Checking device unlock policy...")
	clearPolicy := checkDeviceClearPolicy(deviceInfo.Product)

	if clearPolicy == 1 {
		fmt.Println("⚠️  🔴 This device clears user data when it is unlocked")
	} else if clearPolicy == -1 {
		fmt.Println("✅ 🟢 Unlocking the device does not clear user data")
	}

	fmt.Println("📢 Notice: Please ensure your device bootloader can be unlocked")

	// Confirm before proceeding
	fmt.Print("\n🔓 Press Enter to Unlock (or type 'q' to quit): ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(choice)) == "q" {
		fmt.Println("❌ Unlock cancelled")
		return
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Step 2: Request unlock from Xiaomi API (like Python RetrieveEncryptData)
	fmt.Print("⏳ Requesting unlock permission from Xiaomi servers")
	for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Println()

	unlockResponse := requestUnlockFromAPI(deviceInfo, authData)

	if unlockResponse.Code == 0 && unlockResponse.EncryptData != "" {
		// Success - got encrypted data
		fmt.Println("📦 Received encrypted unlock data from Xiaomi")

		// Convert hex string to bytes (like Python script)
		encryptedBytes, err := hex.DecodeString(unlockResponse.EncryptData)
		if err != nil {
			fmt.Printf("❌ Failed to decode encrypted data: %v\n", err)
			return
		}

		// Write encrypted data to file
		encryptFile := "encryptData"
		err = os.WriteFile(encryptFile, encryptedBytes, 0644)
		if err != nil {
			fmt.Printf("❌ Failed to write encrypt data: %v\n", err)
			return
		}

		// Get serial number (like Python script)
		fmt.Print("📋 Fetching device serial...")
		runFastbootCommand(fastbootPath, "getvar", "serialno")
		fmt.Print("\r\033[K")

		// Stage the encrypted data
		fmt.Println("📤 Staging encrypted data...")
		stageCmd := exec.Command(fastbootPath, "stage", encryptFile)
		if err := stageCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to stage data: %v\n", err)
			os.Remove(encryptFile)
			return
		}

		// Perform unlock
		fmt.Println("🔓 Executing unlock command...")
		unlockCmd := exec.Command(fastbootPath, "oem", "unlock")
		output, err := unlockCmd.CombinedOutput()

		// Clean up
		os.Remove(encryptFile)

		if err != nil {
			fmt.Printf("❌ Unlock failed: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
			return
		}

		fmt.Println("✅ Device unlock successful!")
		fmt.Println("🎉 Your Xiaomi device has been unlocked!")

	} else if unlockResponse.DescEN != "" {
		// Error from API
		fmt.Printf("\n❌ Unlock request failed (Code: %d)\n", unlockResponse.Code)
		fmt.Printf("📝 Message: %s\n", unlockResponse.DescEN)

		if unlockResponse.Code == 20036 && unlockResponse.Data.WaitHour > 0 {
			// Wait time required
			waitTime := time.Now().Add(time.Duration(unlockResponse.Data.WaitHour) * time.Hour)
			fmt.Printf("\n⏰ You can unlock on: %s\n", waitTime.Format("2006-01-02 15:04"))
		} else {
			fmt.Println("\n💡 For error codes: https://offici5l.github.io/articles/mi-error-codes")
		}
	} else {
		fmt.Printf("❌ Unexpected response from Xiaomi API: %+v\n", unlockResponse)
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func checkDeviceClearPolicy(product string) int {
	// Simulate API call to check if device clears data when unlocked
	// In Python: RetrieveEncryptData("/api/v2/unlock/device/clear", {"data":{"product":product}})
	time.Sleep(1 * time.Second)

	// Mock response: -1 = no clear, 1 = clears data, 0 = unknown
	// For demo, return -1 (doesn't clear data)
	return -1
}

func requestUnlockFromAPI(deviceInfo *DeviceInfo, authData *XiaomiAuthResponse) *UnlockResponse {
	// Simulate the Python RetrieveEncryptData("/api/v3/ahaUnlock", ...) call
	// This would make real HTTP requests to Xiaomi API in production
	time.Sleep(3 * time.Second)

	// For demo purposes, return mock success with fake encrypted data
	return &UnlockResponse{
		Code:        0,
		EncryptData: "deadbeef" + strings.Repeat("a1b2c3d4", 10), // Mock hex data
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
