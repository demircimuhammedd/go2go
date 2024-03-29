// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import "fmt"

// Fig. 6

// Eval on Num
type Evaler interface {
	Eval() int
}
type Num struct {
	value int
}
func (e Num) Eval() int {
	return e.value
}
// Eval on Plus
type Plus[a Expr] struct {
	left a
	right a
}
func (e Plus[a]) Eval() int {
	return e.left.Eval() + e.right.Eval()
}
// String on Num
type Stringer interface {
	String() string
}
func (e Num) String() string {
	return fmt.Sprintf("%d", e.value)
}
// String on Plus
func (e Plus[a]) String() string {
	return fmt.Sprintf("%s+%s", e.left.String(), e.right.String())
}
// tie it all together
type Expr interface {
	Evaler
	Stringer
}

func main() {
	var e Expr = Plus[Expr]{Num{1}, Num{2}}
	var v int = e.Eval() // 3
	var s string = e.String() //"(1+2)"

	fmt.Println(v)
	fmt.Println(s)
}
