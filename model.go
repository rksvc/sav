//go:build windows

package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sav/walk"
	"slices"
	"strings"
)

type game struct {
	Name       string
	SaveRoot   string
	BackupPath string
	Paths      []string
}

type model struct {
	config string

	list *listModel
	tree *treeModel

	mw                 *walk.MainWindow
	lb                 *walk.ListBox
	tv                 *walk.TreeView
	name, save, backup *walk.LineEdit
}

func newModel() (*model, error) {
	m := model{list: &listModel{}}

	config, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	config = filepath.Join(config, "sav")
	if err := os.MkdirAll(config, 0755); err != nil {
		return nil, err
	}
	m.config = filepath.Join(config, "config.json")

	f, err := os.Open(m.config)
	if err == nil {
		err = json.NewDecoder(f).Decode(&m.list.games)
		f.Close()
		if err != nil {
			return nil, err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return &m, nil
}

func (m *model) saveConfig(games []game) error {
	f, err := os.Create(m.config)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(games)
}

func (m *model) listSelectedIndexChanged() {
	index := m.lb.CurrentIndex()
	if index == -1 {
		m.try(m.name.SetText(""))
		m.try(m.save.SetText(""))
		m.try(m.backup.SetText(""))
	} else {
		game := m.list.games[index]
		m.try(m.name.SetText(game.Name))
		m.try(m.save.SetText(game.SaveRoot))
		m.try(m.backup.SetText(game.BackupPath))
	}
}

func (m *model) buttonAddGameClick() {
	games := append(m.list.games, game{Name: "New Game"})
	if m.try(m.saveConfig(games)) {
		m.list.games = games
		index := len(m.list.games) - 1
		m.list.PublishItemsInserted(index, index)
		m.try(m.lb.SetCurrentIndex(index))
	}
}

func (m *model) buttonRemoveGameClick() {
	index := m.lb.CurrentIndex()
	game := m.list.games[index]
	if walk.MsgBox(m.mw, "", "Are you sure to delete "+game.Name+"?", walk.MsgBoxYesNo) == walk.DlgCmdYes {
		m.list.games = slices.Delete(m.list.games, index, index+1)
		if m.try(m.saveConfig(m.list.games)) {
			m.list.PublishItemsRemoved(index, index)
			m.try(m.lb.SetCurrentIndex(-1))
		} else {
			m.list.games = slices.Insert(m.list.games, index, game)
		}
	}
}

func (m *model) buttonSaveNameClick() {
	index := m.lb.CurrentIndex()
	oldName := m.list.games[index].Name
	m.list.games[index].Name = m.name.Text()
	if m.try(m.saveConfig(m.list.games)) {
		m.list.PublishItemChanged(index)
	} else {
		m.list.games[index].Name = oldName
	}
}

func (m *model) buttonChooseSaveClick() {
	if folder := m.browseFolder(); folder != "" {
		index := m.lb.CurrentIndex()
		oldSaveRoot := m.list.games[index].SaveRoot
		m.list.games[index].SaveRoot = folder
		if m.try(m.saveConfig(m.list.games)) {
			m.try(m.save.SetText(folder))
		} else {
			m.list.games[index].SaveRoot = oldSaveRoot
		}
	}
}

func (m *model) buttonChooseBackupClick() {
	if folder := m.browseFolder(); folder != "" {
		index := m.lb.CurrentIndex()
		oldBackupPath := m.list.games[index].BackupPath
		m.list.games[index].BackupPath = folder
		if m.try(m.saveConfig(m.list.games)) {
			m.try(m.backup.SetText(folder))
		} else {
			m.list.games[index].BackupPath = oldBackupPath
		}
	}
}

func (m *model) buttonOpenSaveRootClick() {
	exec.Command("explorer.exe", m.save.Text()).Run()
}

func (m *model) buttonOpenBackupClick() {
	backupPath := m.backup.Text()
	name := m.list.games[m.lb.CurrentIndex()].Name
	backup := filepath.Join(backupPath, r.Replace(name)+".zip")
	_, err := os.Stat(backup)
	if err == nil {
		exec.Command("explorer.exe", "/select,", backup).Run()
	} else if errors.Is(err, os.ErrNotExist) {
		exec.Command("explorer.exe", backupPath).Run()
	} else {
		m.error(err)
	}
}

func (m *model) textSaveRootChanged() {
	m.tree = &treeModel{}
	if !m.try(m.tv.SetModel(m.tree)) {
		return
	}
	index := m.lb.CurrentIndex()
	if index == -1 {
		return
	}
	game := m.list.games[index]
	if game.SaveRoot == "" {
		return
	}

	root := &treeItem{name: filepath.Base(game.SaveRoot)}
	if !m.try(filepath.WalkDir(game.SaveRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if path == game.SaveRoot {
			return nil
		}
		path = cutPathPrefix(path, game.SaveRoot)
		cur := root
		for part := range strings.SplitSeq(path, string(filepath.Separator)) {
			index := slices.IndexFunc(cur.children, func(i *treeItem) bool {
				return i.name == part
			})
			if index == -1 {
				i := &treeItem{name: part, parent: cur}
				cur.children = append(cur.children, i)
				cur = i
			} else {
				cur = cur.children[index]
			}
		}
		return nil
	})) {
		return
	}
	walkNodes(root, func(i *treeItem) bool {
		slices.SortFunc(i.children, func(a, b *treeItem) int {
			if len(a.children) > 0 {
				if len(b.children) > 0 {
					return compareFold(a.name, b.name)
				}
				return -1
			} else if len(b.children) > 0 {
				return 1
			}
			return compareFold(a.name, b.name)
		})
		return true
	})
	m.tree = &treeModel{root: root}
	if !m.try(m.tv.SetModel(m.tree)) {
		return
	}

	if !m.try(m.tv.SetExpanded(root, true)) {
		return
	}
	if slices.Contains(game.Paths, "*") {
		if !walkNodes(root, func(i *treeItem) bool { return m.try(m.tv.SetChecked(i, true)) }) {
			return
		}
	} else {
		for _, path := range game.Paths {
			cur := root
			for part := range strings.SplitSeq(path, string(filepath.Separator)) {
				index := slices.IndexFunc(cur.children, func(i *treeItem) bool {
					return i.name == part
				})
				if index == -1 {
					cur = nil
					break
				} else if !m.try(m.tv.SetExpanded(cur, true)) {
					return
				}
				cur = cur.children[index]
			}
			if cur != nil {
				if !walkNodes(cur, func(i *treeItem) bool { return m.try(m.tv.SetChecked(i, true)) }) {
					return
				}
			}
		}
	}
}

func (m *model) treeCheckedChanged(item walk.TreeItem) {
	i := item.(*treeItem)
	checked := m.tv.Checked(i)
	if !walkNodes(i, func(i *treeItem) bool { return m.try(m.tv.SetChecked(i, checked)) }) {
		return
	}
	if checked {
		for cur := i.parent; cur != nil; cur = cur.parent {
			if !slices.ContainsFunc(cur.children, func(i *treeItem) bool {
				return !m.tv.Checked(i)
			}) {
				if !m.try(m.tv.SetChecked(cur, true)) {
					return
				}
			}
		}
	} else {
		for cur := i.parent; cur != nil; cur = cur.parent {
			if !m.try(m.tv.SetChecked(cur, false)) {
				return
			}
		}
	}

	index := m.lb.CurrentIndex()
	if m.tv.Checked(m.tree.root) {
		m.list.games[index].Paths = []string{"*"}
	} else {
		m.list.games[index].Paths = nil
		var walk func(*treeItem)
		walk = func(i *treeItem) {
			if m.tv.Checked(i) {
				var paths []string
				for i.parent != nil {
					paths = append(paths, i.name)
					i = i.parent
				}
				slices.Reverse(paths)
				m.list.games[index].Paths = append(m.list.games[index].Paths, filepath.Join(paths...))
				return
			}
			for _, i := range i.children {
				walk(i)
			}
		}
		walk(m.tree.root)
	}
	m.try(m.saveConfig(m.list.games))
}

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

func (m *model) buttonBackUpClick() {
	game := m.list.games[m.lb.CurrentIndex()]
	if len(game.Paths) == 0 {
		m.info("No file selected.")
		return
	}

	f, err := os.CreateTemp("", "")
	if err != nil {
		m.error(err)
		return
	}

	w := zip.NewWriter(f)
	addFile := func(path string) error {
		dst, err := w.Create(strings.ReplaceAll(cutPathPrefix(path, game.SaveRoot), "\\", "/"))
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
				if err != nil {
					break
				}
			} else {
				err = addFile(name)
			}
		}
	}
	m.try(w.Close())
	m.try(f.Close())
	if err != nil {
		m.error(err)
		os.Remove(f.Name())
		return
	}

	dst := filepath.Join(game.BackupPath, r.Replace(game.Name)+".zip")
	if err = os.Rename(f.Name(), dst); err != nil {
		err = copyContent(f.Name(), dst)
		os.Remove(f.Name())
	}
	if err == nil {
		m.info("Done!")
	} else {
		m.error(err)
	}
}

func (m *model) browseFolder() string {
	var dlg walk.FileDialog
	var owner walk.Form
	if m.mw != nil {
		owner = m.mw
	}
	accept, err := dlg.ShowBrowseFolder(owner)
	if err != nil {
		m.error(err)
		return ""
	} else if accept {
		return dlg.FilePath
	}
	return ""
}

func (m *model) info(msg string) {
	var owner walk.Form
	if m.mw != nil {
		owner = m.mw
	}
	walk.MsgBox(owner, "", msg, walk.MsgBoxIconInformation)
}

func (m *model) error(err error) {
	var owner walk.Form
	if m.mw != nil {
		owner = m.mw
	}
	walk.MsgBox(owner, "", err.Error(), walk.MsgBoxIconError)
}

func (m *model) try(err error) bool {
	if err == nil {
		return true
	}
	m.error(err)
	return false
}
