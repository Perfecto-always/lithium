package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

// Added the components here
type LithiumWindow struct {
	*walk.MainWindow
	edit *walk.TextEdit
	// lb   *walk.ListBox
}

var (
	// file        []string
	currentFile string
	lineno      uint = 1
)

func main() {
	mw := new(LithiumWindow)

	var openAction, showAboutBoxAction *walk.Action

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Icon:     "./lithium.ico",
		Title:    "Lithium",
		// The Menubar for all thing menubar
		MenuItems: []MenuItem{
			Menu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "&Open File",
						Shortcut:    Shortcut{Modifiers: walk.ModControl, Key: walk.KeyO},
						OnTriggered: mw.openAction_Triggered,
					},
					Action{
						Text:        "&Save File",
						Shortcut:    Shortcut{Modifiers: walk.ModControl, Key: walk.KeyS},
						OnTriggered: mw.saveFile_Triggered,
					},
					Separator{},
					Action{
						Text:        "E&xit",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						AssignTo:    &showAboutBoxAction,
						Text:        "About",
						OnTriggered: mw.showAboutBoxAction_Triggered,
					},
				},
			},
		},
		MinSize: Size{Width: 300, Height: 200},
		Layout: VBox{
			MarginsZero: true,
			SpacingZero: true,
		},
		// All things in the view that interactable upon opening
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					// ListBox{
					// 	// AssignTo: &mw.lb,
					// 	MaxSize: Size{Width: 200},
					// 	MinSize: Size{Width: 100}},
					TextEdit{AssignTo: &mw.edit,
						MinSize: Size{Width: 800},
						Font:    Font{PointSize: 16},
					},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	mw.Run()
}

func (mw *LithiumWindow) openAction_Triggered() {
	lineno = 1

	dlg := new(walk.FileDialog)
	dlg.Title = "Select File"
	dlg.Filter = "All Files (*.*)|*.*"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file dialog:", err)
	} else if !ok {
		fmt.Println("No file selected.")
	} else {
		fmt.Println("Selected file:", dlg.FilePath)
		buf, err := os.ReadFile(dlg.FilePath)

		if err != nil {
			fmt.Println("Error reading file: ", dlg.FilePath)
		}

		currentFile = dlg.FilePath
		for i := range buf {
			if buf[i] == 13 {
				if buf[i+1] == 10 {
					lineno++
				}
			}
		}
		mw.edit.SetText(string(buf))
	}
}

func (mw *LithiumWindow) showAboutBoxAction_Triggered() {
	walk.MsgBox(mw, "About Lithium", "Lithium is a supposed to be light weight text-editor.", walk.MsgBoxIconInformation)
}

// Function to save the content of TextEdit to a file
func (mw *LithiumWindow) saveFile_Triggered() {
	if currentFile == "" {
		return
	}
	// Get the text content from the TextEdit
	textContent := mw.edit.Text()

	// Write the text content to the file
	err := os.WriteFile(currentFile, []byte(textContent), 0644)
	if err != nil {
		fmt.Println("Error saving file")
	}
}
