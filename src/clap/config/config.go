package config

import (
	"clap/logger"
	"fmt"
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
func LoadConfig() {
	logger.Infoln("LoadConfig......Hello World!")
	c := make(chan int, 1)
	c <- 1
	fmt.Println(<-c)
}

func sum(a []int, c chan int) {
	sum := 0
	for _, v := range a {
		sum += v
	}
	c <- sum
}
