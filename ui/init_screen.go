package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// InitScreen represents the initialization/splash screen
type InitScreen struct {
	app         fyne.App
	window      fyne.Window
	progressBar *widget.ProgressBar
}

// NewInitScreen creates a new initialization screen
func NewInitScreen(app fyne.App) *InitScreen {
	return &InitScreen{
		app: app,
	}
}

// Show displays the initialization screen and starts the progress
func (i *InitScreen) Show() {
	// Create window
	i.window = i.app.NewWindow(lang.L("title"))
	i.window.Resize(fyne.NewSize(600, 400))
	i.window.CenterOnScreen()
	i.window.SetFixedSize(true)

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
	// App logo/icon
	logoIcon := widget.NewIcon(theme.ComputerIcon())
	logoIcon.Resize(fyne.NewSize(64, 64))

	// App title
	titleLabel := widget.NewLabelWithStyle(
		lang.L("title"),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Loading text
	loadingLabel := widget.NewLabelWithStyle(
		"Initializing...",
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)

	// Progress bar
	i.progressBar = widget.NewProgressBar()
	i.progressBar.Resize(fyne.NewSize(400, 20))

	// Status text
	statusLabel := widget.NewLabelWithStyle(
		"Loading application components...",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	// Arrange components
	logoContainer := container.NewCenter(logoIcon)

	contentContainer := container.NewVBox(
		logoContainer,
		layout.NewSpacer(),
		titleLabel,
		widget.NewSeparator(),
		loadingLabel,
		layout.NewSpacer(),
		i.progressBar,
		layout.NewSpacer(),
		statusLabel,
		layout.NewSpacer(),
	)

	return container.NewCenter(contentContainer)
}

// startProgress animates the progress bar and handles completion
func (i *InitScreen) startProgress() {
	// Simple approach: wait 2 seconds then complete
	time.Sleep(2 * time.Second)

	// Use fyne.Do to ensure UI operations happen on main thread
	fyne.Do(func() {
		// Set progress to complete
		i.progressBar.SetValue(1.0)

		// Small delay before transitioning
		time.Sleep(200 * time.Millisecond)

		// Close init screen and show login screen
		i.transitionToLoginScreen()
	})
}

// transitionToLoginScreen handles the transition from init to login screen
func (i *InitScreen) transitionToLoginScreen() {
	// Close init window
	if i.window != nil {
		i.window.Close()
	}

	// Create and show login screen
	loginScreen := NewLoginScreen(i.app)
	loginScreen.Show()
}
