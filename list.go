//go:build windows

package main

import "sav/walk"

type listModel struct {
	walk.ListModelBase
	games []game
}

func (m *listModel) ItemCount() int {
	return len(m.games)
}

func (m *listModel) Value(index int) any {
	return m.games[index].Name
}
