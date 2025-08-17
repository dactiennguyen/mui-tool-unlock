package platform

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"muitoolunlock/internal/colors"
)

// Setup sets up platform-tools and returns the fastboot path
func Setup() string {
	fmt.Println(colors.Section("🔧 Platform Tools Setup"))
	fmt.Println(colors.Package("Setting up platform-tools..."))

	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Println(colors.Error("Failed to get current directory"))
		return ""
	}
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
		fmt.Println(colors.Success("Platform-tools already available"))
		return fastbootPath
	}

	// Download platform-tools
	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "darwin"
	}

	url := fmt.Sprintf("https://dl.google.com/android/repository/platform-tools-latest-%s.zip", osName)
	zipPath := filepath.Join(baseDir, "platform-tools.zip")

	fmt.Println(colors.Download("Downloading platform-tools..."))
	fmt.Printf("%s %s\n", colors.Info("URL:"), colors.DimText(url))
	if err := downloadFile(url, zipPath); err != nil {
		fmt.Println(colors.Error(fmt.Sprintf("Failed to download platform-tools: %v", err)))
		return ""
	}

	fmt.Println(colors.Package("Extracting platform-tools..."))
	if err := unzipFile(zipPath, baseDir); err != nil {
		fmt.Println(colors.Error(fmt.Sprintf("Failed to extract platform-tools: %v", err)))
		os.Remove(zipPath)
		return ""
	}

	// Clean up zip file
	os.Remove(zipPath)
	fmt.Println(colors.Info("Cleaned up temporary files"))

	// Make fastboot executable on Unix systems
	if runtime.GOOS != "windows" {
		os.Chmod(fastbootPath, 0755)
		fmt.Println(colors.Info("Set executable permissions"))
	}

	fmt.Println(colors.Success("Platform-tools setup completed"))
	return fastbootPath
}

// downloadFile downloads a file from URL to filepath
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

// unzipFile extracts zip file to destination
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
