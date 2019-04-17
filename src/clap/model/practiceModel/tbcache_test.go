package practiceModel

import (
	"testing"
	"reflect"
)

type data struct {
	Hp int
}

func BenchmarkGetallrec(b *testing.B){

	v := data{Hp:2}

	b.StopTimer()
	b.ResetTimer()
	b.StartTimer()

	for i:=0 ;i < b.N;i++{
		v.Hp = 3
	}
}

func BenchmarkTestReflect(b *testing.B){
	v:=data{Hp:2}

	vv := reflect.ValueOf(&v).Elem()
	f:=vv.FieldByName("Hp")

	b.StopTimer()
	b.ResetTimer()
	b.StartTimer()

	for i :=0;i<b.N;i++{
		f.SetInt(3)
	}
}