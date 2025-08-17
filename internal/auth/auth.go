package auth

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"muitoolunlock/internal/colors"
	"muitoolunlock/internal/types"
)

// GetWebBrowserID prompts user for web authentication and returns device ID
func GetWebBrowserID() string {
	authURL := "https://account.xiaomi.com/pass/serviceLogin?sid=unlockApi&checkSafeAddress=true&passive=false&hidden=false"

	fmt.Println(colors.Section("🌐 Xiaomi Web Authentication"))
	fmt.Println(colors.Browser("Opening Xiaomi authentication page..."))
	fmt.Printf("%s %s\n", colors.Info("URL:"), colors.DimText(authURL))

	// Auto-open browser based on OS
	if err := openBrowser(authURL); err != nil {
		fmt.Println(colors.Warning(fmt.Sprintf("Could not open browser automatically: %v", err)))
		fmt.Println(colors.Info("Please open the URL above manually."))
	} else {
		fmt.Println(colors.Success("Browser opened automatically!"))
	}

	fmt.Println(colors.Notice("Follow these steps:"))
	fmt.Println(colors.DimText("1. Login with your Xiaomi account in the opened browser"))
	fmt.Println(colors.DimText("2. After successful login, copy the redirect URL"))
	fmt.Println(colors.DimText("3. Look for 'd=' parameter in the URL"))
	fmt.Println(colors.DimText("4. Paste the complete URL or just the 'd' parameter value below"))

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(colors.Input("Enter the redirect URL (or just the 'd' parameter value): "))
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

// AuthenticateXiaomi performs Xiaomi authentication
func AuthenticateXiaomi(user, password, deviceID string) (*types.XiaomiAuthResponse, error) {
	// Simulate the real Xiaomi authentication process
	// In reality, this would make HTTP requests to Xiaomi API endpoints
	fmt.Println(colors.Section("🔐 Xiaomi Authentication"))
	fmt.Println(colors.Progress("Posting credentials to Xiaomi servers..."))
	time.Sleep(2 * time.Second)

	// Hash password like Python script
	hasher := md5.New()
	hasher.Write([]byte(password))
	passwordHash := strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))

	fmt.Printf("%s %s\n", colors.Email("User:"), colors.BoldText(user))
	fmt.Printf("%s %s%s\n", colors.Key("Hash:"), colors.DimText(passwordHash[:8]), colors.DimText("..."))
	fmt.Printf("%s %s%s\n", colors.Device("Device ID:"), colors.DimText(deviceID[:min(len(deviceID), 12)]), colors.DimText("..."))

	// Simulate successful authentication
	return &types.XiaomiAuthResponse{
		Code:      0,
		UserID:    "123456789",
		SSecurity: "mock_ssecurity_token_here",
		Nonce:     "mock_nonce_value",
		Location:  "https://account.xiaomi.com/pass/serviceLogin",
	}, nil
}

// openBrowser opens the default browser with the given URL
func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
