package pygo

import "errors"

var ErrStopIteration = errors.New("stop iteration")

func XRange(start, stop, step int) chan int {
	var c chan int
	go func() {
		for j := start; j < stop; j += step {
			c <- j
		}
		close(c)
	}()
	return c
}

func Chunk(items []int, size int) chan []int {
	var c chan []int
	start, stop := 0, len(items)-1
	go func() {
		defer close(c)
		for j := start; j < stop; j += size {
			step := j + size
			if step > stop {
				c <- items[j:stop]
				return
			}
			c <- items[j:step]
		}
	}()
	return c
}

func Generator(start, stop, step int) func() (int, error) {
	var i = start
	return func() (int, error) {
		defer func() {
			i += step
		}()
		if i >= stop {
			return 0, ErrStopIteration
		}
		return i, nil
	}
}

func ChunkGenerator(items []int, size int) func() ([]int, error) {
	var j, stop = 0, len(items) - 1
	return func() ([]int, error) {
		defer func() {
			j += size
		}()
		if j > stop {
			return nil, ErrStopIteration
		}
		step := j + size
		if step > stop {
			step = stop
		}
		return items[j:step], nil
	}
}
