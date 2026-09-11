//go:build windows

package main

import (
	"os"
	"sav/walk"
	. "sav/walk/declarative"
	"unsafe"

	"github.com/lxn/win"
)

func main() {
	m, err := newModel()
	if err != nil {
		m.error(err)
		os.Exit(1)
	}

	window := MainWindow{
		Title:    "Game Save Manager",
		AssignTo: &m.mw,
		Layout:   VBox{Spacing: 3},
		Children: []Widget{
			Label{Text: "Games"},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					VSplitter{
						StretchFactor: 3,
						Children: []Widget{
							ListBox{
								Name:                  "lb",
								AssignTo:              &m.lb,
								Model:                 m.list,
								OnCurrentIndexChanged: m.listSelectedIndexChanged,
							},
							Composite{
								Layout: HBox{MarginsZero: true, Spacing: 3},
								Children: []Widget{
									PushButton{
										Text:      "Add",
										OnClicked: m.buttonAddGameClick,
									},
									PushButton{
										Text:      "Remove",
										Enabled:   Bind("lb.CurrentIndex != -1"),
										OnClicked: m.buttonRemoveGameClick,
									},
								},
							},
						},
					},
					VSplitter{
						StretchFactor: 7,
						Enabled:       Bind("lb.CurrentIndex != -1"),
						Children: []Widget{
							Composite{
								Layout: Grid{Columns: 4, Spacing: 3, Margins: Margins{Left: 3, Right: 3}},
								Children: []Widget{
									Label{Text: "Name:"},
									LineEdit{AssignTo: &m.name},
									PushButton{Text: "Save", ColumnSpan: 2, OnClicked: m.buttonSaveNameClick},
									Label{Text: "Save Root:"},
									LineEdit{Name: "save", AssignTo: &m.save, ReadOnly: true, OnTextChanged: m.textSaveRootChanged},
									PushButton{Text: "Choose", OnClicked: m.buttonChooseSaveClick},
									PushButton{Text: "Open", Enabled: Bind(`save.Text != ""`), OnClicked: m.buttonOpenSaveRootClick},
									Label{Text: "Backup Path:"},
									LineEdit{Name: "backup", AssignTo: &m.backup, ReadOnly: true},
									PushButton{Text: "Choose", OnClicked: m.buttonChooseBackupClick},
									PushButton{Text: "Open", Enabled: Bind(`backup.Text != ""`), OnClicked: m.buttonOpenBackupClick},
								},
							},
							TreeView{AssignTo: &m.tv, OnCheckedChanged: m.treeCheckedChanged},
							PushButton{
								Text:      "Back up",
								Enabled:   Bind(`save.Text != "" && backup.Text != ""`),
								OnClicked: m.buttonBackUpClick,
							},
						},
					},
				},
			},
		},
	}
	if err := window.Create(); err != nil {
		m.error(err)
		os.Exit(1)
	}

	var mi win.MONITORINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))
	hMon := win.MonitorFromWindow(m.mw.Handle(), win.MONITOR_DEFAULTTOPRIMARY)
	if win.GetMonitorInfo(hMon, &mi) {
		work := mi.RcWork
		const WIDTH, HEIGHT = 900, 550
		m.try(m.mw.SetBoundsPixels(walk.Rectangle{
			X:      int(work.Left) + (int(work.Right-work.Left)-WIDTH)/2,
			Y:      int(work.Top) + (int(work.Bottom-work.Top)-HEIGHT)/2,
			Width:  WIDTH,
			Height: HEIGHT,
		}))
	}

	m.mw.Run()
}
