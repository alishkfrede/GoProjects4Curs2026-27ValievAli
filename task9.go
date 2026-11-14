package main

import "fmt"

type resource interface {
	name() string
	close() error
}

type file struct {
	fname string
}

func (f *file) name() string  { return f.fname }
func (f *file) close() error  { fmt.Println("closing file", f.fname); return nil }

type connection struct {
	addr string
}

func (c *connection) name() string { return c.addr }
func (c *connection) close() error { fmt.Println("closing conn", c.addr); return nil }

// ресурс, который «ломается» при закрытии
type brokenRes struct{ id int }

func (b *brokenRes) name() string { return fmt.Sprintf("broken-%d", b.id) }
func (b *brokenRes) close() error {
	fmt.Println("closing broken", b.id)
	return fmt.Errorf("boom")
}

func useAll(resources []resource) {
	// закрываем всё в обратном порядке, даже если были ошибки
	defer func() {
		for i := len(resources) - 1; i >= 0; i-- {
			if err := resources[i].close(); err != nil {
				fmt.Println("close error:", err)
			}
		}
	}()

	for _, r := range resources {
		fmt.Println("using", r.name())
	}
}

func main() {
	rs := []resource{
		&file{"a.txt"},
		&connection{"localhost:8080"},
		&brokenRes{1},
	}
	useAll(rs)
}