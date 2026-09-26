//go:build windows

package main

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
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
	return f(root) && !slices.ContainsFunc(root.children, func(i *treeItem) bool { return !walkNodes(i, f) })
}

func compareFold(a, b string) int {
	for {
		ra, sa := utf8.DecodeRuneInString(a)
		rb, sb := utf8.DecodeRuneInString(b)
		if sa == 0 {
			if sb == 0 {
				return 0
			}
			return -1
		} else if sb == 0 {
			return 1
		}

		ra = unicode.ToLower(ra)
		rb = unicode.ToLower(rb)
		if ra != rb {
			if ra > rb {
				return 1
			}
			return -1
		}

		a = a[sa:]
		b = b[sb:]
	}
}

func cutPathPrefix(path, prefix string) string {
	path, _ = strings.CutPrefix(path, prefix)
	return strings.TrimLeft(path, string(filepath.Separator))
}
