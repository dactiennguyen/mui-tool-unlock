<div align="center">

  <a href="https://github.com/your-username/mui-tool-unlock/releases/latest"><img src="https://img.shields.io/badge/MUI%20Tool%20Unlock-%23007ACC?style=flat&logo=go&logoColor=white" alt="MUI Tool Unlock" style="width: 200px; vertical-align: middle;" /> </a><br>

  <img src="https://img.shields.io/github/v/release/your-username/mui-tool-unlock?style=flat&label=Version&labelColor=black&color=brightgreen" alt="Version" /><br><p style="font-weight: bold;">Developed as a GUI tool for device unlocking
  <br>
  with secure authentication and link verification.
  <br>
  Cross-platform support with modern UI.</p>
 
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="MIT License" /></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go Version" /></a>
  <a href="https://fyne.io/"><img src="https://img.shields.io/badge/Fyne-v2-FF6900?style=flat&logo=go&logoColor=white" alt="Fyne UI" /></a>
  
</div>

## ✨ Features

- 🚀 **Fast Initialization** - 2-second splash screen with progress animation
- 🔐 **Secure Authentication** - Multi-step login with link verification  
- 🌍 **Multilingual Support** - English and Vietnamese translations
- 🎨 **Modern UI** - Clean, responsive design with Fyne framework
- 📱 **Fixed Dimensions** - Consistent window sizing and centering
- ⚡ **Real-time Navigation** - Smooth transitions between screens
- 🔓 **Device Unlock** - Professional unlock interface with timing

## 🖥️ Screenshots

### Login Flow
```
Init Screen (2s) → Login → Link Verification → Unlock Screen
```

### Features by Screen
- **Init Screen**: Progress bar with loading animation
- **Login Screen**: Email/password authentication with full-width inputs
- **Link Verification**: Secure link validation with back navigation  
- **Unlock Screen**: "Waiting to connect phone..." → Unlock button

## 📦 Installation

### Requirements
- **Go 1.19+** installed on your system
- **Git** for cloning the repository

### For Windows, Linux, MacOS:

1. **Install Go** (if not already installed):
   ```bash
   # Download from https://golang.org/dl/
   # Or use package managers:
   
   # MacOS with Homebrew
   brew install go
   
   # Ubuntu/Debian
   sudo apt install golang-go
   
   # Windows with Chocolatey  
   choco install golang
   ```

2. **Clone and run the project**:
   ```bash
   git clone https://github.com/your-username/mui-tool-unlock.git
   cd mui-tool-unlock
   go mod tidy
   go run main.go
   ```

3. **Build executable** (optional):
   ```bash
   go build -o mui-tool-unlock
   ./mui-tool-unlock
   ```

### Cross-Platform Build:

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o mui-tool-unlock.exe

# MacOS  
GOOS=darwin GOARCH=amd64 go build -o mui-tool-unlock-mac

# Linux
GOOS=linux GOARCH=amd64 go build -o mui-tool-unlock-linux
```

## 🚀 Usage

1. **Launch the application**:
   ```bash
   go run main.go
   ```

2. **Authentication Flow**:
   - Wait for initialization (2 seconds)
   - Enter any email and password 
   - Click "Login Mui" to proceed
   - Enter verification link
   - Click "Verify Link" to continue

3. **Unlock Process**:
   - Automatic transition to unlock screen
   - Wait for "Waiting to connect phone..." (1 second)
   - Click "Unlock" button when available
   - Success confirmation dialog

## 🛠️ Development

### Project Structure
```
mui-tool-unlock/
├── main.go                 # Application entry point
├── ui/                     # UI components
│   ├── init_screen.go     # Initialization screen
│   ├── login_screen.go    # Login and link verification
│   └── unlock_screen.go   # Device unlock interface
├── translations/          # Internationalization
│   ├── en.json           # English translations
│   └── vn.json           # Vietnamese translations
├── go.mod                # Go module dependencies
└── README.md             # Project documentation
```

### Key Dependencies
- **fyne.io/fyne/v2** - Modern cross-platform GUI framework
- **embed** - Built-in Go package for embedding static files

### Adding New Languages
1. Create new translation file in `translations/` directory
2. Add translation keys following existing structure
3. Translations are automatically loaded via `embed.FS`

## 🎯 Technical Features

- **Fyne UI Framework** - Cross-platform native applications
- **Embedded Translations** - No external files needed for i18n
- **Thread-Safe UI Updates** - Proper goroutine management with `fyne.Do()`
- **State Management** - Clean separation between UI modes
- **Input Validation** - Basic validation for forms and links
- **Memory Efficient** - Minimal resource usage

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔧 Build Requirements

- Go 1.19 or higher
- CGO enabled for Fyne (default)
- Platform-specific dependencies for Fyne

## 📞 Support

If you encounter any issues or have questions:

1. Check existing [Issues](https://github.com/your-username/mui-tool-unlock/issues)
2. Create a new issue with detailed description
3. Include system information and error logs

---

<div align="center">
  <p>Made with ❤️ using Go and Fyne</p>
  <p>Cross-platform • Secure • Modern</p>
</div>
