package main

import "fmt"

type reader interface {
	next() (int, bool)
}

type writer interface {
	write(n int)
}

type sliceReader struct {
	data []int
	pos  int
}

func (r *sliceReader) next() (int, bool) {
	if r.pos >= len(r.data) {
		return 0, false
	}
	v := r.data[r.pos]
	r.pos++
	return v, true
}

type sumWriter struct {
	sum int
}

func (w *sumWriter) write(n int) { w.sum += n }
func (w *sumWriter) total() int  { return w.sum }

func pipeline(r reader, w writer) {
	for {
		v, ok := r.next()
		if !ok {
			return
		}
		w.write(v)
	}
}

func main() {
	r := &sliceReader{data: []int{1, 2, 3, 4, 5, -3}}
	w := &sumWriter{}
	pipeline(r, w)
	fmt.Println("итоговая сумма:", w.total())
}
