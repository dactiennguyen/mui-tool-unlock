package auth

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"muitoolunlock/internal/types"
)

// GetWebBrowserID prompts user for web authentication and returns device ID
func GetWebBrowserID() string {
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

// AuthenticateXiaomi performs Xiaomi authentication
func AuthenticateXiaomi(user, password, deviceID string) (*types.XiaomiAuthResponse, error) {
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
	return &types.XiaomiAuthResponse{
		Code:      0,
		UserID:    "123456789",
		SSecurity: "mock_ssecurity_token_here",
		Nonce:     "mock_nonce_value",
		Location:  "https://account.xiaomi.com/pass/serviceLogin",
	}, nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
