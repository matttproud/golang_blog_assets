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

// go test -benchmem -bench .

package tree

import (
	"testing"
)

func benchmarkJIncorrect(b *testing.B, j int, t Traverser) {
	// Create our system under test (SUT):
	tree := New(j)
	// We create a bunch of recipient channels in advance for our traverser to
	// write into.
	chs := make(chan chan int, b.N)
	for i := 0; i < b.N; i++ {
		chs <- make(chan int, j)
	}
	close(chs)

	// Run our benchmark.
	for ch := range chs { // len(chs) == b.N
		t(tree, ch)
	}
}

func benchmarkJCorrect(b *testing.B, j int, t Traverser) {
	// IMPORTANT: Instruct the benchmarker to not track time and memory allocations!
	b.StopTimer()

	// Create our system under test (SUT):
	tree := New(j)
	// We create a bunch of recipient channels in advance for our traverser to
	// write into.
	chs := make(chan chan int, b.N)
	for i := 0; i < b.N; i++ {
		chs <- make(chan int, j)
	}
	close(chs)

	// IMPORTANT: Instruct the benchmarker to begin tracking time and memory
	// allocations.  Everything after this mark is instrumented!
	b.StartTimer()

	// Run our benchmark.
	for ch := range chs { // len(chs) == b.N
		t(tree, ch)
	}
}

// Choose your poison: benchmarkJIncorrect or benchmarkJCorrect.
var benchmarker = benchmarkJIncorrect

func BenchmarkRecursive1(b *testing.B) {
	benchmarker(b, 1, TraverseRecursive)
}

func BenchmarkRecursive10(b *testing.B) {
	benchmarker(b, 10, TraverseRecursive)
}

func BenchmarkRecursive100(b *testing.B) {
	benchmarker(b, 100, TraverseRecursive)
}

func BenchmarkRecursive1000(b *testing.B) {
	benchmarker(b, 1000, TraverseRecursive)
}

func BenchmarkIterative1(b *testing.B) {
	benchmarker(b, 1, TraverseIterative)
}

func BenchmarkIterative10(b *testing.B) {
	benchmarker(b, 10, TraverseIterative)
}

func BenchmarkIterative100(b *testing.B) {
	benchmarker(b, 100, TraverseIterative)
}

func BenchmarkIterative1000(b *testing.B) {
	benchmarker(b, 1000, TraverseIterative)
}
