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

// go test -test.bench="Benchmark*" -benchmem fifo_test.go
package fifo

import (
	"container/list"
	"testing"
)

// Seen this a few places.  Wastes memory due to re-allocating slices.
func benchRealloc(b *testing.B, cap int) {
	b.StopTimer()
	buf := make([]int, 0, cap)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if l := len(buf); l < cap {
			buf = append(buf, i)
		} else {
			buf = append(buf[1:], i)
		}
	}
}

func BenchmarkRealloc1(b *testing.B) {
	benchRealloc(b, 1)
}

func BenchmarkRealloc10(b *testing.B) {
	benchRealloc(b, 10)
}

func BenchmarkRealloc100(b *testing.B) {
	benchRealloc(b, 100)
}

func BenchmarkRealloc1000(b *testing.B) {
	benchRealloc(b, 1000)
}

// Zero allocation but can be slower than benchRealloc depending on system.
func benchShift(b *testing.B, cap int) {
	b.StopTimer()
	buf := make([]int, 0, cap)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if l := len(buf); l < cap {
			buf = append(buf, i)
		} else {
			copy(buf, buf[1:])
			buf[cap-1] = i
		}
	}
}

func BenchmarkShift1(b *testing.B) {
	benchShift(b, 1)
}

func BenchmarkShift10(b *testing.B) {
	benchShift(b, 10)
}

func BenchmarkShift100(b *testing.B) {
	benchShift(b, 100)
}

func BenchmarkShift1000(b *testing.B) {
	benchShift(b, 1000)
}

// Allocate only on warmup, if at all.
func benchList(b *testing.B, cap int) {
	b.StopTimer()
	lst := list.New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if l := lst.Len(); l < cap {
			lst.PushBack(i)
		} else {
			frnt := lst.Front()
			frnt.Value = i
			lst.MoveToBack(frnt)
		}
	}
}

func BenchmarkList1(b *testing.B) {
	benchList(b, 1)
}

func BenchmarkList10(b *testing.B) {
	benchList(b, 10)
}

func BenchmarkList100(b *testing.B) {
	benchList(b, 100)
}

func BenchmarkList1000(b *testing.B) {
	benchList(b, 1000)
}

// Unsure if this is really valid per the language spec.  AFAIK, it offers no
// concrete guarantee of ordering and sequencing of elements.  Works in
// practice?  There was a thread involving RSC and Ian on the topic a while
// back.
func benchChan(b *testing.B, cap int) {
	b.StopTimer()
	ch := make(chan int, cap)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
	inner:
		select {
		case ch <- i:
			break inner
		default:
			<-ch
			ch <- i
		}
	}
}

func BenchmarkChan1(b *testing.B) {
	benchChan(b, 1)
}

func BenchmarkChan10(b *testing.B) {
	benchChan(b, 10)
}

func BenchmarkChan100(b *testing.B) {
	benchChan(b, 100)
}

func BenchmarkChan1000(b *testing.B) {
	benchChan(b, 1000)
}

// TODO(mtp): Benchmark custom single-linked list with support for only the ops
//            we need.
