package gorgonia

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorgonia.org/tensor"
)

type broadcastOpTest struct {
	name string
	a    Value
	b    Value

	// broadcast axes
	left, right []byte

	// results
	ab  Value
	err bool
}

var broadcastAddTests = []broadcastOpTest{
	{name: "vec-mat",
		a:     tensor.New(tensor.WithShape(2), tensor.WithBacking([]float64{100, 200})),
		b:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		left:  []byte{1},
		right: nil,
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{101, 102, 203, 204})),
		err:   false,
	},

	{name: "mat-vec",
		a:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		b:     tensor.New(tensor.WithShape(2), tensor.WithBacking([]float64{100, 200})),
		left:  nil,
		right: []byte{1},
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{101, 102, 203, 204})),
		err:   false,
	},
	{name: "rowvec-mat",
		a:     tensor.New(tensor.WithShape(2, 1), tensor.WithBacking([]float64{100, 200})),
		b:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		left:  []byte{1},
		right: nil,
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{101, 102, 203, 204})),
		err:   false,
	},
	{name: "mat-rowvec",
		a:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		b:     tensor.New(tensor.WithShape(2, 1), tensor.WithBacking([]float64{100, 200})),
		left:  nil,
		right: []byte{1},
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{101, 102, 203, 204})),
		err:   false,
	},
	{name: "colvec-mat",
		a:     tensor.New(tensor.WithShape(1, 2), tensor.WithBacking([]float64{100, 200})),
		b:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		left:  []byte{0},
		right: nil,
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{101, 202, 103, 204})),
		err:   false,
	},
	{name: "mat-colvec",
		a:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		b:     tensor.New(tensor.WithShape(1, 2), tensor.WithBacking([]float64{100, 200})),
		left:  nil,
		right: []byte{0},
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{101, 202, 103, 204})),
		err:   false,
	},

	// TODO (these would give coverage to all broadcast applications)
	// 	vec-3tensor
	// 	3tensor-vec
	// 	mat-3tensor
	// 	3-tensor-mat
	// and their corresponding errors

	// WILL ERR
	{name: "vec-mat- wrong left pattern axis",
		a:     tensor.New(tensor.WithShape(2), tensor.WithBacking([]float64{100, 200})),
		b:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		left:  []byte{0},
		right: nil,
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{101, 102, 203, 204})),
		err:   true,
	},
	{name: "rowvec-mat: wrong axis",
		a:     tensor.New(tensor.WithShape(2, 1), tensor.WithBacking([]float64{100, 200})),
		b:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		left:  []byte{2},
		right: nil,
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{101, 102, 203, 204})),
		err:   true,
	},

	{name: "impossible mat-mat",
		a:     tensor.New(tensor.WithShape(2, 4), tensor.WithBacking([]float64{1, 2, 3, 4, 5, 6, 7, 8})),
		b:     tensor.New(tensor.WithShape(1, 2), tensor.WithBacking([]float64{100, 200})),
		left:  nil,
		right: []byte{0, 1},
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{101, 102, 203, 204})),
		err:   true,
	},
}

func TestBroadcastAdd(t *testing.T) {
	assert := assert.New(t)
	for i, bat := range broadcastAddTests {
		g := NewGraph()
		a := NodeFromAny(g, bat.a, WithName("a"))
		b := NodeFromAny(g, bat.b, WithName("b"))
		c, err := BroadcastAdd(a, b, bat.left, bat.right)
		if checkErr(t, bat.err, err, bat.name, i) {
			continue
		}
		machine := NewTapeMachine(g)

		if err = machine.RunAll(); err != nil {
			t.Errorf("Test %v(%d): %v", bat.name, i, err)
		}
		assert.Equal(bat.ab.Data(), c.Value().Data(), "Test %v(%v)", bat.name, i)
		machine.Close()
	}
}

var broadcastMulTests = []broadcastOpTest{
	{name: "vec-mat",
		a:     tensor.New(tensor.WithShape(2), tensor.WithBacking([]float64{10, 20})),
		b:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		left:  []byte{1},
		right: nil,
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{10, 20, 60, 80})),
		err:   false,
	},

	{name: "mat-vec",
		a:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		b:     tensor.New(tensor.WithShape(2), tensor.WithBacking([]float64{10, 20})),
		left:  nil,
		right: []byte{1},
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{10, 20, 60, 80})),
		err:   false,
	},
	{name: "rowvec-mat",
		a:     tensor.New(tensor.WithShape(2, 1), tensor.WithBacking([]float64{10, 20})),
		b:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		left:  []byte{1},
		right: nil,
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{10, 20, 60, 80})),
		err:   false,
	},
	{name: "mat-rowvec",
		a:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		b:     tensor.New(tensor.WithShape(2, 1), tensor.WithBacking([]float64{10, 20})),
		left:  nil,
		right: []byte{1},
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{10, 20, 60, 80})),
		err:   false,
	},
	{name: "colvec-mat",
		a:     tensor.New(tensor.WithShape(1, 2), tensor.WithBacking([]float64{10, 20})),
		b:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		left:  []byte{0},
		right: nil,
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{10, 40, 30, 80})),
		err:   false,
	},
	{name: "mat-colvec",
		a:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		b:     tensor.New(tensor.WithShape(1, 2), tensor.WithBacking([]float64{10, 20})),
		left:  nil,
		right: []byte{0},
		ab:    tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{10, 40, 30, 80})),
		err:   false,
	},

	// TODO (these would give coverage to all broadcast applications)
	// 	vec-3tensor
	// 	3tensor-vec
	// 	mat-3tensor
	// 	3-tensor-mat
	// and their corresponding errors

	// WILL ERR
	{name: "vec-mat- wrong left pattern axis",
		a:     tensor.New(tensor.WithShape(2), tensor.WithBacking([]float64{10, 20})),
		b:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		left:  []byte{0},
		right: nil,
		err:   true,
	},
	{name: "rowvec-mat: wrong axis",
		a:     tensor.New(tensor.WithShape(2, 1), tensor.WithBacking([]float64{10, 20})),
		b:     tensor.New(tensor.WithShape(2, 2), tensor.WithBacking([]float64{1, 2, 3, 4})),
		left:  []byte{2},
		right: nil,
		err:   true,
	},

	{name: "impossible mat-mat",
		a:     tensor.New(tensor.WithShape(2, 4), tensor.WithBacking([]float64{1, 2, 3, 4, 5, 6, 7, 8})),
		b:     tensor.New(tensor.WithShape(1, 2), tensor.WithBacking([]float64{10, 20})),
		left:  nil,
		right: []byte{0, 1},
		err:   true,
	},
}

func TestBroadcastMul(t *testing.T) {
	assert := assert.New(t)
	for i, bat := range broadcastMulTests {
		g := NewGraph()
		a := NodeFromAny(g, bat.a, WithName("a"))
		b := NodeFromAny(g, bat.b, WithName("b"))
		c, err := BroadcastMul(a, b, bat.left, bat.right)
		if checkErr(t, bat.err, err, bat.name, i) {
			continue
		}
		machine := NewTapeMachine(g)

		if err = machine.RunAll(); err != nil {
			t.Errorf("Test %v(%d): %v", bat.name, i, err)
		}
		assert.Equal(bat.ab.Data(), c.Value().Data(), "Test %v(%v)", bat.name, i)
		machine.Close()
	}
}
