//go:build windows

package main

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func copyContent(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()
	_, err = io.Copy(d, s)
	return err
}

func walkNodes(root *treeItem, f func(*treeItem) bool) bool {
	return !slices.ContainsFunc(root.children, func(i *treeItem) bool {
		return !f(i) || !walkNodes(i, f)
	})
}

func cutPathPrefix(path, prefix string) string {
	path, _ = strings.CutPrefix(path, prefix)
	return strings.TrimLeft(path, string(filepath.Separator))
}
