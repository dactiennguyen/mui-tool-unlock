package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// LoginScreen represents the login screen
type LoginScreen struct {
	app           fyne.App
	window        fyne.Window
	emailEntry    *widget.Entry
	passEntry     *widget.Entry
	linkEntry     *widget.Entry
	loginButton   *widget.Button
	verifyButton  *widget.Button
	backButton    *widget.Button
	isLinkMode    bool
	mainContainer *fyne.Container
}

// NewLoginScreen creates a new login screen
func NewLoginScreen(app fyne.App) *LoginScreen {
	return &LoginScreen{
		app: app,
	}
}

// Show displays the login screen
func (l *LoginScreen) Show() {
	// Create window
	l.window = l.app.NewWindow(lang.L("title"))
	l.window.Resize(fyne.NewSize(500, 450))
	l.window.CenterOnScreen()
	l.window.SetFixedSize(true)

	// Initialize state
	l.isLinkMode = false

	// Create content
	content := l.createContent()
	l.window.SetContent(content)

	// Show window
	l.window.Show()
}

// createContent creates the login screen content
func (l *LoginScreen) createContent() *fyne.Container {
	// App logo/icon
	logoIcon := widget.NewIcon(theme.AccountIcon())
	logoIcon.Resize(fyne.NewSize(80, 80))
	logoContainer := container.NewCenter(logoIcon)

	// Title
	titleLabel := widget.NewLabelWithStyle(
		lang.L("login_mui_account"),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Subtitle
	subtitleLabel := widget.NewLabelWithStyle(
		lang.L("sign_in_subtitle"),
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)

	// Email field
	emailLabel := widget.NewLabelWithStyle(
		lang.L("email"),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	l.emailEntry = widget.NewEntry()
	l.emailEntry.SetPlaceHolder(lang.L("email_placeholder"))

	// Password field
	passwordLabel := widget.NewLabelWithStyle(
		lang.L("password"),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	l.passEntry = widget.NewPasswordEntry()
	l.passEntry.SetPlaceHolder(lang.L("password_placeholder"))

	// Link field (for link verification mode)
	linkLabel := widget.NewLabelWithStyle(
		lang.L("link"),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	l.linkEntry = widget.NewEntry()
	l.linkEntry.SetPlaceHolder(lang.L("link_placeholder"))

	// Login button
	l.loginButton = widget.NewButton(lang.L("login_mui"), l.handleLogin)
	l.loginButton.Importance = widget.HighImportance

	// Verify button
	l.verifyButton = widget.NewButton(lang.L("verify_link"), l.handleVerifyLink)
	l.verifyButton.Importance = widget.HighImportance

	// Back button
	l.backButton = widget.NewButton(lang.L("back"), l.handleBack)

	// Create the main container that will switch between modes
	l.mainContainer = l.createLoginForm(
		logoContainer, titleLabel, subtitleLabel,
		emailLabel, passwordLabel, linkLabel,
	)

	// Full width container with padding
	return container.NewBorder(
		nil, nil,
		layout.NewSpacer(), layout.NewSpacer(), // left and right spacers for padding
		container.NewPadded(l.mainContainer),
	)
}

// createLoginForm creates the login form layout
func (l *LoginScreen) createLoginForm(logoContainer *fyne.Container, titleLabel, subtitleLabel, emailLabel, passwordLabel, linkLabel *widget.Label) *fyne.Container {
	return container.NewVBox(
		logoContainer,
		layout.NewSpacer(),
		titleLabel,
		subtitleLabel,
		widget.NewSeparator(),
		layout.NewSpacer(),
		emailLabel,
		l.emailEntry,
		layout.NewSpacer(),
		passwordLabel,
		l.passEntry,
		layout.NewSpacer(),
		l.loginButton,
		layout.NewSpacer(),
	)
}

// createLinkForm creates the link verification form layout
func (l *LoginScreen) createLinkForm(logoContainer *fyne.Container, titleLabel, subtitleLabel, linkLabel *widget.Label) *fyne.Container {
	// Change subtitle for link mode
	linkSubtitle := widget.NewLabelWithStyle(
		lang.L("link_subtitle"),
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)

	// Main form content
	formContent := container.NewVBox(
		logoContainer,
		layout.NewSpacer(),
		titleLabel,
		linkSubtitle,
		widget.NewSeparator(),
		layout.NewSpacer(),
		linkLabel,
		l.linkEntry,
		layout.NewSpacer(),
		container.NewCenter(l.verifyButton), // Verify button centered
		layout.NewSpacer(),
	)

	// Create layout with back button at top-left corner
	return container.NewBorder(
		container.NewHBox(l.backButton, layout.NewSpacer()), // top with back button
		nil,         // bottom
		nil,         // left
		nil,         // right
		formContent, // center content
	)
}

// switchToLinkMode switches the UI to link verification mode
func (l *LoginScreen) switchToLinkMode() {
	l.isLinkMode = true

	// Get the components we need to recreate the form
	logoIcon := widget.NewIcon(theme.AccountIcon())
	logoIcon.Resize(fyne.NewSize(80, 80))
	logoContainer := container.NewCenter(logoIcon)

	titleLabel := widget.NewLabelWithStyle(
		lang.L("login_mui_account"),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	linkLabel := widget.NewLabelWithStyle(
		lang.L("link"),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	// Create new content for link mode
	newContent := l.createLinkForm(logoContainer, titleLabel, nil, linkLabel)

	// Update the main container
	l.mainContainer.Objects = newContent.Objects
	l.mainContainer.Refresh()
}

// switchToLoginMode switches the UI back to login mode
func (l *LoginScreen) switchToLoginMode() {
	l.isLinkMode = false

	// Only clear link entry, keep email and password
	l.linkEntry.SetText("")

	// Get the components we need to recreate the form
	logoIcon := widget.NewIcon(theme.AccountIcon())
	logoIcon.Resize(fyne.NewSize(80, 80))
	logoContainer := container.NewCenter(logoIcon)

	titleLabel := widget.NewLabelWithStyle(
		lang.L("login_mui_account"),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	subtitleLabel := widget.NewLabelWithStyle(
		lang.L("sign_in_subtitle"),
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)

	emailLabel := widget.NewLabelWithStyle(
		lang.L("email"),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	passwordLabel := widget.NewLabelWithStyle(
		lang.L("password"),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	linkLabel := widget.NewLabelWithStyle(
		lang.L("link"),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	// Create new content for login mode
	newContent := l.createLoginForm(logoContainer, titleLabel, subtitleLabel, emailLabel, passwordLabel, linkLabel)

	// Update the main container
	l.mainContainer.Objects = newContent.Objects
	l.mainContainer.Refresh()
}

// handleLogin processes login attempt
func (l *LoginScreen) handleLogin() {
	email := l.emailEntry.Text
	password := l.passEntry.Text

	// Basic validation
	if email == "" || password == "" {
		dialog.ShowInformation(
			"Error",
			"Please enter both email and password",
			l.window,
		)
		return
	}

	// Direct switch to link verification mode (no credential validation)
	l.switchToLinkMode()
}

// handleVerifyLink processes link verification
func (l *LoginScreen) handleVerifyLink() {
	link := l.linkEntry.Text

	// Basic validation
	if link == "" {
		dialog.ShowInformation(
			"Error",
			"Please enter a link to verify",
			l.window,
		)
		return
	}

	// Simple validation for demo
	if l.validateLink(link) {
		// Close login window and open unlock screen
		l.window.Close()
		
		// Create and show unlock screen
		unlockScreen := NewUnlockScreen(l.app)
		unlockScreen.Show()
	} else {
		dialog.ShowInformation(
			"Error",
			"Invalid link. Please check and try again.",
			l.window,
		)
	}
}

// handleBack handles back button press
func (l *LoginScreen) handleBack() {
	l.switchToLoginMode()
}

// validateLink validates the provided link
func (l *LoginScreen) validateLink(link string) bool {
	// Simple validation for demo - in a real app, you'd verify against your system
	// For demo purposes, accept links that contain "mui" or start with "http"
	if len(link) < 5 {
		return false
	}

	// Accept any reasonable URL format or containing "mui"
	return link[:4] == "http" || len(link) > 10
}
