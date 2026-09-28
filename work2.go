package main

import(
	"fmt"
	"math"
)

type Shape interface{
	Area() float64
	Perimeter() float64
}

type Circle struct{
	Radius float64
}

func(c Circle) Area() float64 {return math.Pi * c.Radius * c.Radius}
func(c Circle) Perimeter() float64 {return 2*math.Pi * c.Radius}

type Rectangle struct{
	Width, Height float64
}

func (r Rectangle) Area() float64 {return r.Width * r.Height}
func (r Rectangle) Perimeter() float64 {return 2 * (r.Width + r.Height)}

type Triangle struct{
	A, B, C float64
}

func (t Triangle) Perimeter() float64 {return t.A + t.B + t.C}

func (t Triangle) Area() float64{
	p := t.Perimeter() / 2
	return math.Sqrt(p * (p - t.A) * (p-t.B) * (p-t.C))
}

func Total(shapes []Shape) (area, perimeter float64){
	for _, s := range shapes{
		area += s.Area()
		perimeter += s.Perimeter()
	}
	return
}

func main(){
	shapes := []Shape{
		Circle{Radius: 1},
		Rectangle{Width: 3, Height: 4},
		Triangle{A:3, B:4, C:5},
	}

	area, perimeter := Total(shapes)
	fmt.Printf("Площадь: %.4f\n", area)
	fmt.Printf("Периметр: %.4f\n", perimeter)
}