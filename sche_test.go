package sche

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var v = 0
var mu sync.Mutex

func f() {
	fmt.Println("f()")
	mu.Lock()
	v++
	mu.Unlock()
}
func f1(i int) {
	fmt.Println("f1(", i, ")")
	mu.Lock()
	v++
	mu.Unlock()
}
func f2(i int, s string) {
	fmt.Println("f2(", i, s, ")")
	mu.Lock()
	v++
	mu.Unlock()
}

func TestSche(t *testing.T) {
	Sche.Push().After(1 * time.Second).Do(f)
	Sche.Push().After(1*time.Second).Do(f1, 0)
	Sche.Push().After(1*time.Second).Do(f2, 1, "***")
	Sche.Push().After(1*time.Second).Do(f2, 2, "***")
	Sche.Push().After(1*time.Second).Do(f2, 3, "***")
	//Sche.Push().Every(3 * time.Second).Do(f)
	//Sche.Push().Every(2*time.Second).DoTimes(3, f)
	//Sche.Push().Every(1).Second().Do(f).ForTimes(3)
	//Sche.Push().At(20,35,59).Do(f)
	//Sche.Push().Every().Day().At(20, 40, 0).Do(f)
	//Sche.Push().Every().Wednesday().At(20,56,40).Do(f)
	//Sche.Push().Every(2).Wednesday().At(21,0,40).Do(f)
	go Sche.Run()

	time.Sleep(10 * time.Second)
}

func TestJob_check(t *testing.T) {
	j := newJob()

	var a1 []interface{}
	a1 = append(a1, 100)
	j.check(f1, a1)
	j.handler.Call(j.in)

	var a2 []interface{}
	a2 = append(a2, 100)
	a2 = append(a2, "string")
	j.check(f2, a2)
	j.handler.Call(j.in)
}
