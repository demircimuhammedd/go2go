// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import "fmt"

// Featherweight Go, Fig. 3

type Function[a any, b any] interface {
	Apply(x a) b
}
type incr struct{ n int }
func (this incr) Apply(x int) int {
	return x + this.n
}
type pos struct{}
func (this pos) Apply(x int) bool {
	return x > 0
}
type compose[a any, b any, c any] struct {
	f Function[a, b]
	g Function[b, c]
}
func (this compose[a, b, c]) Apply(x a) c {
	return this.g.Apply(this.f.Apply(x))
}
func main() {
	var f Function[int, bool] = compose[int, int, bool]{incr{-5}, pos{}}
	var b bool = f.Apply(3) // false

	fmt.Println(b)
}
