[![GoDoc](https://godoc.org/github.com/geozelot/boxtree?status.svg)](https://godoc.org/github.com/geozelot/boxtree) [![Go Report Card](https://goreportcard.com/badge/github.com/geozelot/boxtree)](https://goreportcard.com/report/github.com/geozelot/boxtree)

# BOXTree for Go

Very fast static, flat **BOX** **Tree** implementation for reverse 2D range searches (box overlap).

**Suitable for high performance spatial bounding box searches for point coordinate input!**

Highly efficient and with almost no memory footprint other than the stored ranges.

Further scientific reading about the adapted traversal algorithm and comparisons between different approaches (in C/C++) can be found [here](https://github.com/lh3/cgranges).


# Behaviour

* BOXTree will build the tree once (**static; no updates after creation**)
* BOXTree returns indices to the initial `[]Box` array
* BOXTree currently supports finding all boxes for a single `[]float64` value pair

# Usage

## API ([GoDoc](https://godoc.org/github.com/geozelot/boxtree))

### `type Box`

`Box{}` is the main interface expected by `NewBOXTree()`; requires `Limits()` method to access box limits.

```go
type Box interface {
    Limits() (Lower, Upper []float64)
}
```

### `type BOXTree`

`BOXTree{}` is the main package object; holds Slice of reference indices and the respective box limits.

```go
type BOXTree struct {
    // contains filtered or unexported fields
}
```

### `func NewBOXTree`

`NewBOXTree()` is the main initialization function; creates the tree from the given Slice of Box.

```go
func NewBOXTree(bxs []Box) *BOXTree
```

### `func (*BOXTree) Overlaps`

`Overlaps()` is the main entry point for box searches; traverses the tree and collects boxes that overlap with the given values.

```go
func (inT *BOXTree) Overlaps(vals []float64) []int
```

## Import
```go
import (
    "github.com/geozelot/boxtree"
)
```

## Examples

#### Simple `Box{}` interface implementation:

```go
// SimpleBox is a simple Struct implicitly implementing the Box interface.
type SimpleBox struct {

  MinX, MinY, MaxX, MaxY float64

}

// Limits accesses the box limits.
func (sb *SimpleBox) Limits() (Lower, Upper []float64) {

  return []float64{ sb.MinX, sb.MinY }, []float64{ sb.MaxX, sb.MaxY }

}
```

#### Test Setup:

```go
package main

import (

    "github.com/geozelot/boxtree"
    "fmt"

)

// defining simple Struct holding box limits
type SimpleBox struct {
  MinX, MinY, MaxX, MaxY float64
}

  // add method to access limits; implicitly implements BOXTree.Box interface
  func (sb *SimpleBox) Limits() (Lower, Upper []float64) {
    return []float64{ sb.MinX, sb.MinY }, []float64{ sb.MaxX, sb.MaxY }
  }

func main() {

  // create typed var
  var tree *boxtree.BOXTree
  
  // create example boxes
  inputBoxes := []boxtree.Box{

    &SimpleBox{ MinX: 4.0, MinY: 6.0, MaxX: 8.0,  MaxY: 10.0 },
    &SimpleBox{ MinX: 5.0, MinY: 5.0, MaxX: 11.0, MaxY: 9.0  },
    &SimpleBox{ MinX: 1.0, MinY: 4.0, MaxX: 4.0,  MaxY: 7.0  },     // match
    &SimpleBox{ MinX: 2.0, MinY: 3.0, MaxX: 3.0,  MaxY: 4.0  },
    &SimpleBox{ MinX: 4.0, MinY: 6.0, MaxX: 8.0,  MaxY: 10.0 },
    &SimpleBox{ MinX: 6.0, MinY: 3.0, MaxX: 8.0,  MaxY: 8.0  },
    &SimpleBox{ MinX: 2.0, MinY: 6.0, MaxX: 7.0,  MaxY: 7.0  },     // match

  }

  // initialize new BOXTree and create tree from inputBoxes
  tree = boxtree.NewBOXree(inputBoxes)

  // create point to match
  point := []float64{ 3.2, 6.3 }

  // parse return Slice with indices referencing inputBoxes
  for _, matchedIndex := range tree.Overlaps(point) {

    // using BOXTree.Box interface method to access limits
    lowerLimits, upperLimits := inputBoxes[matchedIndex].Limits()

    fmt.Printf("Match at inputBoxes index %2d with 2D range [ %v, %v ]\n", matchedIndex, lowerLimits, upperLimits)

    /*
      Match at inputBoxes index  6 with 2D range [ [2 6], [7 7] ]
      Match at inputBoxes index  2 with 2D range [ [1 4], [4 7] ]
    */

  }

}
```

#### Try on [Go Playground](https://play.golang.org/p/xeVFUX1m5vS).

____

##### Inspired by this great [KDTree implementation](https://github.com/mourner/kdbush) for JavaScript and adapted from this excellent [Go port](https://github.com/MadAppGang/kdbush).
