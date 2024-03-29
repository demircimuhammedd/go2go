// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafeslice_test

import (
	"fmt"
	"testing"

	"github.com/bcmills/go2go/unsafeslice"
)

// asCPointer returns b as a C-style pointer and length
func asCPointer(b []byte) (*byte, int) {
	if len(b) == 0 {
		return nil, 0
	}
	return &b[0], len(b)
}

func ExampleOf() {
	original := []byte("Hello, world!")
	p, n := asCPointer(original)

	alias := unsafeslice.Of(p, n)

	fmt.Printf("original: %s\n", original)
	fmt.Printf("alias: %s\n", alias)
	copy(alias, "Adios")
	fmt.Printf("original: %s\n", original)
	fmt.Printf("alias: %s\n", alias)

	// Output:
	// original: Hello, world!
	// alias: Hello, world!
	// original: Adios, world!
	// alias: Adios, world!
}

func ExampleConvert() {
	// For this example, we're going to do a transformation on some ASCII text.
	// That transformation is not endian-sensitive, so we can reinterpret the text
	// as a slice of uint32s to process it word-at-a-time instead of
	// byte-at-a-time.

	const input = "HELLO, WORLD!"

	// Allocate an aligned backing buffer.
	buf := make([]uint32, (len(input)+3)/4)

	// Reinterpret it as a byte slice so that we can copy in our text.
	// The call site here is awkward because we have to specify both types,
	// even though the source type can be inferred.
	alias := unsafeslice.Convert[uint32, byte](buf)
	copy(alias, input)

	// Perform an endian-insensitive transformation word-by-word instead of
	// byte-by-byte.
	for i := range buf {
		buf[i] |= 0x20202020
	}

	// Read the result back out of the byte-slice view to interpret it as text.
	fmt.Printf("%s\n", alias[:len(input)])

	// Output:
	// hello, world!
}

func ExampleConvertAt() {
	// For this example, we're going to do a transformation on some ASCII text.
	// That transformation is not endian-sensitive, so we can reinterpret the text
	// as a slice of uint32s to process it word-at-a-time instead of
	// byte-at-a-time.

	const input = "HELLO, WORLD!"

	// Allocate an aligned backing buffer.
	buf := make([]uint32, (len(input)+3)/4)

	// Reinterpret it as a byte slice so that we can copy in our text.
	var alias []byte
	unsafeslice.ConvertAt(&alias, buf)
	copy(alias, input)

	// Perform an endian-insensitive transformation word-by-word instead of
	// byte-by-byte.
	for i := range buf {
		buf[i] |= 0x20202020
	}

	// Read the result back out of the byte-slice view to interpret it as text.
	fmt.Printf("%s\n", alias[:len(input)])

	// Output:
	// hello, world!
}

type big [1 << 20]byte

func TestOfWithVeryLargeTypeDoesNotPanic(t *testing.T) {
	var x big
	_ = unsafeslice.Of(&x, 1)
}

func TestConvertAt(t *testing.T) {
	u32 := []uint32{0x00102030, 0x40506070}[:1]
	var b []byte
	unsafeslice.ConvertAt(&b, u32)

	if want := len(u32) * 4; len(b) != want {
		t.Errorf("ConvertAt(_, %x): length = %v; want %v", u32, len(b), want)
	}
	if want := cap(u32) * 4; cap(b) != want {
		t.Errorf("ConvertAt(_, %x): capacity = %v; want %v", u32, cap(b), want)
	}
}

func TestConvertAtErrors(t *testing.T) {
	cases := []struct {
		desc string
		dst  *[]uint32
		src  []byte
	}{
		{
			desc: "incompatible capacity",
			src:  []byte("foobar")[:4:6],
			dst:  new([]uint32),
		},
		{
			desc: "incompatible length",
			src:  []byte("foobar\x00\x00")[:6],
			dst:  new([]uint32),
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				if msg := recover(); msg != nil {
					t.Logf("recovered: %v", msg)
				} else {
					t.Errorf("ConvertAt failed to panic as expected.")
				}
			}()

			unsafeslice.ConvertAt(tc.dst, tc.src)
		})
	}
}

type sliceOfBig []big

func TestConvertAtTypeInference(t *testing.T) {
	var b []byte
	var s sliceOfBig

	// This call requires an explicit type conversion. Otherwise, it fails with the error:
	// 	cannot use &s (value of type *sliceOfBig) as *[]big value in argument
	unsafeslice.ConvertAt((*[]big)(&s), b)
}
