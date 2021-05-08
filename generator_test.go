package pygo

import (
	"fmt"
	"testing"
)

const (
	benchStop  = 1000
	benchChunk = 1000
)

func compareSliceInt(a, b []int) error {
	if n, m := len(a), len(b); n != m {
		return fmt.Errorf("failed slice lenght %v != %v", n, m)
	}
	for j, av := range a {
		if bv := b[j]; av != bv {
			return fmt.Errorf("failed slice compare %v != %v", av, bv)
		}
	}
	return nil
}

func compareSlices(a, b [][]int) error {
	if n, m := len(a), len(b); n != m {
		return fmt.Errorf("failed slice lenght %v != %v", n, m)
	}
	for j := range a {
		if err := compareSliceInt(a[j], b[j]); err != nil {
			return fmt.Errorf("failed slice compare %v != %v: %w", a[j], b[j], err)
		}
	}
	return nil
}

func TestXRange(t *testing.T) {
	cases := []struct {
		start    int
		stop     int
		step     int
		expected []int
		err      error
	}{
		{0, 0, 1, []int{}, nil},
		{0, 1, 1, []int{0}, nil},
		{0, 1, 0, nil, ErrOffsetIteration},
		{0, 1, -1, nil, ErrOffsetIteration},
		{0, 0, 1, []int{}, nil},
		{2, 0, 1, []int{}, nil},
		{0, 3, 1, []int{0, 1, 2}, nil},
		{3, 12, 2, []int{3, 5, 7, 9, 11}, nil},
		{-7, 10, 5, []int{-7, -2, 3, 8}, nil},
	}
	for i, c := range cases {
		g, err := XRange(c.start, c.stop, c.step)
		if err != c.err {
			t.Errorf("failed [%d] error check: %v != %v", i, err, c.err)
			continue
		}
		if err != nil {
			continue
		}
		// positive cases
		output := make([]int, 0, len(c.expected))
		for j := range g {
			output = append(output, j)
		}
		err = compareSliceInt(output, c.expected)
		if err != nil {
			t.Errorf("failed [%d]: %v", i, err)
		}
	}
}

func BenchmarkXRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g, err := XRange(0, benchStop, 1)
		if err != nil {
			b.Errorf("failed: %v", err)
		}
		j := 0
		for k := range g {
			if k != j {
				b.Errorf("failed %v != %v", k, j)
			}
			j++
		}
	}
}

func TestChunk(t *testing.T) {
	cases := []struct {
		items    []int
		size     int
		expected [][]int
		err      error
	}{
		{[]int{}, 1, [][]int{}, nil},
		{[]int{1, 2}, 1, [][]int{{1}, {2}}, nil},
		{[]int{1, 2}, 0, [][]int{{1}, {2}}, ErrOffsetIteration},
		{[]int{1, 2}, -1, [][]int{{1}, {2}}, ErrOffsetIteration},
		{[]int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}, nil},
		{[]int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}, nil},
		{[]int{1, 2, 3, 4, 5}, 3, [][]int{{1, 2, 3}, {4, 5}}, nil},
		{[]int{1, 2, 3}, 5, [][]int{{1, 2, 3}}, nil},
	}
	for i, c := range cases {
		chunks, err := Chunk(c.items, c.size)
		if err != c.err {
			t.Errorf("failed [%d] error check: %v != %v", i, err, c.err)
			continue
		}
		if err != nil {
			continue
		}
		// positive cases
		output := make([][]int, 0, len(c.expected))
		for j := range chunks {
			output = append(output, j)
		}
		err = compareSlices(c.expected, output)
		if err != nil {
			t.Errorf("failed [%d]: %v", i, err)
		}
	}
}

func BenchmarkChunk(b *testing.B) {
	items := make([]int, benchChunk)
	for i := range items {
		items[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g, err := Chunk(items, 2)
		if err != nil {
			b.Errorf("failed: %v", err)
		}
		j := 0
		for k := range g {
			err = compareSliceInt(k, []int{j, j + 1})
			if err != nil {
				b.Errorf("failed %v", err)
			}
			j += 2
		}
	}
}

func TestGenerator(t *testing.T) {
	var (
		j   int
		err error
	)
	cases := []struct {
		start    int
		stop     int
		step     int
		expected []int
		err      error
	}{
		{0, 0, 1, []int{}, nil},
		{0, 1, 1, []int{0}, nil},
		{0, 1, 0, nil, ErrOffsetIteration},
		{0, 1, -1, nil, ErrOffsetIteration},
		{0, 0, 1, []int{}, nil},
		{2, 0, 1, []int{}, nil},
		{0, 3, 1, []int{0, 1, 2}, nil},
		{3, 12, 2, []int{3, 5, 7, 9, 11}, nil},
		{-7, 10, 5, []int{-7, -2, 3, 8}, nil},
	}
	for i, c := range cases {
		output := make([]int, 0, len(c.expected))
		g := Generator(c.start, c.stop, c.step)
		for j, err = g(); err == nil; j, err = g() {
			output = append(output, j)
		}
		if (err != ErrStopIteration) && (err != c.err) {
			t.Errorf("failed [%d] error check: %v != %v", i, err, c.err)
			continue
		}
		if c.err != nil {
			continue
		}
		err = compareSliceInt(output, c.expected)
		if err != nil {
			t.Errorf("failed [%d]: %v", i, err)
		}
	}
}

func BenchmarkGenerator(b *testing.B) {
	var (
		j   int
		err error
	)
	for i := 0; i < b.N; i++ {
		g := Generator(0, benchStop, 1)
		k := 0
		for j, err = g(); err == nil; j, err = g() {
			if k != j {
				b.Errorf("failed %v != %v", k, j)
			}
			k++
		}
		if err != ErrStopIteration {
			b.Errorf("failed erorr: %v", err)
		}
	}
}

func TestChunkGenerator(t *testing.T) {
	var (
		chunks []int
		err    error
	)
	cases := []struct {
		items    []int
		size     int
		expected [][]int
		err      error
	}{
		{[]int{}, 1, [][]int{}, nil},
		{[]int{1, 2}, 1, [][]int{{1}, {2}}, nil},
		{[]int{1, 2}, 0, [][]int{{1}, {2}}, ErrOffsetIteration},
		{[]int{1, 2}, -1, [][]int{{1}, {2}}, ErrOffsetIteration},
		{[]int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}, nil},
		{[]int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}, nil},
		{[]int{1, 2, 3, 4, 5}, 3, [][]int{{1, 2, 3}, {4, 5}}, nil},
		{[]int{1, 2, 3}, 5, [][]int{{1, 2, 3}}, nil},
	}
	for i, c := range cases {
		output := make([][]int, 0, len(c.expected))
		g := ChunkGenerator(c.items, c.size)
		for chunks, err = g(); err == nil; chunks, err = g() {
			output = append(output, chunks)
		}
		if (err != ErrStopIteration) && (err != c.err) {
			t.Errorf("failed [%d] error check: %v != %v", i, err, c.err)
			continue
		}
		if c.err != nil {
			continue
		}
		// positive cases
		err = compareSlices(c.expected, output)
		if err != nil {
			t.Errorf("failed [%d]: %v", i, err)
		}
	}
}

func BenchmarkChunkGenerator(b *testing.B) {
	var (
		chunks []int
		err    error
	)
	items := make([]int, benchChunk)
	for i := range items {
		items[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g := ChunkGenerator(items, 2)
		j := 0
		for chunks, err = g(); err == nil; chunks, err = g() {
			if e := compareSliceInt(chunks, []int{j, j + 1}); e != nil {
				b.Errorf("failed %v", e)
			}
			j += 2
		}
		if err != ErrStopIteration {
			b.Errorf("failed stop %v", err)
			continue
		}
	}
}

func TestNewGenStruct(t *testing.T) {
	cases := []struct {
		start    int
		stop     int
		step     int
		expected []int
		err      error
	}{
		{0, 0, 1, []int{}, nil},
		{0, 1, 1, []int{0}, nil},
		{0, 1, 0, nil, ErrOffsetIteration},
		{0, 1, -1, nil, ErrOffsetIteration},
		{0, 0, 1, []int{}, nil},
		{2, 0, 1, []int{}, nil},
		{0, 3, 1, []int{0, 1, 2}, nil},
		{3, 12, 2, []int{3, 5, 7, 9, 11}, nil},
		{-7, 10, 5, []int{-7, -2, 3, 8}, nil},
	}
	for i, c := range cases {
		g, err := NewGenStruct(c.start, c.stop, c.step)
		if err != c.err {
			t.Errorf("failed [%d] error check: %v != %v", i, err, c.err)
			continue
		}
		if err != nil {
			continue
		}
		// positive cases
		output := make([]int, 0, len(c.expected))
		for v, ok := g.Next(); ok; v, ok = g.Next() {
			output = append(output, v)
		}
		err = compareSliceInt(output, c.expected)
		if err != nil {
			t.Errorf("failed [%d]: %v", i, err)
		}
	}
}

func BenchmarkGenStruct_Next(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g, err := NewGenStruct(0, benchStop, 1)
		if err != nil {
			b.Error(err)
		}
		k := 0
		for v, ok := g.Next(); ok; v, ok = g.Next() {
			if k != v {
				b.Errorf("failed %v != %v", k, v)
			}
			k++
		}
	}
}

func TestNewGenStructChunk(t *testing.T) {
	cases := []struct {
		items    []int
		size     int
		expected [][]int
		err      error
	}{
		{[]int{}, 1, [][]int{}, nil},
		{[]int{1, 2}, 1, [][]int{{1}, {2}}, nil},
		{[]int{1, 2}, 0, [][]int{{1}, {2}}, ErrOffsetIteration},
		{[]int{1, 2}, -1, [][]int{{1}, {2}}, ErrOffsetIteration},
		{[]int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}, nil},
		{[]int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}, nil},
		{[]int{1, 2, 3, 4, 5}, 3, [][]int{{1, 2, 3}, {4, 5}}, nil},
		{[]int{1, 2, 3}, 5, [][]int{{1, 2, 3}}, nil},
	}
	for i, c := range cases {
		g, err := NewGenStructChunk(c.items, c.size)
		if err != c.err {
			t.Errorf("failed [%d] error check: %v != %v", i, err, c.err)
			continue
		}
		if err != nil {
			continue
		}
		// positive cases
		output := make([][]int, 0, len(c.expected))
		for v, ok := g.NextChunk(); ok; v, ok = g.NextChunk() {
			output = append(output, v)
		}
		err = compareSlices(c.expected, output)
		if err != nil {
			t.Errorf("failed [%d]: %v", i, err)
		}
	}
}

func BenchmarkGenStruct_NextChunk(b *testing.B) {
	items := make([]int, benchChunk)
	for i := range items {
		items[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g, err := NewGenStructChunk(items, 2)
		if err != nil {
			b.Errorf("failed: %v", err)
		}
		j := 0
		for v, ok := g.NextChunk(); ok; v, ok = g.NextChunk() {
			err = compareSliceInt(v, []int{j, j + 1})
			if err != nil {
				b.Errorf("failed %v", err)
			}
			j += 2
		}
	}
}
