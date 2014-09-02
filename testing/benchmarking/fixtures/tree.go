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
package tree

import (
	"math/rand"
)

// Node models a node in a binary tree.
type Node struct {
	Val         int
	Left, Right *Node
}

// TraverseRecursive traverses a binary tree following "pre-order" semantics,
// emitting all visited data to the channel out.  It does not discard elements
// if the channel blocks but rather blocks.  It closes the channel when it is
// finished.  Under the hood, it uses recursion to achieve the result.
func TraverseRecursive(n *Node, out chan<- int) {
	traverseRecursive(n, out)
	close(out)
}

func traverseRecursive(n *Node, out chan<- int) {
	if n == nil {
		return
	}
	out <- n.Val
	traverseRecursive(n.Left, out)
	traverseRecursive(n.Right, out)
}

type stack []*Node

func (s *stack) Push(n *Node) {
	*s = append(*s, n)
}

func (s *stack) Pop() (*Node, bool) {
	stk := *s
	l := len(stk)
	if l == 0 {
		return nil, false
	}
	n := stk[l-1]
	*s = stk[:l-1]
	return n, true
}

// TraverseIterative traverses a binary tree following "pre-order" semantics,
// emitting all visited data to the channel out.  It does not discard elements
// if the channel blocks but rather blocks.  It closes the channel when it is
// finished.  Under the hood, it uses loop iterations to achieve the result.
func TraverseIterative(n *Node, out chan<- int) {
	var stk stack
	stk.Push(n)
	for {
		nod, ok := stk.Pop()
		if !ok {
			break
		}
		out <- nod.Val
		if nod.Left != nil {
			stk.Push(nod.Left)
		}
		if nod.Right != nil {
			stk.Push(nod.Right)
		}
	}
	close(out)
}

var rnd = rand.New(rand.NewSource(42))

const (
	mkDbl int = iota
	mkLft
	mkRght
)

// New creates a new binary tree composed of randomized nodes of size n.
func New(n int) *Node {
	var stk stack
	prnt := &Node{Val: rnd.Intn(255)}
	stk.Push(prnt)
	n--
	for n > 0 {
		nod, ok := stk.Pop()
		if !ok {
			break
		}
		typ := rnd.Intn(mkRght + 1)
		switch {
		case typ == mkDbl && n >= 2:
			nod.Left = &Node{Val: rnd.Intn(255)}
			nod.Right = &Node{Val: rnd.Intn(255)}
			stk.Push(nod.Left)
			stk.Push(nod.Right)
			n -= 2
		case typ == mkLft:
			nod.Left = &Node{Val: rnd.Intn(255)}
			stk.Push(nod.Left)
			n--
		case typ == mkRght:
			nod.Right = &Node{Val: rnd.Intn(255)}
			stk.Push(nod.Right)
			n--
		}
	}
	return prnt
}

// Traverser is a function that is capable of traversing the binary tree.
type Traverser func(*Node, chan<- int)
