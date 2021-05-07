package pygo

import "errors"

var (
	// ErrStopIteration is error when iteration is finished.
	ErrStopIteration = errors.New("stop iteration")
	// ErrOffsetIteration is error for failed step/size value.
	ErrOffsetIteration = errors.New("iteration offset must be positive")
)

// XRange is range generator.
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

// Generator is function closure int generator.
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

// ChunkGenerator is function closure chunk splitter.
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
