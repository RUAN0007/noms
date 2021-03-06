// Copyright 2017 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package ngql

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/attic-labs/noms/go/chunks"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/testify/suite"
)

type QueryGraphQLSuite struct {
	suite.Suite
	vs *types.ValueStore
}

func TestQueryGraphQL(t *testing.T) {
	suite.Run(t, &QueryGraphQLSuite{})
}

func (suite *QueryGraphQLSuite) SetupTest() {
	cs := chunks.NewTestStore()
	suite.vs = types.NewValueStore(types.NewBatchStoreAdaptor(cs))
}

func (suite *QueryGraphQLSuite) assertQueryResult(v types.Value, q, expect string) {
	buff := &bytes.Buffer{}
	Query(v, q, suite.vs, buff)
	suite.Equal(expect+"\n", string(buff.Bytes()))
}

func (suite *QueryGraphQLSuite) TestScalars() {
	suite.assertQueryResult(types.String("aaa"), "{root}", `{"data":{"root":"aaa"}}`)
	suite.assertQueryResult(types.String(""), "{root}", `{"data":{"root":""}}`)

	suite.assertQueryResult(types.Number(0), "{root}", `{"data":{"root":0}}`)
	suite.assertQueryResult(types.Number(1), "{root}", `{"data":{"root":1}}`)
	suite.assertQueryResult(types.Number(-1), "{root}", `{"data":{"root":-1}}`)
	suite.assertQueryResult(types.Number(1<<31), "{root}", `{"data":{"root":2.147483648e+09}}`)
	suite.assertQueryResult(types.Number(-(1 << 31)), "{root}", `{"data":{"root":-2.147483648e+09}}`)
	suite.assertQueryResult(types.Number(0.001), "{root}", `{"data":{"root":0.001}}`)
	suite.assertQueryResult(types.Number(0.00000001), "{root}", `{"data":{"root":1e-08}}`)

	suite.assertQueryResult(types.Bool(false), "{root}", `{"data":{"root":false}}`)
	suite.assertQueryResult(types.Bool(true), "{root}", `{"data":{"root":true}}`)
}

func (suite *QueryGraphQLSuite) TestStructBasic() {
	s1 := types.NewStruct("Foo", types.StructData{
		"a": types.String("aaa"),
		"b": types.Bool(true),
		"c": types.Number(0.1),
	})

	suite.assertQueryResult(s1, "{root{a}}", `{"data":{"root":{"a":"aaa"}}}`)
	suite.assertQueryResult(s1, "{root{a b}}", `{"data":{"root":{"a":"aaa","b":true}}}`)
	suite.assertQueryResult(s1, "{root{a b c}}", `{"data":{"root":{"a":"aaa","b":true,"c":0.1}}}`)
	suite.assertQueryResult(s1, "{root{a c}}", `{"data":{"root":{"a":"aaa","c":0.1}}}`)
}

func (suite *QueryGraphQLSuite) TestEmptyStruct() {
	s1 := types.NewStruct("", types.StructData{})

	suite.assertQueryResult(s1, "{root{hash}}", `{"data":{"root":{"hash":"c66c33bb6na2m5mk0bek7eqqrl2t7gmv"}}}`)
}

func (suite *QueryGraphQLSuite) TestEmbeddedStruct() {
	s1 := types.NewStruct("Foo", types.StructData{
		"a": types.String("aaa"),
		"b": types.NewStruct("Bar", types.StructData{
			"c": types.Bool(true),
			"d": types.Number(0.1),
		}),
	})

	suite.assertQueryResult(s1, "{root{a}}", `{"data":{"root":{"a":"aaa"}}}`)
	suite.assertQueryResult(s1, "{root{a b {c}}}", `{"data":{"root":{"a":"aaa","b":{"c":true}}}}`)
	suite.assertQueryResult(s1, "{root{a b {c d}}}", `{"data":{"root":{"a":"aaa","b":{"c":true,"d":0.1}}}}`)
}

func (suite *QueryGraphQLSuite) TestListBasic() {
	for _, valuesKey := range []string{"elements", "values"} {
		list := types.NewList()
		suite.assertQueryResult(list, "{root{size}}", `{"data":{"root":{"size":0}}}`)
		suite.assertQueryResult(list, "{root{"+valuesKey+"}}", `{"data":{"root":{}}}`)

		list = types.NewList(types.String("foo"), types.String("bar"), types.String("baz"))

		suite.assertQueryResult(list, "{root{"+valuesKey+"}}", `{"data":{"root":{"`+valuesKey+`":["foo","bar","baz"]}}}`)
		suite.assertQueryResult(list, "{root{size}}", `{"data":{"root":{"size":3}}}`)
		suite.assertQueryResult(list, "{root{"+valuesKey+"(at:1,count:2)}}", `{"data":{"root":{"`+valuesKey+`":["bar","baz"]}}}`)

		list = types.NewList(types.Bool(true), types.Bool(false), types.Bool(false))

		suite.assertQueryResult(list, "{root{"+valuesKey+"}}", `{"data":{"root":{"`+valuesKey+`":[true,false,false]}}}`)
		suite.assertQueryResult(list, "{root{"+valuesKey+"(at:1,count:2)}}", `{"data":{"root":{"`+valuesKey+`":[false,false]}}}`)

		list = types.NewList(types.Number(1), types.Number(1.1), types.Number(-100))

		suite.assertQueryResult(list, "{root{"+valuesKey+"}}", `{"data":{"root":{"`+valuesKey+`":[1,1.1,-100]}}}`)
		suite.assertQueryResult(list, "{root{"+valuesKey+"(at:1,count:2)}}", `{"data":{"root":{"`+valuesKey+`":[1.1,-100]}}}`)

		list = types.NewList(types.String("a"), types.String("b"), types.String("c"))
		suite.assertQueryResult(list, "{root{"+valuesKey+"(at:4)}}", `{"data":{"root":{"`+valuesKey+`":[]}}}`)
		suite.assertQueryResult(list, "{root{"+valuesKey+"(count:0)}}", `{"data":{"root":{"`+valuesKey+`":[]}}}`)
		suite.assertQueryResult(list, "{root{"+valuesKey+"(count:10)}}", `{"data":{"root":{"`+valuesKey+`":["a","b","c"]}}}`)
		suite.assertQueryResult(list, "{root{"+valuesKey+"(at:-1)}}", `{"data":{"root":{"`+valuesKey+`":["a","b","c"]}}}`)
	}
}

func (suite *QueryGraphQLSuite) TestListOfStruct() {
	list := types.NewList(
		types.NewStruct("Foo", types.StructData{
			"a": types.Number(28),
			"b": types.String("foo"),
		}),
		types.NewStruct("Foo", types.StructData{
			"a": types.Number(-20.102),
			"b": types.String("bar"),
		}),
		types.NewStruct("Foo", types.StructData{
			"a": types.Number(5),
			"b": types.String("baz"),
		}),
	)

	suite.assertQueryResult(list, "{root{elements{a b}}}", `{"data":{"root":{"elements":[{"a":28,"b":"foo"},{"a":-20.102,"b":"bar"},{"a":5,"b":"baz"}]}}}`)

	suite.assertQueryResult(list, "{root{elements{a}}}", `{"data":{"root":{"elements":[{"a":28},{"a":-20.102},{"a":5}]}}}`)
}

func (suite *QueryGraphQLSuite) TestSetBasic() {
	for _, valuesKey := range []string{"elements", "values"} {
		set := types.NewSet()
		suite.assertQueryResult(set, "{root{size}}", `{"data":{"root":{"size":0}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"}}", `{"data":{"root":{}}}`)

		set = types.NewSet(types.String("foo"), types.String("bar"), types.String("baz"))

		suite.assertQueryResult(set, "{root{"+valuesKey+"}}", `{"data":{"root":{"`+valuesKey+`":["bar","baz","foo"]}}}`)
		suite.assertQueryResult(set, "{root{size}}", `{"data":{"root":{"size":3}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"(count:2)}}", `{"data":{"root":{"`+valuesKey+`":["bar","baz"]}}}`)

		set = types.NewSet(types.Bool(true), types.Bool(false))

		suite.assertQueryResult(set, "{root{"+valuesKey+"}}", `{"data":{"root":{"`+valuesKey+`":[false,true]}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"(count:1)}}", `{"data":{"root":{"`+valuesKey+`":[false]}}}`)

		set = types.NewSet(types.Number(1), types.Number(1.1), types.Number(-100))

		suite.assertQueryResult(set, "{root{"+valuesKey+"}}", `{"data":{"root":{"`+valuesKey+`":[-100,1,1.1]}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"(count:2)}}", `{"data":{"root":{"`+valuesKey+`":[-100,1]}}}`)

		set = types.NewSet(types.String("a"), types.String("b"), types.String("c"), types.String("d"))
		suite.assertQueryResult(set, "{root{"+valuesKey+"(count:0)}}", `{"data":{"root":{"`+valuesKey+`":[]}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"(count:2)}}", `{"data":{"root":{"`+valuesKey+`":["a","b"]}}}`)

		suite.assertQueryResult(set, "{root{"+valuesKey+"(at:0,count:2)}}", `{"data":{"root":{"`+valuesKey+`":["a","b"]}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"(at:-1,count:2)}}", `{"data":{"root":{"`+valuesKey+`":["a","b"]}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"(at:1,count:2)}}", `{"data":{"root":{"`+valuesKey+`":["b","c"]}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"(at:2)}}", `{"data":{"root":{"`+valuesKey+`":["c","d"]}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"(at:2,count:1)}}", `{"data":{"root":{"`+valuesKey+`":["c"]}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"(at:2,count:0)}}", `{"data":{"root":{"`+valuesKey+`":[]}}}`)
		suite.assertQueryResult(set, "{root{"+valuesKey+"(at:2,count:10)}}", `{"data":{"root":{"`+valuesKey+`":["c","d"]}}}`)
	}
}

func (suite *QueryGraphQLSuite) TestSetOfStruct() {
	set := types.NewSet(
		types.NewStruct("Foo", types.StructData{
			"a": types.Number(28),
			"b": types.String("foo"),
		}),
		types.NewStruct("Foo", types.StructData{
			"a": types.Number(-20.102),
			"b": types.String("bar"),
		}),
		types.NewStruct("Foo", types.StructData{
			"a": types.Number(5),
			"b": types.String("baz"),
		}),
	)

	suite.assertQueryResult(set, "{root{values{a b}}}", `{"data":{"root":{"values":[{"a":-20.102,"b":"bar"},{"a":5,"b":"baz"},{"a":28,"b":"foo"}]}}}`)
	suite.assertQueryResult(set, "{root{values{a}}}", `{"data":{"root":{"values":[{"a":-20.102},{"a":5},{"a":28}]}}}`)
}

func (suite *QueryGraphQLSuite) TestMapBasic() {
	for _, entriesKey := range []string{"elements", "entries"} {

		m := types.NewMap()
		suite.assertQueryResult(m, "{root{size}}", `{"data":{"root":{"size":0}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"}}", `{"data":{"root":{}}}`)

		m = types.NewMap(
			types.String("a"), types.Number(1),
			types.String("b"), types.Number(2),
			types.String("c"), types.Number(3),
			types.String("d"), types.Number(4),
		)

		suite.assertQueryResult(m, "{root{"+entriesKey+"{key value}}}", `{"data":{"root":{"`+entriesKey+`":[{"key":"a","value":1},{"key":"b","value":2},{"key":"c","value":3},{"key":"d","value":4}]}}}`)
		suite.assertQueryResult(m, "{root{size}}", `{"data":{"root":{"size":4}}}`)
	}
}

func (suite *QueryGraphQLSuite) TestMapOfStruct() {
	m := types.NewMap(
		types.String("foo"), types.NewStruct("Foo", types.StructData{
			"a": types.Number(28),
			"b": types.String("foo"),
		}),
		types.String("bar"), types.NewStruct("Foo", types.StructData{
			"a": types.Number(-20.102),
			"b": types.String("bar"),
		}),
		types.String("baz"), types.NewStruct("Foo", types.StructData{
			"a": types.Number(5),
			"b": types.String("baz"),
		}),
	)

	suite.assertQueryResult(m, "{root{entries{key value{a}}}}", `{"data":{"root":{"entries":[{"key":"bar","value":{"a":-20.102}},{"key":"baz","value":{"a":5}},{"key":"foo","value":{"a":28}}]}}}`)
	suite.assertQueryResult(m, "{root{entries(count:1){value{a b}}}}", `{"data":{"root":{"entries":[{"value":{"a":-20.102,"b":"bar"}}]}}}`)
	suite.assertQueryResult(m, "{root{entries(count:3){key}}}", `{"data":{"root":{"entries":[{"key":"bar"},{"key":"baz"},{"key":"foo"}]}}}`)
}

func (suite *QueryGraphQLSuite) TestRef() {
	r := suite.vs.WriteValue(types.Number(100))

	suite.assertQueryResult(r, "{root{targetValue}}", `{"data":{"root":{"targetValue":100}}}`)
	suite.assertQueryResult(r, "{root{targetHash}}", `{"data":{"root":{"targetHash":"fpbhln9asjlalp10btna9ocuc4nj9v15"}}}`)
	suite.assertQueryResult(r, "{root{targetValue targetHash}}", `{"data":{"root":{"targetHash":"fpbhln9asjlalp10btna9ocuc4nj9v15","targetValue":100}}}`)

	r = suite.vs.WriteValue(types.NewStruct("Foo", types.StructData{
		"a": types.Number(28),
		"b": types.String("foo"),
	}))

	suite.assertQueryResult(r, "{root{targetValue{a}}}", `{"data":{"root":{"targetValue":{"a":28}}}}`)
	suite.assertQueryResult(r, "{root{targetValue{a b}}}", `{"data":{"root":{"targetValue":{"a":28,"b":"foo"}}}}`)

	r = suite.vs.WriteValue(types.NewList(types.String("foo"), types.String("bar"), types.String("baz")))

	suite.assertQueryResult(r, "{root{targetValue{values}}}", `{"data":{"root":{"targetValue":{"values":["foo","bar","baz"]}}}}`)
	suite.assertQueryResult(r, "{root{targetValue{values(at:1,count:2)}}}", `{"data":{"root":{"targetValue":{"values":["bar","baz"]}}}}`)
}

func (suite *QueryGraphQLSuite) TestListOfUnionOfStructs() {
	list := types.NewList(
		types.NewStruct("Foo", types.StructData{
			"a": types.Number(28),
			"b": types.String("baz"),
		}),
		types.NewStruct("Bar", types.StructData{
			"b": types.String("bar"),
		}),
		types.NewStruct("Baz", types.StructData{
			"c": types.Bool(true),
		}),
	)

	suite.assertQueryResult(list,
		fmt.Sprintf("{root{values{... on %s{a b} ... on %s{b} ... on %s{c}}}}",
			getTypeName(list.Get(0).Type()),
			getTypeName(list.Get(1).Type()),
			getTypeName(list.Get(2).Type())),
		`{"data":{"root":{"values":[{"a":28,"b":"baz"},{"b":"bar"},{"c":true}]}}}`)
}

func (suite *QueryGraphQLSuite) TestListOfUnionOfStructsConflictingFieldTypes() {
	list := types.NewList(
		types.NewStruct("Foo", types.StructData{
			"a": types.Number(28),
		}),
		types.NewStruct("Bar", types.StructData{
			"a": types.String("bar"),
		}),
		types.NewStruct("Baz", types.StructData{
			"a": types.Bool(true),
		}),
	)

	suite.assertQueryResult(list,
		fmt.Sprintf("{root{values{... on %s{a} ... on %s{b: a} ... on %s{c: a}}}}",
			getTypeName(list.Get(0).Type()),
			getTypeName(list.Get(1).Type()),
			getTypeName(list.Get(2).Type())),
		`{"data":{"root":{"values":[{"a":28},{"b":"bar"},{"c":true}]}}}`)
}

func (suite *QueryGraphQLSuite) TestListOfUnionOfScalars() {
	list := types.NewList(
		types.Number(28),
		types.String("bar"),
		types.Bool(true),
	)

	suite.assertQueryResult(list, "{root{values{... on BooleanValue{b: scalarValue} ... on StringValue{s: scalarValue} ... on NumberValue{n: scalarValue}}}}", `{"data":{"root":{"values":[{"n":28},{"s":"bar"},{"b":true}]}}}`)
}

func (suite *QueryGraphQLSuite) TestCyclicStructs() {
	typ := types.MakeStructTypeFromFields("A", types.FieldMap{
		"a": types.StringType,
		"b": types.MakeSetType(types.MakeCycleType(0)),
	})

	// Struct A {
	//  a: "aaa"
	//  b: Set(Struct A {
	// 	 a: "bbb"
	// 	 b: Set()
	//  })
	// }

	s1 := types.NewStructWithType(typ, types.ValueSlice{
		types.String("aaa"),
		types.NewSet(types.NewStructWithType(typ, types.ValueSlice{types.String("bbb"), types.NewSet()})),
	})

	suite.assertQueryResult(s1, "{root{a b{values{a}}}}", `{"data":{"root":{"a":"aaa","b":{"values":[{"a":"bbb"}]}}}}`)
}

func (suite *QueryGraphQLSuite) TestCyclicStructsWithUnion() {
	typ := types.MakeStructTypeFromFields("A", types.FieldMap{
		"a": types.StringType,
		"b": types.MakeUnionType(types.NumberType, types.MakeCycleType(0)),
	})

	// Struct A {
	//  a: "aaa"
	//  b: Struct A {
	// 	 a: "bbb"
	// 	 b: 42
	//  })
	// }

	s1 := types.NewStructWithType(typ, types.ValueSlice{
		types.String("aaa"),
		types.NewStructWithType(typ, types.ValueSlice{types.String("bbb"), types.Number(42)}),
	})

	suite.assertQueryResult(s1,
		fmt.Sprintf(`{root{a b {... on %s{a}}}}`, getTypeName(s1.Type())),
		`{"data":{"root":{"a":"aaa","b":{"a":"bbb"}}}}`)
}

func (suite *QueryGraphQLSuite) TestNestedCollection() {
	list := types.NewList(
		types.NewSet(
			types.NewMap(types.Number(10), types.String("foo")),
			types.NewMap(types.Number(20), types.String("bar")),
		),
		types.NewSet(
			types.NewMap(types.Number(30), types.String("baz")),
			types.NewMap(types.Number(40), types.String("bat")),
		),
	)

	suite.assertQueryResult(list, "{root{size}}", `{"data":{"root":{"size":2}}}`)
	suite.assertQueryResult(list, "{root{values(count:1){size}}}", `{"data":{"root":{"values":[{"size":2}]}}}`)
	suite.assertQueryResult(list, "{root{values(at:1,count:1){values(count:1){entries{key value}}}}}", `{"data":{"root":{"values":[{"values":[{"entries":[{"key":30,"value":"baz"}]}]}]}}}`)
}

func (suite *QueryGraphQLSuite) TestLoFi() {
	b := types.NewBlob(bytes.NewBufferString("I am a blob"))

	suite.assertQueryResult(b, "{root}", `{"data":{"root":"h6jkv35uum62a7ovu14uvmhaf0sojgh6"}}`)

	t := types.StringType
	suite.assertQueryResult(t, "{root}", `{"data":{"root":"pej65tf21rubhu9cb0oi5gqrkgf26aql"}}`)
}

func (suite *QueryGraphQLSuite) TestError() {
	buff := &bytes.Buffer{}
	Error(errors.New("Some error string"), buff)
	suite.Equal(buff.String(), `{"data":null,"errors":[{"message":"Some error string","locations":null}]}
`)
}

func (suite *QueryGraphQLSuite) TestMapArgs() {
	for _, entriesKey := range []string{"elements", "entries"} {

		m := types.NewMap(
			types.String("a"), types.Number(1),
			types.String("c"), types.Number(2),
			types.String("e"), types.Number(3),
			types.String("g"), types.Number(4),
		)

		// count
		suite.assertQueryResult(m, "{root{"+entriesKey+"(count:0){value}}}", `{"data":{"root":{"`+entriesKey+`":[]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(count:2){value}}}", `{"data":{"root":{"`+entriesKey+`":[{"value":1},{"value":2}]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(count:3){key}}}", `{"data":{"root":{"`+entriesKey+`":[{"key":"a"},{"key":"c"},{"key":"e"}]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(count: -1){key}}}", `{"data":{"root":{"`+entriesKey+`":[]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(count:5){value}}}", `{"data":{"root":{"`+entriesKey+`":[{"value":1},{"value":2},{"value":3},{"value":4}]}}}`)

		// at
		suite.assertQueryResult(m, "{root{"+entriesKey+"(at:0){value}}}", `{"data":{"root":{"`+entriesKey+`":[{"value":1},{"value":2},{"value":3},{"value":4}]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(at:-1){value}}}", `{"data":{"root":{"`+entriesKey+`":[{"value":1},{"value":2},{"value":3},{"value":4}]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(at:2){value}}}", `{"data":{"root":{"`+entriesKey+`":[{"value":3},{"value":4}]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(at:5){value}}}", `{"data":{"root":{"`+entriesKey+`":[]}}}`)

		// at & count
		suite.assertQueryResult(m, "{root{"+entriesKey+"(at:0,count:2){value}}}", `{"data":{"root":{"`+entriesKey+`":[{"value":1},{"value":2}]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(at:-1,count:2){value}}}", `{"data":{"root":{"`+entriesKey+`":[{"value":1},{"value":2}]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(at:1,count:2){value}}}", `{"data":{"root":{"`+entriesKey+`":[{"value":2},{"value":3}]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(at:2,count:1){value}}}", `{"data":{"root":{"`+entriesKey+`":[{"value":3}]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(at:2,count:0){value}}}", `{"data":{"root":{"`+entriesKey+`":[]}}}`)
		suite.assertQueryResult(m, "{root{"+entriesKey+"(at:2,count:10){value}}}", `{"data":{"root":{"`+entriesKey+`":[{"value":3},{"value":4}]}}}`)

		// key
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"e"){key value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"e","value":3}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"g"){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"value":4}]}}}`)
		// "f", no count/through so asking for exact match
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"f"){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[]}}}`)
		// "x" is larger than end
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"x"){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[]}}}`)

		// key & at
		// at is ignored when key is present
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"e",at:2){key value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"e","value":3}]}}}`)

		// key & count
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"c", count: 2){key value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"c","value":2},{"key":"e","value":3}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"c", count: 0){key value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"c", count: -1){key value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"e", count: 5){key value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"e","value":3},{"key":"g","value":4}]}}}`)

		// through
		suite.assertQueryResult(m, `{root{`+entriesKey+`(through:"c"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"a"},{"key":"c"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(through:"b"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"a"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(through:"0"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[]}}}`)

		// key & through
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"c", through:"c"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"c"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"c",through:"e"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"c"},{"key":"e"}]}}}`)

		// through & count
		suite.assertQueryResult(m, `{root{`+entriesKey+`(through:"c",count:1){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"a"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(through:"b",count:0){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(through:"0",count:10){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[]}}}`)

		// at & through
		suite.assertQueryResult(m, `{root{`+entriesKey+`(at:0,through:"a"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"a"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(at:1,through:"e"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"c"},{"key":"e"}]}}}`)

		// at & count & through
		suite.assertQueryResult(m, `{root{`+entriesKey+`(at:0,count:2,through:"a"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"a"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(at:0,count:2,through:"e"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"a"},{"key":"c"}]}}}`)

		// key & count & through
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"c",count:2,through:"c"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"c"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(key:"c",count:2,through:"g"){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"c"},{"key":"e"}]}}}`)
	}
}

func (suite *QueryGraphQLSuite) TestMapKeysArg() {
	for _, entriesKey := range []string{"elements", "entries"} {
		m := types.NewMap(
			types.String("a"), types.Number(1),
			types.String("c"), types.Number(2),
			types.String("e"), types.Number(3),
			types.String("g"), types.Number(4),
		)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(keys:["c","a"]){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"value":2},{"value":1}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(keys:[]){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[]}}}`)

		m = types.NewMap(
			types.Number(1), types.String("a"),
			types.Number(2), types.String("c"),
			types.Number(3), types.String("e"),
			types.Number(4), types.String("g"),
		)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(keys:[4,1]){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"value":"g"},{"value":"a"}]}}}`)

		// Ignore other args
		suite.assertQueryResult(m, `{root{`+entriesKey+`(keys:[4,1],key:2){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"value":"g"},{"value":"a"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(keys:[4,1],count:0){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"value":"g"},{"value":"a"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(keys:[4,1],at:4){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"value":"g"},{"value":"a"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(keys:[4,1],through:1){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"value":"g"},{"value":"a"}]}}}`)
	}
}

func (suite *QueryGraphQLSuite) TestSetArgs() {
	for _, valuesKey := range []string{"elements", "values"} {
		s := types.NewSet(
			types.String("a"),
			types.String("c"),
			types.String("e"),
			types.String("g"),
		)

		// count
		suite.assertQueryResult(s, "{root{"+valuesKey+"(count:0)}}", `{"data":{"root":{"`+valuesKey+`":[]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(count:2)}}", `{"data":{"root":{"`+valuesKey+`":["a","c"]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(count:3)}}", `{"data":{"root":{"`+valuesKey+`":["a","c","e"]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(count: -1)}}", `{"data":{"root":{"`+valuesKey+`":[]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(count:5)}}", `{"data":{"root":{"`+valuesKey+`":["a","c","e","g"]}}}`)

		// at
		suite.assertQueryResult(s, "{root{"+valuesKey+"(at:0)}}", `{"data":{"root":{"`+valuesKey+`":["a","c","e","g"]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(at:-1)}}", `{"data":{"root":{"`+valuesKey+`":["a","c","e","g"]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(at:2)}}", `{"data":{"root":{"`+valuesKey+`":["e","g"]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(at:5)}}", `{"data":{"root":{"`+valuesKey+`":[]}}}`)

		// at & count
		suite.assertQueryResult(s, "{root{"+valuesKey+"(at:0,count:2)}}", `{"data":{"root":{"`+valuesKey+`":["a","c"]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(at:-1,count:2)}}", `{"data":{"root":{"`+valuesKey+`":["a","c"]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(at:1,count:2)}}", `{"data":{"root":{"`+valuesKey+`":["c","e"]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(at:2,count:1)}}", `{"data":{"root":{"`+valuesKey+`":["e"]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(at:2,count:0)}}", `{"data":{"root":{"`+valuesKey+`":[]}}}`)
		suite.assertQueryResult(s, "{root{"+valuesKey+"(at:2,count:10)}}", `{"data":{"root":{"`+valuesKey+`":["e","g"]}}}`)

		// key
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"e")}}`,
			`{"data":{"root":{"`+valuesKey+`":["e"]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"g")}}`,
			`{"data":{"root":{"`+valuesKey+`":["g"]}}}`)
		// "f", no count/through so asking for exact match
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"f")}}`,
			`{"data":{"root":{"`+valuesKey+`":[]}}}`)
		// "x" is larger than end
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"x")}}`,
			`{"data":{"root":{"`+valuesKey+`":[]}}}`)
		// exact match
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"0")}}`,
			`{"data":{"root":{"`+valuesKey+`":[]}}}`)

		// key & at
		// at is ignored when key is present
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"e",at:2)}}`,
			`{"data":{"root":{"`+valuesKey+`":["e"]}}}`)

		// key & count
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"c", count: 2)}}`,
			`{"data":{"root":{"`+valuesKey+`":["c","e"]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"c", count: 0)}}`,
			`{"data":{"root":{"`+valuesKey+`":[]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"c", count: -1)}}`,
			`{"data":{"root":{"`+valuesKey+`":[]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"e", count: 5)}}`,
			`{"data":{"root":{"`+valuesKey+`":["e","g"]}}}`)

		// through
		suite.assertQueryResult(s, `{root{`+valuesKey+`(through:"c")}}`,
			`{"data":{"root":{"`+valuesKey+`":["a","c"]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(through:"b")}}`,
			`{"data":{"root":{"`+valuesKey+`":["a"]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(through:"0")}}`,
			`{"data":{"root":{"`+valuesKey+`":[]}}}`)

		// key & through
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"c", through:"c")}}`,
			`{"data":{"root":{"`+valuesKey+`":["c"]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"c",through:"e")}}`,
			`{"data":{"root":{"`+valuesKey+`":["c","e"]}}}`)

		// through & count
		suite.assertQueryResult(s, `{root{`+valuesKey+`(through:"c",count:1)}}`,
			`{"data":{"root":{"`+valuesKey+`":["a"]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(through:"b",count:0)}}`,
			`{"data":{"root":{"`+valuesKey+`":[]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(through:"0",count:10)}}`,
			`{"data":{"root":{"`+valuesKey+`":[]}}}`)

		// at & through
		suite.assertQueryResult(s, `{root{`+valuesKey+`(at:0,through:"a")}}`,
			`{"data":{"root":{"`+valuesKey+`":["a"]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(at:1,through:"e")}}`,
			`{"data":{"root":{"`+valuesKey+`":["c","e"]}}}`)

		// at & count & through
		suite.assertQueryResult(s, `{root{`+valuesKey+`(at:0,count:2,through:"a")}}`,
			`{"data":{"root":{"`+valuesKey+`":["a"]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(at:0,count:2,through:"e")}}`,
			`{"data":{"root":{"`+valuesKey+`":["a","c"]}}}`)

		// key & count & through
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"c",count:2,through:"c")}}`,
			`{"data":{"root":{"`+valuesKey+`":["c"]}}}`)
		suite.assertQueryResult(s, `{root{`+valuesKey+`(key:"c",count:2,through:"g")}}`,
			`{"data":{"root":{"`+valuesKey+`":["c","e"]}}}`)
	}
}

func (suite *QueryGraphQLSuite) TestMapValues() {
	m := types.NewMap(
		types.String("a"), types.Number(1),
		types.String("c"), types.Number(2),
		types.String("e"), types.Number(3),
		types.String("g"), types.Number(4),
	)

	suite.assertQueryResult(m, "{root{values}}", `{"data":{"root":{"values":[1,2,3,4]}}}`)

	// count
	suite.assertQueryResult(m, "{root{values(count:0)}}", `{"data":{"root":{"values":[]}}}`)
	suite.assertQueryResult(m, "{root{values(count:2)}}", `{"data":{"root":{"values":[1,2]}}}`)
	suite.assertQueryResult(m, "{root{values(count:3)}}", `{"data":{"root":{"values":[1,2,3]}}}`)
	suite.assertQueryResult(m, "{root{values(count: -1)}}", `{"data":{"root":{"values":[]}}}`)
	suite.assertQueryResult(m, "{root{values(count:5)}}", `{"data":{"root":{"values":[1,2,3,4]}}}`)

	// at
	suite.assertQueryResult(m, "{root{values(at:0)}}", `{"data":{"root":{"values":[1,2,3,4]}}}`)
	suite.assertQueryResult(m, "{root{values(at:-1)}}", `{"data":{"root":{"values":[1,2,3,4]}}}`)
	suite.assertQueryResult(m, "{root{values(at:2)}}", `{"data":{"root":{"values":[3,4]}}}`)
	suite.assertQueryResult(m, "{root{values(at:5)}}", `{"data":{"root":{"values":[]}}}`)

	// at & count
	suite.assertQueryResult(m, "{root{values(at:0,count:2)}}", `{"data":{"root":{"values":[1,2]}}}`)
	suite.assertQueryResult(m, "{root{values(at:-1,count:2)}}", `{"data":{"root":{"values":[1,2]}}}`)
	suite.assertQueryResult(m, "{root{values(at:1,count:2)}}", `{"data":{"root":{"values":[2,3]}}}`)
	suite.assertQueryResult(m, "{root{values(at:2,count:1)}}", `{"data":{"root":{"values":[3]}}}`)
	suite.assertQueryResult(m, "{root{values(at:2,count:0)}}", `{"data":{"root":{"values":[]}}}`)
	suite.assertQueryResult(m, "{root{values(at:2,count:10)}}", `{"data":{"root":{"values":[3,4]}}}`)

	// key
	suite.assertQueryResult(m, `{root{values(key:"e")}}`, `{"data":{"root":{"values":[3]}}}`)
	suite.assertQueryResult(m, `{root{values(key:"g")}}`, `{"data":{"root":{"values":[4]}}}`)
	// "f", no count/through so asking for exact match
	suite.assertQueryResult(m, `{root{values(key:"f")}}`, `{"data":{"root":{"values":[]}}}`)
	// "x" is larger than end
	suite.assertQueryResult(m, `{root{values(key:"x")}}`, `{"data":{"root":{"values":[]}}}`)

	// key & at
	// at is ignored when key is present
	suite.assertQueryResult(m, `{root{values(key:"e",at:2)}}`, `{"data":{"root":{"values":[3]}}}`)

	// key & count
	suite.assertQueryResult(m, `{root{values(key:"c",count:2)}}`, `{"data":{"root":{"values":[2,3]}}}`)
	suite.assertQueryResult(m, `{root{values(key:"c",count:0)}}`, `{"data":{"root":{"values":[]}}}`)
	suite.assertQueryResult(m, `{root{values(key:"c",count:-1)}}`, `{"data":{"root":{"values":[]}}}`)
	suite.assertQueryResult(m, `{root{values(key:"e",count:5)}}`, `{"data":{"root":{"values":[3,4]}}}`)

	// through
	suite.assertQueryResult(m, `{root{values(through:"c")}}`, `{"data":{"root":{"values":[1,2]}}}`)
	suite.assertQueryResult(m, `{root{values(through:"b")}}`, `{"data":{"root":{"values":[1]}}}`)
	suite.assertQueryResult(m, `{root{values(through:"0")}}`, `{"data":{"root":{"values":[]}}}`)

	// key & through
	suite.assertQueryResult(m, `{root{values(key:"c", through:"c")}}`, `{"data":{"root":{"values":[2]}}}`)
	suite.assertQueryResult(m, `{root{values(key:"c",through:"e")}}`, `{"data":{"root":{"values":[2,3]}}}`)

	// through & count
	suite.assertQueryResult(m, `{root{values(through:"c",count:1)}}`, `{"data":{"root":{"values":[1]}}}`)
	suite.assertQueryResult(m, `{root{values(through:"b",count:0)}}`, `{"data":{"root":{"values":[]}}}`)
	suite.assertQueryResult(m, `{root{values(through:"0",count:10)}}`, `{"data":{"root":{"values":[]}}}`)

	// at & through
	suite.assertQueryResult(m, `{root{values(at:0,through:"a")}}`, `{"data":{"root":{"values":[1]}}}`)
	suite.assertQueryResult(m, `{root{values(at:1,through:"e")}}`, `{"data":{"root":{"values":[2,3]}}}`)

	// at & count & through
	suite.assertQueryResult(m, `{root{values(at:0,count:2,through:"a")}}`, `{"data":{"root":{"values":[1]}}}`)
	suite.assertQueryResult(m, `{root{values(at:0,count:2,through:"e")}}`, `{"data":{"root":{"values":[1,2]}}}`)

	// key & count & through
	suite.assertQueryResult(m, `{root{values(key:"c",count:2,through:"c")}}`,
		`{"data":{"root":{"values":[2]}}}`)
	suite.assertQueryResult(m, `{root{values(key:"c",count:2,through:"g")}}`,
		`{"data":{"root":{"values":[2,3]}}}`)
}

func (suite *QueryGraphQLSuite) TestMapKeys() {
	m := types.NewMap(
		types.String("a"), types.Number(1),
		types.String("c"), types.Number(2),
		types.String("e"), types.Number(3),
		types.String("g"), types.Number(4),
	)

	suite.assertQueryResult(m, "{root{keys}}", `{"data":{"root":{"keys":["a","c","e","g"]}}}`)

	// count
	suite.assertQueryResult(m, "{root{keys(count:0)}}", `{"data":{"root":{"keys":[]}}}`)
	suite.assertQueryResult(m, "{root{keys(count:2)}}", `{"data":{"root":{"keys":["a","c"]}}}`)
	suite.assertQueryResult(m, "{root{keys(count:3)}}", `{"data":{"root":{"keys":["a","c","e"]}}}`)
	suite.assertQueryResult(m, "{root{keys(count: -1)}}", `{"data":{"root":{"keys":[]}}}`)
	suite.assertQueryResult(m, "{root{keys(count:5)}}", `{"data":{"root":{"keys":["a","c","e","g"]}}}`)

	// at
	suite.assertQueryResult(m, "{root{keys(at:0)}}", `{"data":{"root":{"keys":["a","c","e","g"]}}}`)
	suite.assertQueryResult(m, "{root{keys(at:-1)}}", `{"data":{"root":{"keys":["a","c","e","g"]}}}`)
	suite.assertQueryResult(m, "{root{keys(at:2)}}", `{"data":{"root":{"keys":["e","g"]}}}`)
	suite.assertQueryResult(m, "{root{keys(at:5)}}", `{"data":{"root":{"keys":[]}}}`)

	// at & count
	suite.assertQueryResult(m, "{root{keys(at:0,count:2)}}", `{"data":{"root":{"keys":["a","c"]}}}`)
	suite.assertQueryResult(m, "{root{keys(at:-1,count:2)}}", `{"data":{"root":{"keys":["a","c"]}}}`)
	suite.assertQueryResult(m, "{root{keys(at:1,count:2)}}", `{"data":{"root":{"keys":["c","e"]}}}`)
	suite.assertQueryResult(m, "{root{keys(at:2,count:1)}}", `{"data":{"root":{"keys":["e"]}}}`)
	suite.assertQueryResult(m, "{root{keys(at:2,count:0)}}", `{"data":{"root":{"keys":[]}}}`)
	suite.assertQueryResult(m, "{root{keys(at:2,count:10)}}", `{"data":{"root":{"keys":["e","g"]}}}`)

	// key
	suite.assertQueryResult(m, `{root{keys(key:"e")}}`, `{"data":{"root":{"keys":["e"]}}}`)
	suite.assertQueryResult(m, `{root{keys(key:"g")}}`, `{"data":{"root":{"keys":["g"]}}}`)
	// "f", no count/through so asking for exact match
	suite.assertQueryResult(m, `{root{keys(key:"f")}}`, `{"data":{"root":{"keys":[]}}}`)
	// "x" is larger than end
	suite.assertQueryResult(m, `{root{keys(key:"x")}}`, `{"data":{"root":{"keys":[]}}}`)

	// key & at
	// at is ignored when key is present
	suite.assertQueryResult(m, `{root{keys(key:"e",at:2)}}`, `{"data":{"root":{"keys":["e"]}}}`)

	// key & count
	suite.assertQueryResult(m, `{root{keys(key:"c",count:2)}}`, `{"data":{"root":{"keys":["c","e"]}}}`)
	suite.assertQueryResult(m, `{root{keys(key:"c",count:0)}}`, `{"data":{"root":{"keys":[]}}}`)
	suite.assertQueryResult(m, `{root{keys(key:"c",count:-1)}}`, `{"data":{"root":{"keys":[]}}}`)
	suite.assertQueryResult(m, `{root{keys(key:"e",count:5)}}`, `{"data":{"root":{"keys":["e","g"]}}}`)

	// through
	suite.assertQueryResult(m, `{root{keys(through:"c")}}`, `{"data":{"root":{"keys":["a","c"]}}}`)
	suite.assertQueryResult(m, `{root{keys(through:"b")}}`, `{"data":{"root":{"keys":["a"]}}}`)
	suite.assertQueryResult(m, `{root{keys(through:"0")}}`, `{"data":{"root":{"keys":[]}}}`)

	// key & through
	suite.assertQueryResult(m, `{root{keys(key:"c", through:"c")}}`, `{"data":{"root":{"keys":["c"]}}}`)
	suite.assertQueryResult(m, `{root{keys(key:"c",through:"e")}}`, `{"data":{"root":{"keys":["c","e"]}}}`)

	// through & count
	suite.assertQueryResult(m, `{root{keys(through:"c",count:1)}}`, `{"data":{"root":{"keys":["a"]}}}`)
	suite.assertQueryResult(m, `{root{keys(through:"b",count:0)}}`, `{"data":{"root":{"keys":[]}}}`)
	suite.assertQueryResult(m, `{root{keys(through:"0",count:10)}}`, `{"data":{"root":{"keys":[]}}}`)

	// at & through
	suite.assertQueryResult(m, `{root{keys(at:0,through:"a")}}`, `{"data":{"root":{"keys":["a"]}}}`)
	suite.assertQueryResult(m, `{root{keys(at:1,through:"e")}}`, `{"data":{"root":{"keys":["c","e"]}}}`)

	// at & count & through
	suite.assertQueryResult(m, `{root{keys(at:0,count:2,through:"a")}}`, `{"data":{"root":{"keys":["a"]}}}`)
	suite.assertQueryResult(m, `{root{keys(at:0,count:2,through:"e")}}`, `{"data":{"root":{"keys":["a","c"]}}}`)

	// key & count & through
	suite.assertQueryResult(m, `{root{keys(key:"c",count:2,through:"c")}}`,
		`{"data":{"root":{"keys":["c"]}}}`)
	suite.assertQueryResult(m, `{root{keys(key:"c",count:2,through:"g")}}`,
		`{"data":{"root":{"keys":["c","e"]}}}`)
}

func (suite *QueryGraphQLSuite) TestMapNullable() {
	// When selecting the result based on keys the values may be null.
	m := types.NewMap(
		types.String("a"), types.Number(1),
		types.String("c"), types.Number(2),
	)

	for _, entriesKey := range []string{"elements", "entries"} {
		suite.assertQueryResult(m, `{root{`+entriesKey+`(keys:["a","b","c"]){value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"value":1},{"value":null},{"value":2}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(keys:["a","b","c"]){key}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"a"},{"key":"b"},{"key":"c"}]}}}`)
		suite.assertQueryResult(m, `{root{`+entriesKey+`(keys:["a","b","c"]){key value}}}`,
			`{"data":{"root":{"`+entriesKey+`":[{"key":"a","value":1},{"key":"b","value":null},{"key":"c","value":2}]}}}`)
	}
	suite.assertQueryResult(m, `{root{values(keys:["a","b","c"])}}`,
		`{"data":{"root":{"values":[1,null,2]}}}`)
	suite.assertQueryResult(m, `{root{keys(keys:["a","b","c"])}}`,
		`{"data":{"root":{"keys":["a","b","c"]}}}`)
}
