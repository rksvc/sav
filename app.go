package main

import (
	"archive/zip"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

//go:embed icons/circle-question-mark.svg
var questionSVG []byte

//go:embed icons/circle-plus.svg
var plusSVG []byte

//go:embed icons/trash.svg
var trashSVG []byte

//go:embed icons/folder-open.svg
var folderSVG []byte

//go:embed icons/square-arrow-out-up-right.svg
var externalSVG []byte

//go:embed icons/chevron-right.svg
var chevronSVG []byte

var (
	question = ui.MustParseSVG(questionSVG)
	plus     = ui.MustParseSVG(plusSVG)
	trash    = ui.MustParseSVG(trashSVG)
	folder   = ui.MustParseSVG(folderSVG)
	external = ui.MustParseSVG(externalSVG)
	chevron  = ui.MustParseSVG(chevronSVG)
)

var r = strings.NewReplacer(
	`<`, "",
	`>`, "",
	`:`, "",
	`"`, "",
	`/`, "",
	`\`, "",
	`|`, "",
	`?`, "",
	`*`, "",
)

type game struct {
	Name       string
	SaveRoot   string
	BackupPath string
	Paths      []string
}

type treeItem struct {
	name     string
	open     bool
	checked  bool
	parent   *treeItem
	children []*treeItem
}

type app struct {
	win     *mygo.Window
	config  string
	err     bool
	info    bool
	asking  bool
	message string

	games    []game
	selected int
	list     ui.ListState

	saveRoot   string
	backupPath string

	root *treeItem
}

func newApp() (*app, error) {
	a := app{selected: -1}

	config, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	config = filepath.Join(config, "sav")
	if err := os.MkdirAll(config, 0755); err != nil {
		return nil, err
	}
	a.config = filepath.Join(config, "config.json")

	f, err := os.Open(a.config)
	if err == nil {
		err = json.NewDecoder(f).Decode(&a.games)
		f.Close()
		if err != nil {
			return nil, err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return &a, nil
}

func (a *app) save(games []game) bool {
	f, err := os.Create(a.config)
	if err != nil {
		a.err, a.message = true, err.Error()
		return false
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(games)
	if err != nil {
		a.err, a.message = true, err.Error()
		return false
	}
	return true
}

func (a *app) listSelectedIndexChanged() {
	if a.selected == -1 {
		a.saveRoot = ""
		a.backupPath = ""
		a.root = nil
	} else {
		game := a.games[a.selected]
		a.saveRoot = game.SaveRoot
		a.backupPath = game.BackupPath
		a.rebuildTree()
	}
}

func (a *app) rebuildTree() {
	game := a.games[a.selected]
	if game.SaveRoot == "" {
		a.root = nil
		return
	}

	a.root = &treeItem{name: filepath.Base(game.SaveRoot)}
	var walk func(*treeItem, string, fs.DirEntry) error
	walk = func(item *treeItem, path string, e fs.DirEntry) error {
		if e.IsDir() {
			dirs, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for _, dir := range dirs {
				i := &treeItem{name: dir.Name(), parent: item}
				item.children = append(item.children, i)
				if err = walk(i, filepath.Join(path, dir.Name()), dir); err != nil {
					return err
				}
			}
			if item.children == nil {
				item.children = make([]*treeItem, 0)
			}
		}
		return nil
	}
	info, err := os.Lstat(game.SaveRoot)
	if err != nil {
		a.err, a.message = true, err.Error()
		return
	}
	if err = walk(a.root, game.SaveRoot, fs.FileInfoToDirEntry(info)); err != nil {
		a.err, a.message = true, err.Error()
		return
	}

	walkNodes(a.root, func(i *treeItem) {
		slices.SortFunc(i.children, func(a, b *treeItem) int {
			if a.children != nil {
				if b.children != nil {
					return compareFold(a.name, b.name)
				}
				return -1
			} else if b.children != nil {
				return 1
			}
			return compareFold(a.name, b.name)
		})
	})

	a.root.open = true
	if slices.Contains(game.Paths, "*") {
		walkNodes(a.root, func(i *treeItem) { i.checked = true })
	} else {
		for _, path := range game.Paths {
			cur := a.root
			for part := range strings.SplitSeq(path, string(filepath.Separator)) {
				index := slices.IndexFunc(cur.children, func(i *treeItem) bool {
					return i.name == part
				})
				if index == -1 {
					cur = nil
					break
				}
				cur.open = true
				cur = cur.children[index]
			}
			if cur != nil {
				walkNodes(cur, func(i *treeItem) { i.checked = true })
			}
		}
	}
}

func (a *app) backUp() {
	game := a.games[a.selected]
	if len(game.Paths) == 0 {
		a.info, a.message = true, "No file selected."
		return
	}

	f, err := os.CreateTemp("", "")
	if err != nil {
		a.err, a.message = true, err.Error()
		return
	}

	w := zip.NewWriter(f)
	addFile := func(path string) error {
		file, _ := strings.CutPrefix(path, game.SaveRoot)
		file = strings.TrimLeft(file, string(filepath.Separator))
		dst, err := w.Create(strings.ReplaceAll(file, "\\", "/"))
		if err != nil {
			return err
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(dst, src)
		src.Close()
		return err
	}

	if slices.Contains(game.Paths, "*") {
		err = w.AddFS(os.DirFS(game.SaveRoot))
	} else {
		for _, path := range game.Paths {
			name := filepath.Join(game.SaveRoot, path)
			var info os.FileInfo
			info, err = os.Stat(name)
			if err != nil {
				break
			}

			if info.IsDir() {
				err = filepath.WalkDir(name, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					} else if d.IsDir() {
						return nil
					}
					return addFile(path)
				})
			} else {
				err = addFile(name)
			}
			if err != nil {
				break
			}
		}
	}
	if err := w.Close(); err != nil {
		a.err, a.message = true, err.Error()
	}
	if err := f.Close(); err != nil {
		a.err, a.message = true, err.Error()
	}
	if err != nil {
		a.err, a.message = true, err.Error()
		os.Remove(f.Name())
		return
	}

	dst := filepath.Join(game.BackupPath, r.Replace(game.Name)+".zip")
	if err = os.Rename(f.Name(), dst); err != nil {
		err = copyContent(f.Name(), dst)
		os.Remove(f.Name())
	}
	if err == nil {
		a.info, a.message = true, "Done!"
	} else {
		a.err, a.message = true, err.Error()
	}
}

func (a *app) view(c *ui.Context) {
	a.list.Selected = &a.selected

	t := c.Theme()
	var gray ui.Color
	if t.Dark {
		gray = ui.RGB(33, 37, 41)
	} else {
		gray = ui.RGB(241, 243, 245)
	}

	ui.Row(c).Fill().Children(func() {
		ui.Column(c).Padding(8).Grow(0.3).BorderEndWidth(1).BorderColor(t.Border).AlignSelf(ui.Stretch).Background(gray).Children(func() {
			ui.Row(c).Justify(ui.SpaceBetween).Margin(0, 0, 6, 0).Children(func() {
				ui.Row(c).Gap(4).Children(func() {
					ui.Text(c, "Games").Bold().MarginStart(6)
					ui.Icon(c, question).Size(t.FontSize-2, t.FontSize-2).Tooltip("Double-click the name to rename it")
				})
				if ui.Button(c, "").Children(func() {
					ui.Icon(c, plus)
				}).Clicked() {
					games := append(a.games, game{Name: "New Game"})
					if a.save(games) {
						a.games = games
						a.selected = len(a.games) - 1
						a.listSelectedIndexChanged()
					}
				}
			})
			if ui.List(c, &a.list, len(a.games), func(i int) {
				ui.Row(c).Padding(4, 6).Justify(ui.SpaceBetween).Children(func() {
					if ui.EditableText(c, &a.games[i].Name).Changed() {
						a.save(a.games)
					}
					if i == a.selected && ui.Icon(c, trash).MarginStart(6).Clicked() {
						a.asking = true
					}
				})
			}).Gap(2).Grow(1).Changed() {
				a.listSelectedIndexChanged()
			}
		})

		ui.Column(c).Grow(0.7).AlignSelf(ui.Stretch).Margin(10, 8).Gap(8).Disabled(a.selected == -1).FontSize(14).Children(func() {
			t.Spacing = 2.5
			ui.Row(c).Children(func() {
				ui.Text(c, "Save").MarginEnd(6)
				pick, open := filePicker(c, &a.saveRoot)
				if pick.Clicked() {
					paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
						Parent:    a.win,
						Directory: true,
					})
					if err != nil {
						a.err, a.message = true, err.Error()
					} else if len(paths) > 0 {
						a.games[a.selected].SaveRoot = paths[0]
						if a.save(a.games) {
							a.saveRoot = paths[0]
							a.rebuildTree()
						} else {
							a.games[a.selected].SaveRoot = a.saveRoot
						}
					}
				}
				if open.Clicked() {
					exec.Command("explorer.exe", a.saveRoot).Run()
				}

				ui.Text(c, "to").MarginX(6)
				pick, open = filePicker(c, &a.backupPath)
				if pick.Clicked() {
					paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
						Parent:    a.win,
						Directory: true,
					})
					if err != nil {
						a.err, a.message = true, err.Error()
					} else if len(paths) > 0 {
						a.games[a.selected].BackupPath = paths[0]
						if a.save(a.games) {
							a.backupPath = paths[0]
						} else {
							a.games[a.selected].BackupPath = a.backupPath
						}
					}
				}
				if open.Clicked() {
					name := a.games[a.selected].Name
					backup := filepath.Join(a.backupPath, r.Replace(name)+".zip")
					_, err := os.Stat(backup)
					if err == nil {
						exec.Command("explorer.exe", "/select,", backup).Run()
					} else if errors.Is(err, os.ErrNotExist) {
						exec.Command("explorer.exe", a.backupPath).Run()
					} else {
						a.err, a.message = true, err.Error()
					}
				}
			})

			t.Spacing = 3
			ui.ScrollBoth(c).Background(gray.Alpha(0.5)).Padding(4).Border(1, t.Border).Radius(8).Grow(1).Children(func() {
				tree := ui.TreeBase(c)
				tree.Root.Children(func() {
					if a.root != nil {
						var node func(*treeItem)
						node = func(item *treeItem) {
							var i ui.TreeItemParts
							if item.children == nil {
								i = tree.Item(nil)
							} else {
								i = tree.Item(&item.open)
							}
							ui.Row(c).Gap(4).PaddingStart((14 + 4) * float32(i.Depth())).Children(func() {
								marginStart := float32(14 + 4)
								if item.children != nil {
									marginStart = 0
									arrow := ui.Icon(c, chevron)
									if item.open {
										arrow.Rotate(90)
									}
									if arrow.Clicked() {
										item.open = !item.open
									}
								}
								if ui.Checkbox(c, &item.checked, "").PaddingX(1).MarginStart(marginStart).Changed() {
									walkNodes(item, func(i *treeItem) { i.checked = item.checked })
									if item.checked {
										for cur := item.parent; cur != nil; cur = cur.parent {
											if slices.ContainsFunc(cur.children, func(i *treeItem) bool { return !i.checked }) {
												break
											} else {
												cur.checked = true
											}
										}
									} else {
										for cur := item.parent; cur != nil; cur = cur.parent {
											cur.checked = false
										}
									}

									if a.root.checked {
										a.games[a.selected].Paths = []string{"*"}
									} else {
										a.games[a.selected].Paths = nil
										var walk func(*treeItem)
										walk = func(i *treeItem) {
											if i.checked {
												var paths []string
												for i.parent != nil {
													paths = append(paths, i.name)
													i = i.parent
												}
												slices.Reverse(paths)
												a.games[a.selected].Paths = append(a.games[a.selected].Paths, filepath.Join(paths...))
												return
											}
											for _, i := range i.children {
												walk(i)
											}
										}
										walk(a.root)
									}
									a.save(a.games)
								}
								ui.Text(c, item.name).SingleLine()
							})
							i.Group(func() {
								for _, i := range item.children {
									node(i)
								}
							})
						}
						node(a.root)
					}
				})
			})
			if ui.Button(c, "Back up").Disabled(a.saveRoot == "" || a.backupPath == "").Clicked() {
				a.backUp()
			}
		})
	})

	ui.AlertDialog(c, &a.err, "Error", a.message, "OK")
	ui.AlertDialog(c, &a.info, "Info", a.message, "OK")
	var game game
	if a.selected != -1 {
		game = a.games[a.selected]
	}
	if ui.AlertDialog(c, &a.asking, "Delete "+game.Name+"?", "You can't undo this.", "Cancel", "Delete") == 1 {
		a.games = slices.Delete(a.games, a.selected, a.selected+1)
		if a.save(a.games) {
			a.selected = -1
			a.listSelectedIndexChanged()
		} else {
			a.games = slices.Insert(a.games, a.selected, game)
		}
	}
}

func filePicker(c *ui.Context, value *string) (ui.Element, ui.Element) {
	t := c.Theme()

	input := ui.TextInput(c, value).Grow(1).ReadOnly(true)
	pick := ui.Button(c, "").AlignSelf(ui.Stretch).Radius(0).BorderWidth(1, 0).Children(func() {
		ui.Icon(c, folder)
	})
	open := ui.Button(c, "").AlignSelf(ui.Stretch).Disabled(*value == "").DimDisabled(false).Children(func() {
		ui.Icon(c, external)
	})
	if *value == "" {
		open.Background(t.Surface.Alpha(0.5)).TextColor(t.Text.Alpha(0.5))
	}

	if c.Direction() == ui.LeftToRight {
		input.Radius(t.Radius, 0, 0, t.Radius)
		open.Radius(0, t.Radius, t.Radius, 0)
	} else {
		input.Radius(0, t.Radius, t.Radius, 0)
		open.Radius(t.Radius, 0, 0, t.Radius)
	}
	return pick, open
}
