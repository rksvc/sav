package main

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	a, err := newApp()
	if err != nil {
		mygo.Dialog.Error("", err.Error())
	}
	mygo.App.WhenReady(func() {
		a.win = mygo.NewWindow(mygo.WindowOptions{
			Title:   "Game Save Manager",
			Content: ui.View(a.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		mygo.Dialog.Error("", err.Error())
	}
}
