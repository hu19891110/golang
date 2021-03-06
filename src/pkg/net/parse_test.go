// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"bufio"
	"os"
	"testing"
	"runtime"
)

func TestReadLine(t *testing.T) {
	// /etc/services file does not exist on windows and Plan 9.
	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		return
	}
	filename := "/etc/services" // a nice big file

	fd, err := os.Open(filename)
	if err != nil {
		t.Fatalf("open %s: %v", filename, err)
	}
	br := bufio.NewReader(fd)

	file, err := open(filename)
	if file == nil {
		t.Fatalf("net.open(%s) = nil", filename)
	}

	lineno := 1
	byteno := 0
	for {
		bline, berr := br.ReadString('\n')
		if n := len(bline); n > 0 {
			bline = bline[0 : n-1]
		}
		line, ok := file.readLine()
		if (berr != nil) != !ok || bline != line {
			t.Fatalf("%s:%d (#%d)\nbufio => %q, %v\nnet => %q, %v",
				filename, lineno, byteno, bline, berr, line, ok)
		}
		if !ok {
			break
		}
		lineno++
		byteno += len(line) + 1
	}
}
