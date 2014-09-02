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
	"bytes"
	"container/list"
	"fmt"
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

// Channels act as first-in-first-out queues. For example, if one goroutine
// sends values on a channel and a second goroutine receives them, the values
// are received in the order sent.  -  https://golang.org/ref/spec#Channel_types
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

type intListNode struct {
	val  int
	next *intListNode
}

type intList struct {
	len         int
	front, back *intListNode
}

func (l *intList) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "[intList len=%d", l.len)
	for nod := l.front; nod != nil; nod = nod.next {
		fmt.Fprintf(&buf, " %d", nod.val)
	}
	fmt.Fprintf(&buf, "]")
	return buf.String()
}

func (l *intList) pushBack(v int) {
	l.len++
	nod := &intListNode{val: v}
	if l.front == nil && l.back == nil {
		l.front = nod
		l.back = nod
		return
	}
	l.back.next = nod
	l.back = nod
}

func (l *intList) moveToBack(n *intListNode) {
	if n.next == nil { // one-off for cap == 1
		return
	}
	l.front = n.next // Assuming a FIFO and n == l.front
	l.back.next = n
	n.next = nil
}

// Allocate only on warmup, if at all.
func benchIntList(b *testing.B, cap int) {
	b.StopTimer()
	lst := intList{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if l := lst.len; l < cap {
			lst.pushBack(i)
		} else {
			frnt := lst.front
			frnt.val = i
			lst.moveToBack(frnt)
		}
	}
}

func BenchmarkIntList1(b *testing.B) {
	benchIntList(b, 1)
}

func BenchmarkIntList10(b *testing.B) {
	benchIntList(b, 10)
}

func BenchmarkIntList100(b *testing.B) {
	benchIntList(b, 100)
}

func BenchmarkIntList1000(b *testing.B) {
	benchIntList(b, 1000)
}
