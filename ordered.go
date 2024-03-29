// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ordered illustrates a workaround for the GeneralAbsDifference example in the
// Type Parameters draft, and the limitations of that workaround.
//
// Unfortunately, that workaround itself no longer works due to
// bugs and feature reductions in the implementation of generics
// that occurred after the proposal process.
package main

import (
	"fmt"
	"math"
)

// Numeric is a constraint that matches any numeric type.
// It would likely be in a constraints package in the standard library.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~complex64 | ~complex128
}

// NumericAbs matches numeric types with an Abs method.
type NumericAbs[T any] interface {
	Numeric
	Abs() T
}

// AbsDifference computes the absolute value of the difference of
// a and b, where the absolute value is determined by the Abs method.
func AbsDifference[T NumericAbs[T]](a, b T) T {
	d := a - b
	return d.Abs()
}

// OrderedNumeric matches numeric types that support the < operator.
type OrderedNumeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Complex matches the two complex types, which do not have a < operator.
type Complex interface {
	~complex64 | ~complex128
}

// OrderedAbs is a helper type that defines an Abs method for
// ordered numeric types.
type OrderedAbs[T OrderedNumeric] T

func (a OrderedAbs[T]) Abs() OrderedAbs[T] {
	if a.V < 0 {
		return OrderedAbs[T](-a)
	}
	return OrderedAbs[T](a)
}

// ComplexAbs is a helper type that defines an Abs method for
// complex types.
type ComplexAbs[T Complex] T

func (a ComplexAbs[T]) Abs() ComplexAbs[T] {
	d := math.Hypot(float64(real(a)), float64(imag(a)))
	return ComplexAbs[T]{T(complex(d, 0))}
}

// OrderedAbsDifference returns the absolute value of the difference
// between a and b, where a and b are of an ordered type.
func OrderedAbsDifference[T OrderedNumeric](a, b T) T {
	return T(AbsDifference[OrderedAbs[T]](OrderedAbs[T](a), OrderedAbs[T](b)))
}

// ComplexAbsDifference returns the absolute value of the difference
// between a and b, where a and b are of a complex type.
func ComplexAbsDifference[T Complex](a, b T) T {
	return T(AbsDifference[ComplexAbs[T]](ComplexAbs[T](a), ComplexAbs[T](b)))
}

func asSame[T any](_ T, b interface{}) T {
	return b.(T)
}

// GeneralAbsDifference implements AbsDifference for any *built-in* numeric type T.
//
// However, it panics for defined numeric types that are not built-in:
// handling those cases under the current design would require the use of reflection.
func GeneralAbsDifference[T Numeric](a, b T) T {
	switch a := (interface{})(a).(type) {
	case int:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case int8:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case int16:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case int32:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case int64:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case uint:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case uint8:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case uint16:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case uint32:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case uint64:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case uintptr:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case float32:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))
	case float64:
		return asSame(b, OrderedAbsDifference(a, asSame(a, b)))

	case complex64:
		return asSame(b, ComplexAbsDifference(a, asSame(a, b)))
	case complex128:
		return asSame(b, ComplexAbsDifference(a, asSame(a, b)))

	default:
		panic(fmt.Sprintf("%T is not a builtin numeric type", a))
	}
}

type MyInt int

func main() {
	fmt.Println(GeneralAbsDifference(42, 64))
	fmt.Println(GeneralAbsDifference(42+3i, 64+0i))
	fmt.Println(GeneralAbsDifference(MyInt(42), MyInt(64)))
}
