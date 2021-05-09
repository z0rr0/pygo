package pygo

import "errors"

var (
	// ErrStopIteration is error when iteration is finished.
	ErrStopIteration = errors.New("stop iteration")
	// ErrOffsetIteration is error for failed step/size value.
	ErrOffsetIteration = errors.New("iteration offset must be positive")
)

// XRange is a range generator.
func XRange(start, stop, step int) (chan int, error) {
	if step < 1 {
		return nil, ErrOffsetIteration
	}
	c := make(chan int)
	go func() {
		for j := start; j < stop; j += step {
			c <- j
		}
		close(c)
	}()
	return c, nil
}

// Chunk splits items by chunks with maximum length=size.
func Chunk(items []int, size int) (chan []int, error) {
	if size < 1 {
		return nil, ErrOffsetIteration
	}
	c := make(chan []int)
	start, stop := 0, len(items)
	go func() {
		defer close(c)
		for j := start; j < stop; j += size {
			step := j + size
			if step > stop {
				c <- items[j:stop]
				return // last chunk
			}
			c <- items[j:step]
		}
	}()
	return c, nil
}

// Generator is a function closure int generator.
func Generator(start, stop, step int) func() (int, error) {
	if step < 1 {
		return func() (int, error) { return 0, ErrOffsetIteration }
	}
	i := start
	return func() (int, error) {
		defer func() { i += step }()
		if i >= stop {
			step = 0
			return 0, ErrStopIteration
		}
		return i, nil
	}
}

// ChunkGenerator is a function closure chunk splitter.
func ChunkGenerator(items []int, size int) func() ([]int, error) {
	if size < 1 {
		return func() ([]int, error) { return nil, ErrOffsetIteration }
	}
	j, stop := 0, len(items)
	return func() ([]int, error) {
		defer func() { j += size }()
		if j >= stop {
			size = 0
			return nil, ErrStopIteration
		}
		step := j + size
		if step > stop {
			step = stop
		}
		return items[j:step], nil
	}
}

// GenStruct is a struct generator.
type GenStruct struct {
	stop  int
	step  int
	value int
	items []int
}

// NewGenStruct returns a new struct generator.
func NewGenStruct(start, stop, step int) (*GenStruct, error) {
	if step < 1 {
		return nil, ErrOffsetIteration
	}
	return &GenStruct{stop: stop, step: step, value: start}, nil
}

// NewGenStructChunk returns new chunk splitter.
func NewGenStructChunk(items []int, size int) (*GenStruct, error) {
	if size < 1 {
		return nil, ErrOffsetIteration
	}
	return &GenStruct{stop: len(items), step: size, items: items}, nil
}

// Next returns a new generation value and flag that it is not the end.
func (g *GenStruct) Next() (int, bool) {
	defer func() { g.value += g.step }()
	return g.value, g.value < g.stop
}

// NextChunk returns a new generated chunk and flag that it is not the end.
func (g *GenStruct) NextChunk() ([]int, bool) {
	var i = g.value + g.step
	if i > g.stop {
		i = g.stop
	}
	defer func() { g.value = i }()
	return g.items[g.value:i], g.value < g.stop
}
