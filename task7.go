package main

import "fmt"

type Record struct {
	Name string
	Age  int
}

type validator interface {
	valid(r Record) bool
}

type nonEmptyName struct{}

func (nonEmptyName) valid(r Record) bool { return r.Name != "" }

type ageRange struct{}

func (ageRange) valid(r Record) bool { return r.Age >= 18 && r.Age <= 99 }

func validateAll(records []Record, v validator) []Record {
	res := make([]Record, 0)
	for _, r := range records {
		if v.valid(r) {
			res = append(res, r)
		}
	}
	return res
}

func main() {
	records := []Record{
		{"Алиса", 25},
		{"", 30},
		{"Боб", 17},
		{"Карл", 100},
		{"Дина", 40},
	}
	fmt.Println("с именем:", validateAll(records, nonEmptyName{}))
	fmt.Println("18-99:", validateAll(records, ageRange{}))
}