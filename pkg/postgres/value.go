package postgres

import "strings"

// Integer
type intVal struct {
	v int
}

func NewInt(v int) Value {
	return intVal{v: v}
}

// always valid
func (v intVal) Valid() bool {
	return true
}

func (v intVal) Arg() interface{} {
	return v.v
}

// String
type stringVal struct {
	v string
}

func NewString(v string) Value {
	return stringVal{v: v}
}

func (v stringVal) Valid() bool {
	return strings.TrimSpace(v.v) != ""
}

func (v stringVal) Arg() interface{} {
	return v.v
}

// Slice string
type sliceStringVal struct {
	v []string
}

func NewSliceString(v []string) Value {
	return sliceStringVal{v: v}
}

func (v sliceStringVal) Valid() bool {
	return len(v.v) != 0
}

func (v sliceStringVal) Arg() interface{} {
	return v.v
}
