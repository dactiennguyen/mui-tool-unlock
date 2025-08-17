package ui

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// InitScreen represents the initialization/splash screen
type InitScreen struct {
	app         fyne.App
	window      fyne.Window
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
}

// NewInitScreen creates a new initialization screen
func NewInitScreen(app fyne.App) *InitScreen {
	return &InitScreen{
		app: app,
	}
}

// Show displays the initialization screen and starts the progress
func (i *InitScreen) Show() {
	// Create compact window
	i.window = i.app.NewWindow("🔓 MUI Tool Unlock")
	i.window.Resize(fyne.NewSize(500, 400))
	i.window.CenterOnScreen()
	i.window.SetFixedSize(true)

	// Set window properties for better appearance
	i.window.SetPadded(true)

	// Create content
	content := i.createContent()
	i.window.SetContent(content)

	// Show window
	i.window.Show()

	// Start progress animation in a safe way
	go i.startProgress()
}

// createContent creates the initialization screen content
func (i *InitScreen) createContent() *fyne.Container {
	// App logo/icon - Compact size
	logoIcon := widget.NewIcon(theme.ComputerIcon())
	logoIcon.Resize(fyne.NewSize(64, 64))

	// App title - Centered
	titleLabel := widget.NewLabelWithStyle(
		"🔓 MUI Tool Unlock",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Subtitle - Centered
	subtitleLabel := widget.NewLabelWithStyle(
		"Xiaomi Device Bootloader Unlocker",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	// Version info - Centered
	versionLabel := widget.NewLabelWithStyle(
		"Version 1.5.9",
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)

	// Loading text - Centered
	loadingLabel := widget.NewLabelWithStyle(
		"🚀 Initializing...",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Progress bar - Compact size
	i.progressBar = widget.NewProgressBar()
	i.progressBar.Resize(fyne.NewSize(350, 20))

	// Status text - Centered
	i.statusLabel = widget.NewLabelWithStyle(
		"🔍 Checking platform tools...",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	// Main content with better spacing
	mainContent := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(logoIcon),
		container.NewCenter(titleLabel),
		container.NewCenter(subtitleLabel),
		container.NewCenter(versionLabel),
		widget.NewSeparator(),
		container.NewCenter(loadingLabel),
		layout.NewSpacer(),
		container.NewCenter(i.progressBar),
		container.NewCenter(i.statusLabel),
		layout.NewSpacer(),
	)

	return mainContent
}

// startProgress performs real platform setup with progress tracking
func (i *InitScreen) startProgress() {
	go func() {
		// Quick check first - if tools exist, skip all UI and go straight to login
		if i.checkPlatformToolsExist() {
			i.transitionToLoginScreen()
			return
		}

		// Tools don't exist, show download progress
		success := i.downloadAndSetupPlatformTools()

		fyne.Do(func() {
			if success {
				// Setup completed successfully
				i.statusLabel.SetText("✅ Setup completed!")
				i.progressBar.SetValue(1.0)

				// Direct transition without delays
				go func() {
					time.Sleep(200 * time.Millisecond) // Brief success display
					i.transitionToLoginScreen()
				}()
			} else {
				// Setup failed - show detailed error
				i.statusLabel.SetText("❌ Setup failed!")
				i.progressBar.SetValue(0.0)

				errorDialog := dialog.NewError(
					fmt.Errorf("❌ Platform Tools Setup Failed\n\nPossible causes:\n• No internet connection\n• Firewall blocking downloads\n• Insufficient disk space\n\nPlease check your connection and try again."),
					i.window,
				)
				errorDialog.SetDismissText("Retry")
				errorDialog.SetOnClosed(func() {
					go i.startProgress()
				})
				errorDialog.Show()
			}
		})
	}()
}

// checkPlatformToolsExist checks if platform tools already exist (no UI updates)
func (i *InitScreen) checkPlatformToolsExist() bool {
	baseDir, err := os.Getwd()
	if err != nil {
		return false
	}

	var fastbootName string
	if runtime.GOOS == "windows" {
		fastbootName = "fastboot.exe"
	} else {
		fastbootName = "fastboot"
	}

	fastbootPath := filepath.Join(baseDir, "platform-tools", fastbootName)
	_, err = os.Stat(fastbootPath)
	return err == nil
}

// downloadAndSetupPlatformTools performs download and setup with minimal UI updates
func (i *InitScreen) downloadAndSetupPlatformTools() bool {
	baseDir, err := os.Getwd()
	if err != nil {
		return false
	}

	// Show download status
	fyne.Do(func() {
		i.statusLabel.SetText("⬇️ Downloading platform tools...")
		i.progressBar.SetValue(0.3)
	})

	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "darwin"
	}

	url := fmt.Sprintf("https://dl.google.com/android/repository/platform-tools-latest-%s.zip", osName)
	zipPath := filepath.Join(baseDir, "platform-tools.zip")

	// Download
	if err := i.downloadFileWithProgress(url, zipPath); err != nil {
		return false
	}

	// Show extraction status
	fyne.Do(func() {
		i.statusLabel.SetText("📦 Extracting...")
		i.progressBar.SetValue(0.7)
	})

	// Extract with optimized progress
	if err := i.unzipFileOptimized(zipPath, baseDir); err != nil {
		os.Remove(zipPath)
		return false
	}

	// Cleanup and permissions
	os.Remove(zipPath)
	if runtime.GOOS != "windows" {
		fastbootPath := filepath.Join(baseDir, "platform-tools", "fastboot")
		os.Chmod(fastbootPath, 0755)
	}

	return true
}

// downloadFileWithProgress downloads a file with simplified progress updates
func (i *InitScreen) downloadFileWithProgress(url, filepath string) error {
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

	// Simple progress update during download (no delays)
	fyne.Do(func() {
		i.statusLabel.SetText("📊 Downloading...")
		i.progressBar.SetValue(0.5)
	})

	_, err = io.Copy(out, resp.Body)

	if err == nil {
		fyne.Do(func() {
			i.statusLabel.SetText("✅ Download completed!")
			i.progressBar.SetValue(0.7)
		})
	}

	return err
}

// unzipFileOptimized extracts zip file with minimal UI updates for speed
func (i *InitScreen) unzipFileOptimized(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Simple progress update - no per-file updates
	fyne.Do(func() {
		i.statusLabel.SetText("📦 Extracting files...")
		i.progressBar.SetValue(0.8)
	})

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
		} else {
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
	}

	// Single completion update
	fyne.Do(func() {
		i.statusLabel.SetText("✅ Extraction completed!")
		i.progressBar.SetValue(0.95)
	})

	return nil
}

// transitionToLoginScreen handles the transition from init to login screen
func (i *InitScreen) transitionToLoginScreen() {
	// Use fyne.Do for thread safety
	fyne.Do(func() {
		// Close init window
		if i.window != nil {
			i.window.Close()
		}

		// Create and show login screen
		loginScreen := NewLoginScreen(i.app)
		loginScreen.Show()
	})
}
