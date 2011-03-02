// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The filepath package implements utility routines for manipulating
// filename paths in a way compatible with the target operating
// system.
package filepath

import (
	"io/ioutil"
	"os"
	"strings"
)

// BUG(niemeyer): Windows support is missing in Clean, Join, Ext, Walk, Base, IsAbs and Match.

// Clean returns the shortest path name equivalent to path
// by purely lexical processing.  It applies the following rules
// iteratively until no further processing can be done:
//
//	1. Replace multiple Separator elements with a single one.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path:
//	   that is, replace "/.." by "/" at the beginning of a path,
//         assuming Separator is '/'.
//
// If the result of this process is an empty string, Clean
// returns the string ".".
//
// See also Rob Pike, ``Lexical File Names in Plan 9 or
// Getting Dot-Dot right,''
// http://plan9.bell-labs.com/sys/doc/lexnames.html
func Clean(path string) string {
	if path == "" {
		return "."
	}

	rooted := path[0] == Separator
	n := len(path)

	// Invariants:
	//	reading from path; r is index of next byte to process.
	//	writing to buf; w is index of next byte to write.
	//	dotdot is index in buf where .. must stop, either because
	//		it is the leading slash or it is a leading ../../.. prefix.
	buf := []byte(path)
	r, w, dotdot := 0, 0, 0
	if rooted {
		r, w, dotdot = 1, 1, 1
	}

	for r < n {
		switch {
		case path[r] == Separator:
			// empty path element
			r++
		case path[r] == '.' && (r+1 == n || path[r+1] == Separator):
			// . element
			r++
		case path[r] == '.' && path[r+1] == '.' && (r+2 == n || path[r+2] == Separator):
			// .. element: remove to last separator
			r += 2
			switch {
			case w > dotdot:
				// can backtrack
				w--
				for w > dotdot && buf[w] != Separator {
					w--
				}
			case !rooted:
				// cannot backtrack, but not rooted, so append .. element.
				if w > 0 {
					buf[w] = Separator
					w++
				}
				buf[w] = '.'
				w++
				buf[w] = '.'
				w++
				dotdot = w
			}
		default:
			// real path element.
			// add slash if needed
			if rooted && w != 1 || !rooted && w != 0 {
				buf[w] = Separator
				w++
			}
			// copy element
			for ; r < n && path[r] != Separator; r++ {
				buf[w] = path[r]
				w++
			}
		}
	}

	// Turn empty string into "."
	if w == 0 {
		buf[w] = '.'
		w++
	}

	return string(buf[0:w])
}

// Split splits path immediately following the final Separator,
// partitioning it into a directory and a file name components.
// In operating systems where VolumeSeparator is not empty and
// is found in path after any Separator, Split splits the
// volume name from the file name instead.
// If there are no separators in path, Split returns an empty base
// and file set to path.
func Split(path string) (base, file string) {
	i := strings.LastIndex(path, string(Separator))
	return path[:i+1], path[i+1:]
}

// Join joins any number of path elements into a single path, adding
// a Separator if necessary.  All empty strings are ignored.
func Join(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return Clean(strings.Join(elem[i:], string(Separator)))
		}
	}
	return ""
}

// Ext returns the file name extension used by path.
// The extension is the suffix beginning at the final dot
// in the final Separator-partitioned element of path;
// it is empty if there is no dot.
func Ext(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != Separator; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

// Visitor methods are invoked for corresponding file tree entries
// visited by Walk. The parameter path is the full path of f relative
// to root.
type Visitor interface {
	VisitDir(path string, f *os.FileInfo) bool
	VisitFile(path string, f *os.FileInfo)
}

func walk(path string, f *os.FileInfo, v Visitor, errors chan<- os.Error) {
	if !f.IsDirectory() {
		v.VisitFile(path, f)
		return
	}

	if !v.VisitDir(path, f) {
		return // skip directory entries
	}

	list, err := ioutil.ReadDir(path)
	if err != nil {
		if errors != nil {
			errors <- err
		}
	}

	for _, e := range list {
		walk(Join(path, e.Name), e, v, errors)
	}
}

// Walk walks the file tree rooted at root, calling v.VisitDir or
// v.VisitFile for each directory or file in the tree, including root.
// If v.VisitDir returns false, Walk skips the directory's entries;
// otherwise it invokes itself for each directory entry in sorted order.
// An error reading a directory does not abort the Walk.
// If errors != nil, Walk sends each directory read error
// to the channel.  Otherwise Walk discards the error.
func Walk(root string, v Visitor, errors chan<- os.Error) {
	f, err := os.Lstat(root)
	if err != nil {
		if errors != nil {
			errors <- err
		}
		return // can't progress
	}
	walk(root, f, v, errors)
}

// Base returns the last path element of the Separator-partitioned name.
// Trailing Separator elements are removed before extracting the last
// element.  If the name is empty, "." is returned.  If it consists
// entirely of Separator elements, a single Separator is returned.
func Base(name string) string {
	if name == "" {
		return "."
	}
	// Strip trailing slashes.
	for len(name) > 0 && name[len(name)-1] == Separator {
		name = name[0 : len(name)-1]
	}
	// Find the last element
	if i := strings.LastIndex(name, string(Separator)); i >= 0 {
		name = name[i+1:]
	}
	// If empty now, it had only slashes.
	if name == "" {
		return string(Separator)
	}
	return name
}

// IsAbs returns true if the path is absolute.
func IsAbs(path string) bool {
	return len(path) > 0 && path[0] == Separator
}
