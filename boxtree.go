// MIT License
//
// Copyright (c) 2020 geozelot (André Siefken)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package boxtree provides a very fast, static, flat (augmented) 2D interval Tree for reverse 2D range searches (box overlap).
package boxtree

import (
	"math"
	"math/rand"
)

// Box is the main interface expected by NewBOXTree(); requires Limits method to access box limits.
type Box interface {
	Limits() (Lower, Upper []float64)
}

// BOXTree is the main package object;
// holds Slice of reference indices and the respective box limits.
type BOXTree struct {
	idxs []int
	lmts [][]float64
}

// buildTree is the internal tree construction function;
// creates, sorts and augments nodes into Slices.
func (boT *BOXTree) buildTree(bxs []Box) {

	boT.idxs = make([]int, len(bxs))
	boT.lmts = make([][]float64, 3*len(bxs))

	for i, v := range bxs {

		boT.idxs[i] = i
		l, u := v.Limits()

		boT.lmts[3*i] = l
		boT.lmts[3*i+1] = u
		boT.lmts[3*i+2] = []float64{0}

	}

	sort(boT.lmts, boT.idxs, 0)
	augment(boT.lmts, boT.idxs, 0)

}

// Overlaps is the main entry point for box searches;
// traverses the tree and collects boxes that overlap with the given values.
func (boT *BOXTree) Overlaps(vals []float64) []int {

	stk := []int{0, len(boT.idxs) - 1, 0}
	res := []int{}

	for len(stk) > 0 {

		ax := stk[len(stk)-1]
		stk = stk[:len(stk)-1]
		rb := stk[len(stk)-1]
		stk = stk[:len(stk)-1]
		lb := stk[len(stk)-1]
		stk = stk[:len(stk)-1]

		if lb == rb+1 {
			continue
		}

		cn := int(math.Ceil(float64(lb+rb) / 2.0))
		nm := boT.lmts[3*cn+2][0]

		_ax := (ax + 1) % 2

		if vals[ax] <= nm {

			stk = append(stk, lb)
			stk = append(stk, cn-1)
			stk = append(stk, _ax)

		}

		l := boT.lmts[3*cn]

		if l[ax] <= vals[ax] {

			stk = append(stk, cn+1)
			stk = append(stk, rb)
			stk = append(stk, _ax)

			u := boT.lmts[3*cn+1]

			if vals[ax] <= u[ax] && vals[_ax] <= u[_ax] && l[_ax] <= vals[_ax] {
				res = append(res, boT.idxs[cn])
			}

		}

	}

	return res

}

// NewBOXTree is the main initialization function;
// creates the tree from the given Slice of Box.
func NewBOXTree(bxs []Box) *BOXTree {

	boT := BOXTree{}
	boT.buildTree(bxs)

	return &boT

}

// augment is an internal utility function, adding maximum value of all child nodes to the current node.
func augment(lmts [][]float64, idxs []int, ax int) {

	if len(idxs) < 1 {
		return
	}

	max := 0.0

	for idx := range idxs {

		if lmts[3*idx+1][ax] > max {
			max = lmts[3*idx+1][ax]
		}

	}

	r := len(idxs) >> 1

	lmts[3*r+2][0] = max

	augment(lmts[:3*r], idxs[:r], (ax+1)%2)
	augment(lmts[3*r+3:], idxs[r+1:], (ax+1)%2)

}

// sort is an internal utility function, sorting the tree by lowest limits using Random Pivot QuickSearch
func sort(lmts [][]float64, idxs []int, ax int) {

	if len(idxs) < 2 {
		return
	}

	l, r := 0, len(idxs)-1

	p := rand.Int() % len(idxs)

	idxs[p], idxs[r] = idxs[r], idxs[p]
	lmts[3*p], lmts[3*p+1], lmts[3*p+2], lmts[3*r], lmts[3*r+1], lmts[3*r+2] = lmts[3*r], lmts[3*r+1], lmts[3*r+2], lmts[3*p], lmts[3*p+1], lmts[3*p+2]

	for i := range idxs {

		if lmts[3*i][ax] < lmts[3*r][ax] {

			idxs[l], idxs[i] = idxs[i], idxs[l]
			lmts[3*l], lmts[3*l+1], lmts[3*l+2], lmts[3*i], lmts[3*i+1], lmts[3*i+2] = lmts[3*i], lmts[3*i+1], lmts[3*i+2], lmts[3*l], lmts[3*l+1], lmts[3*l+2]

			l++

		}

	}

	idxs[l], idxs[r] = idxs[r], idxs[l]
	lmts[3*l], lmts[3*l+1], lmts[3*l+2], lmts[3*r], lmts[3*r+1], lmts[3*r+2] = lmts[3*r], lmts[3*r+1], lmts[3*r+2], lmts[3*l], lmts[3*l+1], lmts[3*l+2]

	sort(lmts[:3*l], idxs[:l], (ax+1)%2)
	sort(lmts[3*l+3:], idxs[l+1:], (ax+1)%2)

}
