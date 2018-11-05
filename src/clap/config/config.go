package config

import (
	"strconv"
)

type Human struct {
	string
	pp
	age   int
	grade int
}
type pp struct {
	pname string
	shuak int
}

func (Hu Human) String() string {
	return "Hname: " + Hu.pname + " age: " + strconv.Itoa(Hu.age)
}

//LoadConfig println fmt

func sum(a []int, c chan int) {
	sum := 0
	for _, v := range a {
		sum += v
	}
	c <- sum
}
