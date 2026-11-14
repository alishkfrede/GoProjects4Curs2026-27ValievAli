package main

import "fmt"

type source interface {
	next() ([]string, bool) // порция, есть ли ещё
}

// источник из среза, отдаёт по одной строке
type sliceSource struct {
	data []string
	pos  int
}

func (s *sliceSource) next() ([]string, bool) {
	if s.pos >= len(s.data) {
		return nil, false
	}
	item := s.data[s.pos]
	s.pos++
	return []string{item}, true
}

// генерирующий источник
type genSource struct {
	limit int
	cur   int
}

func (g *genSource) next() ([]string, bool) {
	if g.cur >= g.limit {
		return nil, false
	}
	g.cur++
	return []string{fmt.Sprintf("item-%d", g.cur)}, true
}

func collectAll(s source) []string {
	var all []string
	for {
		chunk, ok := s.next()
		if !ok {
			break
		}
		all = append(all, chunk...)
	}
	return all
}

func main() {
	s1 := &sliceSource{data: []string{"alpha", "beta", "gamma"}}
	s2 := &genSource{limit: 4}

	fmt.Println("slice:", collectAll(s1))
	fmt.Println("gen:", collectAll(s2))
}