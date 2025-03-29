package utils

import (
	"log"
	"testing"
	"time"
)

func TestFuturePromiseCummulative(t *testing.T) {
	p := NewPromise[int](4)

	f := p.Future()
	if !f.IsEmpty() {
		t.Fatal("Fatal #1")
	}

	if !f.IsOpen() {
		t.Fatal("Fatal #2")
	}

	res := p.Put(10)

	if f.IsEmpty() || !res {
		t.Fatal("Fatal #3")
	}

	if f.GetOr(100) != 10 {
		t.Fatal("Fatal #4")
	}

	if !f.IsEmpty() {
		t.Fatal("Fatal #5")
	}

	fres := f.Result()

	if fres.err == nil {
		t.Fatal("Fatal #6")
	}

	res = p.Put(10)

	if f.IsEmpty() || !res {
		t.Fatal("Fatal #7")
	}

	fres = f.Result()
	if fres.err != nil || fres.result != 10 {
		t.Fatal("Fatal #8")
	}

	if !f.IsEmpty() {
		t.Fatal("Fatal #9")
	}

	go func(p *Promise[int]) {
		time.Sleep(time.Second * 1)
		pres := p.Put(20)
		if !pres {
			log.Fatalln("Log #16")
		}
	}(&p)

	fval := f.Get()

	if fval != 20 {
		t.Fatal("Fatal #10")
	}

	go func(p *Promise[int]) {
		time.Sleep(time.Second * 1)
		pres := p.Put(30)
		if !pres {
			log.Fatalln("Log 17")
		}
	}(&p)

	fres = f.WaitResult()

	if fres.result != 30 || fres.err != nil {
		t.Fatal("Fatal #12")
	}

	res = p.Put(10)

	if !f.IsEmpty() || res {
		t.Fatal("Fatal #14")
	}

	f.Close()

	if f.IsOpen() || f.GetOr(100) != 100 || f.Result().err == nil || f.WaitResult().err == nil {
		t.Fatal("Fatal #15")
	}
}
