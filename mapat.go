// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// mapat illustrates the closest API to a conventional functional "Map"
// that I could express under the Type Parameters draft design.
package main

import (
	"context"
	"fmt"
	"syscall"
)

// MapAt copies all of the elements in src to dst, using the supplied function
// to convert between the two.
func MapAt[T2, T1 any](dst []T2, src []T1, convert func(T1) T2) int {
	for i, x := range src {
		if i > len(dst) {
			return i
		}
		dst[i] = convert(x)
	}
	return len(src)
}

func main() {
	errnos := []syscall.Errno{syscall.ENOSYS, syscall.EINVAL}
	errs := make([]error, len(errnos))

	MapAt(errs, errnos, func(err syscall.Errno) error { return err })

	for _, err := range errs {
		fmt.Println(err)
	}

	cancelers := []context.CancelFunc{
		func(){fmt.Println("cancel 1")},
		func(){fmt.Println("cancel 2")},
	}
	funcs := make([]func(), len(cancelers))

	MapAt(funcs, cancelers, func(f context.CancelFunc) func() { return f })

	for _, f := range funcs {
		f()
	}
}
