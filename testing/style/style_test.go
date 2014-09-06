// Copyright 2014 Matt T. Proud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package tree provides examples of how to properly benchmark Go code that
// requires initial setup.
package style

import (
	"reflect"
	"testing"
)

func TestIterative(t *testing.T) {
	for i, test := range []struct {
		in, out []int
	}{
		{},
		{in: []int{1}, out: []int{2}},
		{in: []int{1, 2}, out: []int{2, 3}},
	} {
		out := Increment(test.in)
		if lactual, lexpected := len(out), len(test.out); lactual != lexpected {
			t.Fatalf("%d. got unexpected length %d instead of %d", i, lactual, lexpected)
		}
		for actualIdx, actualVal := range out {
			if expectedVal := test.out[actualIdx]; expectedVal != actualVal {
				t.Fatalf("%d.%d got unexpected value %d instead of %d", i, actualIdx, actualVal, expectedVal)
			}
		}
	}
}

func TestIntSliceIterative(t *testing.T) {
	for i, test := range []struct {
		in, out []int
	}{
		{},
		{in: []int{1}, out: []int{2}},
		{in: []int{1, 2}, out: []int{2, 3}},
	} {
		if out, testOut := IntSlice(Increment(test.in)), IntSlice(test.out); !testOut.Equal(out) {
			t.Fatalf("%d. got unexpected value %s instead of %s", i, out, test.out)
		}
	}
}

func TestReflect(t *testing.T) {
	for i, test := range []struct {
		in, out []int
	}{
		{},
		{in: []int{1}, out: []int{2}},
		{in: []int{1, 2}, out: []int{2, 3}},
	} {
		if out := Increment(test.in); !reflect.DeepEqual(test.out, out) {
			t.Fatalf("%d. got unexpected value %#v instead of %#v", i, out, test.out)
		}
	}
}

func TestReflectRenamed(t *testing.T) {
	for i, test := range []struct {
		in, want []int
	}{
		{},
		{in: []int{1}, want: []int{2}},
		{in: []int{1, 2}, want: []int{2, 3}},
	} {
		if got := Increment(test.in); !reflect.DeepEqual(got, test.want) {
			t.Fatalf("%d. got unexpected value %s instead of %s", i, got, test.want)
		}
	}
}

func TestReflectReordered(t *testing.T) {
	for i, test := range []struct {
		in, out []int
	}{
		{},
		{in: []int{1}, out: []int{2}},
		{in: []int{1, 2}, out: []int{2, 3}},
	} {
		if out := Increment(test.in); !reflect.DeepEqual(out, test.out) {
			t.Fatalf("%d. got unexpected value %s instead of %s", i, out, test.out)
		}
	}
}

func TestReflectPristine(t *testing.T) {
	for _, test := range []struct {
		in, want []int
	}{
		{},
		{in: []int{1}, want: []int{2}},
		{in: []int{1, 2}, want: []int{2, 3}},
	} {
		if got := Increment(test.in); !reflect.DeepEqual(got, test.want) {
			t.Fatalf("Increment(%#v) = %#v; want %#v", test.in, got, test.want)
		}
	}
}
