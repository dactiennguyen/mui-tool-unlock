package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// UnlockScreen represents the unlock screen
type UnlockScreen struct {
	app           fyne.App
	window        fyne.Window
	waitingLabel  *widget.Label
	unlockButton  *widget.Button
	mainContainer *fyne.Container
	isWaiting     bool
}

// NewUnlockScreen creates a new unlock screen
func NewUnlockScreen(app fyne.App) *UnlockScreen {
	return &UnlockScreen{
		app:       app,
		isWaiting: true,
	}
}

// Show displays the unlock screen
func (u *UnlockScreen) Show() {
	// Create window with fixed dimensions
	u.window = u.app.NewWindow(lang.L("title"))
	u.window.Resize(fyne.NewSize(600, 400))
	u.window.CenterOnScreen()
	u.window.SetFixedSize(true)

	// Create content
	content := u.createContent()
	u.window.SetContent(content)

	// Show window
	u.window.Show()

	// Start waiting sequence
	go u.startWaitingSequence()
}

// createContent creates the unlock screen content
func (u *UnlockScreen) createContent() *fyne.Container {
	// App logo/icon
	logoIcon := widget.NewIcon(theme.ComputerIcon())
	logoIcon.Resize(fyne.NewSize(80, 80))
	logoContainer := container.NewCenter(logoIcon)

	// Title
	titleLabel := widget.NewLabelWithStyle(
		lang.L("unlock_title"),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Waiting text (initially visible)
	u.waitingLabel = widget.NewLabelWithStyle(
		lang.L("waiting_to_connect"),
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	// Unlock button (initially hidden)
	u.unlockButton = widget.NewButton(lang.L("unlock"), u.handleUnlock)
	u.unlockButton.Importance = widget.HighImportance
	u.unlockButton.Resize(fyne.NewSize(200, 60))
	u.unlockButton.Hide() // Start hidden

	// Main container that will show either waiting text or unlock button
	u.mainContainer = container.NewVBox(
		logoContainer,
		layout.NewSpacer(),
		titleLabel,
		layout.NewSpacer(),
		u.waitingLabel,
		layout.NewSpacer(),
		container.NewCenter(u.unlockButton),
		layout.NewSpacer(),
	)

	// Full width container with padding
	return container.NewBorder(
		nil, nil,
		layout.NewSpacer(), layout.NewSpacer(), // left and right spacers for padding
		container.NewPadded(u.mainContainer),
	)
}

// startWaitingSequence handles the waiting -> unlock button sequence
func (u *UnlockScreen) startWaitingSequence() {
	// Wait 1 second
	time.Sleep(1 * time.Second)

	// Use fyne.Do to ensure UI operations happen on main thread
	fyne.Do(func() {
		// Hide waiting text and show unlock button
		u.waitingLabel.Hide()
		u.unlockButton.Show()
		u.isWaiting = false
		
		// Refresh the container
		u.mainContainer.Refresh()
	})
}

// handleUnlock processes unlock button press
func (u *UnlockScreen) handleUnlock() {
	// Show success dialog
	dialog.ShowInformation(
		"Success",
		"Device unlocked successfully!",
		u.window,
	)
	
	// Here you could add actual unlock logic
	// For now, just show success message
}
