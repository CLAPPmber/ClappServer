package practiceModel

import (
	"testing"
	"math/rand"
)

func BenchmarkGetallrec(b *testing.B){

	cluesrs := []Cluser{
		{Account:"userone",Password:"123123"},
		{Account:"usertt1",Password:"123123"},
		{Account:"usernew",Password:"123123"},
		{Account:"admin1",Password:"asdasd"},
		{Account:"admin123",Password:"asdasd"},
		{Account:"admliang",Password:"asdasd"},
	}
	b.StartTimer()
	for i:=0;i<b.N;i++{
		b.StopTimer()
		index := rand.Intn(6)
		b.StartTimer()
		Getallrec(cluesrs[index])
	}
}