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

//go:embed icons/chevron-down.svg
var chevronDownSVG []byte

//go:embed icons/chevron-right.svg
var chevronRightSVG []byte

var (
	chevronDown  = ui.MustParseSVG(chevronDownSVG)
	chevronRight = ui.MustParseSVG(chevronRightSVG)
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
	expanded bool
	checked  bool
	parent   *treeItem
	children []*treeItem
}

type app struct {
	win    *mygo.Window
	config string

	games    []game
	selected int
	list     ui.ListState

	name       string
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
		mygo.Dialog.Error("", err.Error())
		return false
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(games)
	if err != nil {
		mygo.Dialog.Error("", err.Error())
		return false
	}
	return true
}

func (a *app) listSelectedIndexChanged() {
	if a.selected == -1 {
		a.name = ""
		a.saveRoot = ""
		a.backupPath = ""
		a.root = nil
	} else {
		game := a.games[a.selected]
		a.name = game.Name
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
		mygo.Dialog.Error("", err.Error())
		return
	}
	if err = walk(a.root, game.SaveRoot, fs.FileInfoToDirEntry(info)); err != nil {
		mygo.Dialog.Error("", err.Error())
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

	a.root.expanded = true
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
				cur.expanded = true
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
		mygo.Dialog.Message(mygo.MessageOptions{
			Parent:  a.win,
			Message: "No file selected.",
		})
		return
	}

	f, err := os.CreateTemp("", "")
	if err != nil {
		mygo.Dialog.Error("", err.Error())
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
		mygo.Dialog.Error("", err.Error())
	}
	if err := f.Close(); err != nil {
		mygo.Dialog.Error("", err.Error())
	}
	if err != nil {
		mygo.Dialog.Error("", err.Error())
		os.Remove(f.Name())
		return
	}

	dst := filepath.Join(game.BackupPath, r.Replace(game.Name)+".zip")
	if err = os.Rename(f.Name(), dst); err != nil {
		err = copyContent(f.Name(), dst)
		os.Remove(f.Name())
	}
	if err == nil {
		mygo.Dialog.Message(mygo.MessageOptions{
			Parent:  a.win,
			Message: "Done!",
		})
	} else {
		mygo.Dialog.Error("", err.Error())
	}
}

func (a *app) view(c *ui.Context) {
	a.list.Label = func(i int) string { return a.games[i].Name }
	a.list.Selected = &a.selected

	t := c.Theme()
	t.Spacing = 3

	ui.Column(c).Fill().Padding(8).Gap(4).Children(func() {
		ui.Text(c, "Games")
		ui.Row(c).Grow(1).Gap(6).Children(func() {
			ui.Column(c).Grow(0.3).AlignSelf(ui.Stretch).Gap(6).Children(func() {
				if ui.List(c, &a.list, len(a.games), func(i int) {
					ui.Text(c, a.games[i].Name).Padding(4, 6)
				}).Padding(3).Border(1, t.Border).Radius(8).Gap(3).Grow(1).Changed() {
					a.listSelectedIndexChanged()
				}
				ui.Row(c).Gap(6).Children(func() {
					if ui.Button(c, "Add").Grow(0.5).Clicked() {
						games := append(a.games, game{Name: "New Game"})
						if a.save(games) {
							a.games = games
							a.selected = len(a.games) - 1
							a.listSelectedIndexChanged()
						}
					}
					if ui.Button(c, "Delete").Grow(0.5).Disabled(a.selected == -1).Clicked() {
						game := a.games[a.selected]
						res, err := mygo.Dialog.Message(mygo.MessageOptions{
							Parent:  a.win,
							Type:    mygo.MessageWarning,
							Message: "Delete " + game.Name + "?",
							Buttons: []string{"Delete", "Cancel"},
						})
						if err != nil {
							mygo.Dialog.Error("", err.Error())
						} else if res.Button == 0 {
							a.games = slices.Delete(a.games, a.selected, a.selected+1)
							if a.save(a.games) {
								a.selected = -1
								a.listSelectedIndexChanged()
							} else {
								a.games = slices.Insert(a.games, a.selected, game)
							}
						}
					}
				})
			})
			ui.Column(c).Grow(0.7).AlignSelf(ui.Stretch).Gap(6).Disabled(a.selected == -1).Children(func() {
				ui.Form(c, func() {
					ui.Field(c, "Name", func() {
						ui.Row(c).Gap(6).Children(func() {
							ui.TextInput(c, &a.name).Grow(1)
							if ui.Button(c, "Save").Clicked() {
								oldName := a.games[a.selected].Name
								a.games[a.selected].Name = a.name
								if !a.save(a.games) {
									a.games[a.selected].Name = oldName
								}
							}
						})
					})
					ui.Field(c, "Save Root", func() {
						ui.Row(c).Gap(6).Children(func() {
							ui.TextInput(c, &a.saveRoot).Grow(1).ReadOnly(true)
							if ui.Button(c, "Choose").Clicked() {
								paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
									Parent:    a.win,
									Directory: true,
								})
								if err != nil {
									mygo.Dialog.Error("", err.Error())
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
							if ui.Button(c, "Open").Disabled(a.saveRoot == "").Clicked() {
								exec.Command("explorer.exe", a.saveRoot).Run()
							}
						})
					})
					ui.Field(c, "Backup Path", func() {
						ui.Row(c).Gap(6).Children(func() {
							ui.TextInput(c, &a.backupPath).Grow(1).ReadOnly(true)
							if ui.Button(c, "Choose").Clicked() {
								paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
									Parent:    a.win,
									Directory: true,
								})
								if err != nil {
									mygo.Dialog.Error("", err.Error())
								} else if len(paths) > 0 {
									a.games[a.selected].BackupPath = paths[0]
									if a.save(a.games) {
										a.backupPath = paths[0]
									} else {
										a.games[a.selected].BackupPath = a.backupPath
									}
								}
							}
							if ui.Button(c, "Open").Disabled(a.backupPath == "").Clicked() {
								name := a.games[a.selected].Name
								backup := filepath.Join(a.backupPath, r.Replace(name)+".zip")
								_, err := os.Stat(backup)
								if err == nil {
									exec.Command("explorer.exe", "/select,", backup).Run()
								} else if errors.Is(err, os.ErrNotExist) {
									exec.Command("explorer.exe", a.backupPath).Run()
								} else {
									mygo.Dialog.Error("", err.Error())
								}
							}
						})
					})
				}).Gap(6)
				ui.ScrollBoth(c).Padding(4).Border(1, t.Border).Radius(8).Grow(1).Children(func() {
					if a.root != nil {
						ui.Column(c).Children(func() {
							a.tree(c, a.root, 0)
						})
					}
				})
				if ui.Button(c, "Back up").Disabled(a.saveRoot == "" || a.backupPath == "").Clicked() {
					a.backUp()
				}
			})
		})
	})
}

func (a *app) tree(c *ui.Context, item *treeItem, depth float32) {
	ui.Row(c).Gap(4).PaddingStart((14 + 4) * depth).Children(func() {
		marginStart := float32(14 + 4)
		if item.children != nil {
			marginStart = 0
			if ui.ButtonBase(c).Children(func() {
				if item.expanded {
					ui.Icon(c, chevronDown)
				} else {
					ui.Icon(c, chevronRight)
				}
			}).Clicked() {
				item.expanded = !item.expanded
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
	if item.expanded {
		for _, i := range item.children {
			a.tree(c, i, depth+1)
		}
	}
}
