package main

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Main menu
	mainMenu := tview.NewList().
		AddItem("Encrypt", "E", 'E', nil).
		AddItem("Decrypt", "D", 'D', nil)

	// Set initial focus
	app.SetFocus(mainMenu)

	// Input form for encryption
	encryptForm := tview.NewForm().
		AddInputField("File/Directory Path:", "", 50, nil, nil).
		AddInputField("Number of Parts:", "", 10, nil, nil).
		AddPasswordField("Password:", "", 50, '*', nil).
		AddButton("Encrypt", func() {
			// Call core.Encrypt() with input values
			// Show loading animation
			go func() {
				for progress := 0; progress <= 100; progress++ {
					progressText := fmt.Sprintf("Encryption progress: %d%%", progress)
					app.QueueUpdateDraw(func() {
						progress.SetText(progressText)
					})
					time.Sleep(100 * time.Millisecond)
				}
			}()
		})

	// Input form for decryption
	decryptForm := tview.NewForm().
		AddInputField("Directory Path:", "", 50, nil, nil).
		AddPasswordField("Password:", "", 50, '*', nil).
		AddButton("Decrypt", func() {
			// Call core.Decrypt() with input values
			// Show loading animation
			go func() {
				for progress := 0; progress <= 100; progress++ {
					progressText := fmt.Sprintf("Decryption progress: %d%%", progress)
					app.QueueUpdateDraw(func() {
						progress.SetText(progressText)
					})
					time.Sleep(100 * time.Millisecond)
				}
			}()
		})

	// Progress indicator
	progress := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)

	// Main layout
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainMenu, 0, 1, true).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(progress, 0, 1, false)), 0, 1, true

	// Switch between input forms based on user selection
	mainMenu.SetSelectedFunc(func(index int, string) {
		switch index {
		case 0:
			layout.AddItem(encryptForm, 0, 1, true)
		case 1:
			layout.AddItem(decryptForm, 0, 1, true)
		}
		app.Draw()
	})

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
