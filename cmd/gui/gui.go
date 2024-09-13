package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/voodooEntity/go-tachicrypt/src/core"
	"image/color"
	"strconv"
	"strings"
)

var tachicrypt *TachiCrypt

func main() {
	tachicrypt = NewTachiCrypt(app.New())
	tachicrypt.ShowMain()
}

type TachiCrypt struct {
	App            fyne.App
	Windows        map[string]fyne.Window
	TargetType     string
	TargetPath     string
	OutputPath     string
	PartsAmount    int
	MasterPassword string
}

func (tc *TachiCrypt) CreateWindow(ident string, title string, width float32, height float32, content *fyne.Container) {
	if _, ok := tc.Windows[ident]; !ok {
		tc.Windows[ident] = tc.App.NewWindow(title)
		tc.Windows[ident].SetContent(content)
		tc.Windows[ident].Resize(fyne.NewSize(width, height))
		if ident == "main" {
			tc.Windows[ident].ShowAndRun()
		} else {
			tc.Windows[ident].Show()
		}
	}
}

func (tc *TachiCrypt) GetWindow(ident string) fyne.Window {
	if _, ok := tc.Windows[ident]; ok {
		return tc.Windows[ident]
	}
	return nil
}

func (tc *TachiCrypt) UpdateWindow(ident string, title string, content *fyne.Container, width float32, height float32) {
	cw := tc.GetWindow(ident)
	//cw.SetTitle(title)
	cw.SetContent(content)
	cw.Resize(fyne.NewSize(width, height))
}

func (tc *TachiCrypt) UpdateWindowContent(ident string, content *fyne.Container) {
	cw := tc.GetWindow(ident)
	cw.SetContent(content)
	cw.Content().Refresh()
}

func NewTachiCrypt(app fyne.App) *TachiCrypt {
	return &TachiCrypt{
		App:        app,
		TargetPath: "",
		OutputPath: "",
		Windows:    make(map[string]fyne.Window),
	}
}

func (tc *TachiCrypt) ShowMain() {
	fmt.Println("yes we get here")

	text := canvas.NewText("You choose to encrypt. Please fill the following fields.", color.White)
	text.Alignment = fyne.TextAlignTrailing
	text.TextStyle = fyne.TextStyle{Italic: true}
	textContainer := container.New(layout.NewHBoxLayout(), text)

	encrypt := widget.NewButton("Encrypt", func() {
		tc.ShowEncDataChoose()
	})

	decrypt := widget.NewButton("Decrypt", func() {
		tc.ShowDecDataChoose()
	})

	buttonContainer := container.New(layout.NewHBoxLayout(), encrypt, decrypt)

	cont := container.New(layout.NewVBoxLayout(), textContainer, buttonContainer)
	tc.CreateWindow("main", "~ TachiCrypt ~", 300, 200, cont)
}

func (tc *TachiCrypt) ShowEncDataChoose() {
	text := canvas.NewText("You choose to encrypt. Please either choose a target directory or file to encrypt and proceed.", color.White)
	text.Alignment = fyne.TextAlignTrailing
	text.TextStyle = fyne.TextStyle{Italic: true}
	textContainer := container.New(layout.NewHBoxLayout(), text)

	pathOutput := GetTextElement("Chosen data path: ", color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	pathContainer := container.New(layout.NewVBoxLayout(), pathOutput)

	chooseFile := widget.NewButton("Choose file", func() {
		onChosen := func(f fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println(err)
				return
			}
			if f == nil {
				return
			}
			tc.TargetPath = f.URI().String()
			tc.TargetType = "file"
			fmt.Printf("chosen: %v", f.URI())
			pathContainer.RemoveAll()
			pathContainer.Add(GetTextElement("Chosen data path: "+f.URI().String(), color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true}))
		}
		dialog.ShowFileOpen(onChosen, tc.GetWindow("main"))
	})

	chooseDirectory := widget.NewButton("Choose directory", func() {
		fmt.Println("choose directory")
		onChosen := func(f fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println(err)
				return
			}
			if f == nil {
				return
			}
			tc.TargetPath = f.Path()
			tc.TargetType = "directory"
			fmt.Printf("chosen: %v", f.Path())
			pathContainer.RemoveAll()
			pathContainer.Add(GetTextElement("Chosen data path: "+f.Path(), color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true}))
		}
		dialog.ShowFolderOpen(onChosen, tc.GetWindow("main"))
	})

	proceed := widget.NewButton("Proceed", func() {
		fmt.Println("proceed")
		if tc.TargetPath != "" {
			fmt.Println(tc.TargetPath, tc.TargetType)
			tc.ShowEncOutputPath()
		}
	})

	buttonContainer := container.New(layout.NewHBoxLayout(), chooseDirectory, chooseFile)

	cont := container.New(layout.NewVBoxLayout(), textContainer, buttonContainer, pathContainer, proceed)
	tc.UpdateWindow("main", " ~ TachiCrypt : Encrypt : Choose data ~ ", cont, 300, 300)
}

func (tc *TachiCrypt) ShowEncOutputPath() {
	text := canvas.NewText("Please choose the target directory where the encrypted files should be written to.", color.White)
	text.Alignment = fyne.TextAlignTrailing
	text.TextStyle = fyne.TextStyle{Italic: true}
	textContainer := container.New(layout.NewHBoxLayout(), text)

	pathOutput := GetTextElement("Chosen output path: ", color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	pathContainer := container.New(layout.NewVBoxLayout(), pathOutput)

	chooseDirectory := widget.NewButton("Choose directory", func() {
		fmt.Println("choose directory")
		onChosen := func(f fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println(err)
				return
			}
			if f == nil {
				return
			}
			tc.OutputPath = f.Path()
			fmt.Printf("chosen: %v", f.Path())
			pathContainer.RemoveAll()
			pathContainer.Add(GetTextElement("Chosen output path: "+f.Path(), color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true}))
		}
		dialog.ShowFolderOpen(onChosen, tc.GetWindow("main"))
	})

	proceed := widget.NewButton("Proceed", func() {
		fmt.Println("proceed")
		if tc.OutputPath != "" {
			fmt.Println(tc.OutputPath)
			tc.ShowEncGetMasterPwdAndPartsAmount()
		}
	})

	buttonContainer := container.New(layout.NewHBoxLayout(), chooseDirectory)

	cont := container.New(layout.NewVBoxLayout(), textContainer, buttonContainer, pathContainer, proceed)
	tc.UpdateWindow("main", " ~ TachiCrypt : Encrypt : Choose output path ~ ", cont, 300, 300)
}

func (tc *TachiCrypt) ShowEncGetMasterPwdAndPartsAmount() {
	text := widget.NewLabelWithStyle("Please enter your masterlock password and the amount of encrypted parts to be created. The masterlock password will be used to encrypt the 'masterlock' file which includes all the information to decrypt the encrypted parts. The amount of encrypted parts will decide on how many parts will be created as result - every part is necessary to decrypt the data. The amount needs to be an integer > 0.", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	text.Wrapping = fyne.TextWrapWord
	textContainer := container.New(layout.NewAdaptiveGridLayout(1), text)

	// password input
	passwordText := GetTextElement("Enter your masterlock password: ", color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	passwordInput := widget.NewEntry()
	passwordInput.SetPlaceHolder("")
	passwordInput.Resize(fyne.NewSize(100, 20))
	passwordContainer := container.New(layout.NewAdaptiveGridLayout(2), passwordText, passwordInput)

	// password input
	partsAmountText := GetTextElement("Enter amount of parts: ", color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	partsAmountInput := widget.NewEntry()
	partsAmountInput.Resize(fyne.NewSize(100, 20))
	partsAmountInput.SetPlaceHolder("")
	partsAmountContainer := container.New(layout.NewAdaptiveGridLayout(2), partsAmountText, partsAmountInput)

	proceed := widget.NewButton("Proceed", func() {
		if partsAmountInput.Text != "" {
			pa, err := strconv.Atoi(partsAmountInput.Text)
			if err != nil {
				fmt.Println(err)
				return
			}
			tc.PartsAmount = pa
		}
		if passwordInput.Text == "" {
			fmt.Println("no master password given")
			return
		}
		tc.MasterPassword = passwordInput.Text
		fmt.Println("Password for encryption provided |" + tc.MasterPassword + "|")
		fmt.Println("proceed")
		fmt.Println(tc.TargetPath, tc.TargetType)
		tc.ShowEncRunEncryption()
	})

	cont := container.New(layout.NewVBoxLayout(), textContainer, passwordContainer, partsAmountContainer, proceed)
	tc.UpdateWindow("main", " ~ TachiCrypt : Encrypt : Password & PartAmount ~ ", cont, 300, 300)
}

func (tc *TachiCrypt) ShowEncRunEncryption() {
	text := widget.NewLabelWithStyle("Running encyption", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	text.Wrapping = fyne.TextWrapWord
	textContainer := container.New(layout.NewAdaptiveGridLayout(1), text)

	infinite := widget.NewProgressBarInfinite()

	cont := container.New(layout.NewVBoxLayout(), textContainer, infinite)
	q := make(chan bool, 1)
	c := core.New()
	go func() {
		c.Hide(strings.TrimPrefix(tc.TargetPath, "file://"), tc.PartsAmount, tc.OutputPath, tc.MasterPassword)
		q <- true
	}()
	tc.UpdateWindow("main", " ~ TachiCrypt : Encrypt : Password & PartAmount ~ ", cont, 300, 300)
	for {
		select {
		case r := <-q:
			if r {
				tc.ShowEncFinal()
				return
			}
		}
	}
}

func (tc *TachiCrypt) ShowEncFinal() {
	text := widget.NewLabelWithStyle("Encryption done. Please check the provided output directory for the resulting files.", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	text.Wrapping = fyne.TextWrapWord
	textContainer := container.New(layout.NewAdaptiveGridLayout(1), text)

	exit := widget.NewButton("Exit", func() {
		tc.App.Quit()
	})

	cont := container.New(layout.NewAdaptiveGridLayout(1), textContainer, exit)
	tc.UpdateWindow("main", " ~ TachiCrypt : Encrypt : Final ~ ", cont, 300, 300)
}

func (tc *TachiCrypt) ShowDecDataChoose() {
	text := canvas.NewText("You choose to decrypt. Please choose the directory containing the encrypted data parts and the masterlock file.", color.White)
	text.Alignment = fyne.TextAlignTrailing
	text.TextStyle = fyne.TextStyle{Italic: true}
	textContainer := container.New(layout.NewAdaptiveGridLayout(1), text)

	pathOutput := GetTextElement("Chosen data path: ", color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	pathContainer := container.New(layout.NewVBoxLayout(), pathOutput)

	chooseDirectory := widget.NewButton("Choose directory", func() {
		fmt.Println("choose directory")
		onChosen := func(f fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println(err)
				return
			}
			if f == nil {
				return
			}
			tc.TargetPath = f.Path()
			tc.TargetType = "directory"
			fmt.Printf("chosen: %v", f.Path())
			pathContainer.RemoveAll()
			pathContainer.Add(GetTextElement("Chosen data path: "+f.Path(), color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true}))
		}
		dialog.ShowFolderOpen(onChosen, tc.GetWindow("main"))
	})

	proceed := widget.NewButton("Proceed", func() {
		fmt.Println("proceed")
		if tc.TargetPath != "" {
			fmt.Println(tc.TargetPath, tc.TargetType)
			tc.ShowDecOutputPath()
		}
	})

	buttonContainer := container.New(layout.NewHBoxLayout(), chooseDirectory)

	cont := container.New(layout.NewVBoxLayout(), textContainer, buttonContainer, pathContainer, proceed)
	tc.UpdateWindow("main", " ~ TachiCrypt : Decrypt : Choose data ~ ", cont, 300, 300)
}

func (tc *TachiCrypt) ShowDecOutputPath() {
	text := canvas.NewText("Please choose the target directory where the decrypted files should be written to.", color.White)
	text.Alignment = fyne.TextAlignTrailing
	text.TextStyle = fyne.TextStyle{Italic: true}
	textContainer := container.New(layout.NewHBoxLayout(), text)

	pathOutput := GetTextElement("Chosen output path: ", color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	pathContainer := container.New(layout.NewVBoxLayout(), pathOutput)

	chooseDirectory := widget.NewButton("Choose directory", func() {
		fmt.Println("choose directory")
		onChosen := func(f fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println(err)
				return
			}
			if f == nil {
				return
			}
			tc.OutputPath = f.Path()
			fmt.Printf("chosen: %v", f.Path())
			pathContainer.RemoveAll()
			pathContainer.Add(GetTextElement("Chosen output path: "+f.Path(), color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true}))
		}
		dialog.ShowFolderOpen(onChosen, tc.GetWindow("main"))
	})

	proceed := widget.NewButton("Proceed", func() {
		fmt.Println("proceed")
		if tc.OutputPath != "" {
			fmt.Println(tc.OutputPath)
			tc.ShowDecGetMasterPwd()
		}
	})

	buttonContainer := container.New(layout.NewHBoxLayout(), chooseDirectory)

	cont := container.New(layout.NewVBoxLayout(), textContainer, buttonContainer, pathContainer, proceed)
	tc.UpdateWindow("main", " ~ TachiCrypt : Decrypt : Choose output path ~ ", cont, 300, 300)
}

func (tc *TachiCrypt) ShowDecGetMasterPwd() {
	text := widget.NewLabelWithStyle("Please enter your masterlock password.", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	text.Wrapping = fyne.TextWrapWord
	textContainer := container.New(layout.NewAdaptiveGridLayout(1), text)

	// password input
	passwordText := GetTextElement("Enter your masterlock password: ", color.White, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	passwordInput := widget.NewEntry()
	passwordInput.SetPlaceHolder("")
	passwordInput.Resize(fyne.NewSize(100, 20))
	passwordContainer := container.New(layout.NewAdaptiveGridLayout(2), passwordText, passwordInput)

	proceed := widget.NewButton("Proceed", func() {
		if passwordInput.Text == "" {
			fmt.Println("no master password given")
			return
		}
		tc.MasterPassword = passwordInput.Text
		fmt.Println("proceed")
		fmt.Println("Password for decryption provided |" + tc.MasterPassword + "|")
		fmt.Println(tc.TargetPath, tc.TargetType)
		tc.ShowDecRunDecryption()
	})

	cont := container.New(layout.NewVBoxLayout(), textContainer, passwordContainer, proceed)
	tc.UpdateWindow("main", " ~ TachiCrypt : Decrypt : Password ~ ", cont, 300, 300)
}

func (tc *TachiCrypt) ShowDecRunDecryption() {
	text := widget.NewLabelWithStyle("Running decryption", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	text.Wrapping = fyne.TextWrapWord
	textContainer := container.New(layout.NewAdaptiveGridLayout(1), text)

	infinite := widget.NewProgressBarInfinite()

	cont := container.New(layout.NewVBoxLayout(), textContainer, infinite)
	q := make(chan bool, 1)
	c := core.New()
	go func() {
		c.Unhide(strings.TrimPrefix(tc.TargetPath, "file://"), tc.OutputPath, tc.MasterPassword)
		q <- true
	}()
	tc.UpdateWindow("main", " ~ TachiCrypt : Decrypt : Run ~ ", cont, 300, 300)
	for {
		select {
		case r := <-q:
			if r {
				tc.ShowEncFinal()
				return
			}
		}
	}
}

func (tc *TachiCrypt) ShowDecFinal() {
	text := widget.NewLabelWithStyle("Decryption done. Please check the provided output directory for the resulting files.", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	text.Wrapping = fyne.TextWrapWord
	textContainer := container.New(layout.NewAdaptiveGridLayout(1), text)

	exit := widget.NewButton("Exit", func() {
		tc.App.Quit()
	})

	cont := container.New(layout.NewVBoxLayout(), textContainer, exit)
	tc.UpdateWindow("main", " ~ TachiCrypt : Encrypt : Final ~ ", cont, 300, 300)
}

func GetTextElement(text string, col color.Gray16, align fyne.TextAlign, style fyne.TextStyle) *canvas.Text {
	t := canvas.NewText(text, col)
	t.Alignment = align
	t.TextStyle = style
	return t
}
