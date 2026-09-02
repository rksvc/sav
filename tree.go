//go:build windows

package main

import "sav/walk"

type treeItem struct {
	name     string
	parent   *treeItem
	children []*treeItem
}

func (i *treeItem) Text() string {
	return i.name
}

func (i *treeItem) Parent() walk.TreeItem {
	if i.parent == nil {
		return nil
	}
	return i.parent
}

func (i *treeItem) ChildCount() int {
	return len(i.children)
}

func (i *treeItem) ChildAt(index int) walk.TreeItem {
	return i.children[index]
}

type treeModel struct {
	walk.TreeModelBase
	root *treeItem
}

func (m *treeModel) RootCount() int {
	if m.root == nil {
		return 0
	}
	return 1
}

func (m *treeModel) RootAt(index int) walk.TreeItem {
	if index == 0 {
		return m.root
	}
	return nil
}
