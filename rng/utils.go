package rng

// array-based stack
type worklist []int

func (w *worklist) push(i int) {
	*w = append(*w, i)
}

func (w *worklist) pop() int {
	l := len(*w) - 1
	n := (*w)[l]
	(*w) = (*w)[:l]
	return n
}
