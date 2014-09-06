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

// Increment consumes a slice of integers and returns a new slice that contains
// a copy of the original values but each value having been respectively
// incremented by one.
func Increment(in []int) []int {
	if in == nil {
		return nil
	}
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = v + 1
	}
	return out
}

// Hold your horses, and ignore sort.IntSlice for a moment.
type IntSlice []int

func (s IntSlice) Equal(o IntSlice) bool {
	if len(s) != len(o) {
		return false
	}
	for i, v := range s {
		if other := o[i]; other != v {
			return false
		}
	}
	return true
}
