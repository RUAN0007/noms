// This file was generated by nomdl/codegen.

package test

import (
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __testPackageInFile_struct_CachedRef = __testPackageInFile_struct_Ref()

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func __testPackageInFile_struct_Ref() ref.Ref {
	p := types.PackageDef{
		NamedTypes: types.MapOfStringToTypeRefDef{

			"Struct": types.MakeStructTypeRef("Struct",
				[]types.Field{
					types.Field{"s", types.MakePrimitiveTypeRef(types.StringKind), false},
					types.Field{"b", types.MakePrimitiveTypeRef(types.BoolKind), false},
				},
				types.Choices{},
			),
		},
	}.New()
	return types.RegisterPackage(&p)
}

// ListOfStruct

type ListOfStruct struct {
	l types.List
}

func NewListOfStruct() ListOfStruct {
	return ListOfStruct{types.NewList()}
}

type ListOfStructDef []StructDef

func (def ListOfStructDef) New() ListOfStruct {
	l := make([]types.Value, len(def))
	for i, d := range def {
		l[i] = d.New().NomsValue()
	}
	return ListOfStruct{types.NewList(l...)}
}

func (l ListOfStruct) Def() ListOfStructDef {
	d := make([]StructDef, l.Len())
	for i := uint64(0); i < l.Len(); i++ {
		d[i] = StructFromVal(l.l.Get(i)).Def()
	}
	return d
}

func ListOfStructFromVal(val types.Value) ListOfStruct {
	// TODO: Validate here
	return ListOfStruct{val.(types.List)}
}

func (l ListOfStruct) NomsValue() types.Value {
	return l.l
}

func (l ListOfStruct) Equals(other types.Value) bool {
	if other, ok := other.(ListOfStruct); ok {
		return l.l.Equals(other.l)
	}
	return false
}

func (l ListOfStruct) Ref() ref.Ref {
	return l.l.Ref()
}

func (l ListOfStruct) Chunks() []types.Future {
	return l.l.Chunks()
}

// A Noms Value that describes ListOfStruct.
var __typeRefForListOfStruct = types.MakeCompoundTypeRef("", types.ListKind, types.MakeTypeRef("Struct", __testPackageInFile_struct_CachedRef))

func (m ListOfStruct) TypeRef() types.TypeRef {
	return __typeRefForListOfStruct
}

func init() {
	types.RegisterFromValFunction(__typeRefForListOfStruct, func(v types.Value) types.NomsValue {
		return ListOfStructFromVal(v)
	})
}

func (l ListOfStruct) Len() uint64 {
	return l.l.Len()
}

func (l ListOfStruct) Empty() bool {
	return l.Len() == uint64(0)
}

func (l ListOfStruct) Get(i uint64) Struct {
	return StructFromVal(l.l.Get(i))
}

func (l ListOfStruct) Slice(idx uint64, end uint64) ListOfStruct {
	return ListOfStruct{l.l.Slice(idx, end)}
}

func (l ListOfStruct) Set(i uint64, val Struct) ListOfStruct {
	return ListOfStruct{l.l.Set(i, val.NomsValue())}
}

func (l ListOfStruct) Append(v ...Struct) ListOfStruct {
	return ListOfStruct{l.l.Append(l.fromElemSlice(v)...)}
}

func (l ListOfStruct) Insert(idx uint64, v ...Struct) ListOfStruct {
	return ListOfStruct{l.l.Insert(idx, l.fromElemSlice(v)...)}
}

func (l ListOfStruct) Remove(idx uint64, end uint64) ListOfStruct {
	return ListOfStruct{l.l.Remove(idx, end)}
}

func (l ListOfStruct) RemoveAt(idx uint64) ListOfStruct {
	return ListOfStruct{(l.l.RemoveAt(idx))}
}

func (l ListOfStruct) fromElemSlice(p []Struct) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = v.NomsValue()
	}
	return r
}

type ListOfStructIterCallback func(v Struct, i uint64) (stop bool)

func (l ListOfStruct) Iter(cb ListOfStructIterCallback) {
	l.l.Iter(func(v types.Value, i uint64) bool {
		return cb(StructFromVal(v), i)
	})
}

type ListOfStructIterAllCallback func(v Struct, i uint64)

func (l ListOfStruct) IterAll(cb ListOfStructIterAllCallback) {
	l.l.IterAll(func(v types.Value, i uint64) {
		cb(StructFromVal(v), i)
	})
}

type ListOfStructFilterCallback func(v Struct, i uint64) (keep bool)

func (l ListOfStruct) Filter(cb ListOfStructFilterCallback) ListOfStruct {
	nl := NewListOfStruct()
	l.IterAll(func(v Struct, i uint64) {
		if cb(v, i) {
			nl = nl.Append(v)
		}
	})
	return nl
}

// Struct

type Struct struct {
	m types.Map
}

func NewStruct() Struct {
	return Struct{types.NewMap(
		types.NewString("$name"), types.NewString("Struct"),
		types.NewString("$type"), types.MakeTypeRef("Struct", __testPackageInFile_struct_CachedRef),
		types.NewString("s"), types.NewString(""),
		types.NewString("b"), types.Bool(false),
	)}
}

type StructDef struct {
	S string
	B bool
}

func (def StructDef) New() Struct {
	return Struct{
		types.NewMap(
			types.NewString("$name"), types.NewString("Struct"),
			types.NewString("$type"), types.MakeTypeRef("Struct", __testPackageInFile_struct_CachedRef),
			types.NewString("s"), types.NewString(def.S),
			types.NewString("b"), types.Bool(def.B),
		)}
}

func (s Struct) Def() (d StructDef) {
	d.S = s.m.Get(types.NewString("s")).(types.String).String()
	d.B = bool(s.m.Get(types.NewString("b")).(types.Bool))
	return
}

var __typeRefForStruct = types.MakeTypeRef("Struct", __testPackageInFile_struct_CachedRef)

func (m Struct) TypeRef() types.TypeRef {
	return __typeRefForStruct
}

func init() {
	types.RegisterFromValFunction(__typeRefForStruct, func(v types.Value) types.NomsValue {
		return StructFromVal(v)
	})
}

func StructFromVal(val types.Value) Struct {
	// TODO: Validate here
	return Struct{val.(types.Map)}
}

func (s Struct) NomsValue() types.Value {
	return s.m
}

func (s Struct) Equals(other types.Value) bool {
	if other, ok := other.(Struct); ok {
		return s.m.Equals(other.m)
	}
	return false
}

func (s Struct) Ref() ref.Ref {
	return s.m.Ref()
}

func (s Struct) Chunks() []types.Future {
	return s.m.Chunks()
}

func (s Struct) S() string {
	return s.m.Get(types.NewString("s")).(types.String).String()
}

func (s Struct) SetS(val string) Struct {
	return Struct{s.m.Set(types.NewString("s"), types.NewString(val))}
}

func (s Struct) B() bool {
	return bool(s.m.Get(types.NewString("b")).(types.Bool))
}

func (s Struct) SetB(val bool) Struct {
	return Struct{s.m.Set(types.NewString("b"), types.Bool(val))}
}