// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// nomono is a non-monomorphizable program under the Type Parameters draft
// design.
package main

import "fmt"

type Expander[a any] struct{
	x *Expander[Expander[a]]
}

func (x Expander[a]) Expand() interface{} {
	return x.x
}

func main() {
	var a interface{} = Expander[struct{}]{}.x
	fmt.Sprintf("Expander[struct{}]{}.Expand(): %T\n", a)
}
