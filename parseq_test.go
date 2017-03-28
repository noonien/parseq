package parseq_test

import (
	"net/url"
	"testing"

	"github.com/noonien/parseq"
	"github.com/stretchr/testify/require"
)

type t1 struct {
	Str  string
	Int  int      `query:"int"`
	Bool bool     `json:"bool"`
	Arr  []string `query:"arr"`
	Ign  string   `parseq:"-"`
	Ign2 string   `json:"-"`
}

type t2 struct {
	t1
}

func TestNames(t *testing.T) {
	tests := []struct {
		query string
		in    interface{}
		out   interface{}
	}{
		{"Str=str", &t1{}, &t1{Str: "str"}},
		{"int=42", &t1{}, &t1{Int: 42}},
		{"bool=true", &t1{}, &t1{Bool: true}},
		{"arr=i1&arr=i2", &t1{}, &t1{Arr: []string{"i1", "i2"}}},
		{"ign=str", &t1{}, &t1{}},
		{"ign2=str", &t1{}, &t1{}},
		{"int=42", &t2{}, &t2{t1{Int: 42}}},
	}

	require := require.New(t)
	for _, test := range tests {
		q, err := url.ParseQuery(test.query)
		require.Nil(err, "there should be no errors while parsing \"%s\"", test.query)

		err = parseq.Unmarshal(q, test.in)
		require.Nil(err, "there should be no errors while unmarshaling \"%s\" into %T", test.query, test.in)

		require.Equal(test.out, test.in)
	}

}
