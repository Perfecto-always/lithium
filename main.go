package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

type LithiumWindow struct {
	*walk.MainWindow
}

var file []string

func main() {
	mw := new(LithiumWindow)

	var openAction, showAboutBoxAction *walk.Action
	var edit *walk.TextEdit

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
		Layout:  VBox{},
		// All things in the view that interactable upon opening
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					ListBox{
						MaxSize: Size{Width: 200},
						MinSize: Size{Width: 100}},
					TextEdit{AssignTo: &edit,
						MinSize: Size{Width: 800}},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	mw.Run()
}

func (mw *LithiumWindow) openAction_Triggered() {
	dlg := new(walk.FileDialog)
	dlg.Title = "Select File"
	dlg.Filter = "All Files (*.*)|*.*"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file dialog:", err)
	} else if !ok {
		fmt.Println("No file selected.")
	} else {
		fmt.Println("Selected file:", dlg.FilePath)
		file[0] = dlg.FilePath
		// Here you can do something with the selected file path
	}
}

func (mw *LithiumWindow) showAboutBoxAction_Triggered() {
	walk.MsgBox(mw, "About Lithium", "Lithium is a supposed to be light weight text-editor.", walk.MsgBoxIconInformation)
}
