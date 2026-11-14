package main

import "fmt"

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

type visitor interface {
	visit(n int)
}

type printVisitor struct{}

func (printVisitor) visit(n int) { fmt.Print(n, " ") }

type sumVisitor struct {
	sum int
}

func (s *sumVisitor) visit(n int) { s.sum += n }

func walk(root *Node, v visitor) {
	if root == nil {
		return
	}
	walk(root.Left, v)
	v.visit(root.Value)
	walk(root.Right, v)
}

func main() {
	root := &Node{4,
		&Node{2, &Node{1, nil, nil}, &Node{3, nil, nil}},
		&Node{6, &Node{5, nil, nil}, &Node{7, nil, nil}},
	}

	fmt.Print("обход: ")
	walk(root, printVisitor{})
	fmt.Println()

	sv := &sumVisitor{}
	walk(root, sv)
	fmt.Println("сумма:", sv.sum)
}