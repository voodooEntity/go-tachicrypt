package main

import (
	"github.com/rivo/tview"
)

var app *tview.Application

func main() {
	app = tview.NewApplication()

	encryptButton := tview.NewButton("Encrypt").SetSelectedFunc(func() {

	})
	encryptButton.SetBorder(true).SetRect(0, 0, 22, 3)

	decryptButton := tview.NewButton("Decrypt").SetSelectedFunc(func() {

	})
	decryptButton.SetBorder(true).SetRect(0, 0, 22, 3)

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Left (1/2 x width of Top)"), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

func menu() {

}
