# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include ../../Make.inc
O:=$(HOST_O)

TARG=6c

HFILES=\
	gc.h\
	../6l/6.out.h\
	../cc/cc.h\

OFILES=\
	cgen.$O\
	list.$O\
	sgen.$O\
	swt.$O\
	txt.$O\
	pgen.$O\
	pswt.$O\
	div.$O\
	mul.$O\
	reg.$O\
	peep.$O\
	machcap.$O\
	../6l/enam.$O\

LIB=\
	../cc/cc.a\

include ../../Make.ccmd

%.$O: ../cc/%.c
	$(HOST_CC) $(HOST_CFLAGS) -c -I. -o $@ ../cc/$*.c
