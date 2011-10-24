# Modified the basic makefiles referred to from the
# Go home page.
#
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=mflag
GOFILES=\
		errors.go\
		strconv.go\
		flagset.go\
        mflag.go\

include $(GOROOT)/src/Make.pkg

ex.%:
	bash -c 'for d in examples/*; do cd $$d && echo "$@" | sed "s/ex\.//" | xargs gomake && cd -; done'
