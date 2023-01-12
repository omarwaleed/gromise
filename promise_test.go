package gromise

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestPromiseCreate(t *testing.T) {
	p := New(func(a ...any) any {
		return true
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	p.Run(&wg)
	wg.Wait()
	if p == nil {
		t.Fatal("Promise should not be nil")
	}
	if p.Error != nil {
		t.Fatal("Promise should not have an error")
	}
	if p.Result == nil {
		t.Fatal("Promise should have a result")
	}
	if *p.Result != true {
		t.Fatal("Promise result should be true")
	}
}

func TestPromiseCreateAndReject(t *testing.T) {
	p := New(func(a ...any) any {
		panic(errors.New("Promise rejected"))
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	p.Run(&wg)
	wg.Wait()
	if p == nil {
		t.Fatal("Promise should not be nil")
	}
	if p.Error == nil {
		t.Fatal("Promise should have an error")
	}
}

func TestPromiseAll(t *testing.T) {
	p1 := New(func(a ...any) any {
		return "hello"
	})
	p2 := New(func(a ...any) any {
		return "world"
	})
	promises := []*Promise{p1, p2}
	results, err := All(promises)
	if err != nil {
		t.Fatal("Promise should not have an error")
	}
	if len(results) != 2 {
		t.Fatal("Promise should have 2 results")
	}
	if results[0] != "hello" {
		t.Fatal("Promise result should be hello")
	}
	if results[1] != "world" {
		t.Fatal("Promise result should be world")
	}
}

func TestPromiseAllReject(t *testing.T) {
	p1 := New(func(a ...any) any {
		return "hello"
	})
	p2 := New(func(a ...any) any {
		panic(errors.New("Promise rejected"))
	})
	promises := []*Promise{p1, p2}
	results, err := All(promises)
	if err == nil {
		t.Fatal("Promise should have an error")
	}
	if results != nil {
		t.Fatal("Promise results should be nil")
	}
}

func TestPromiseAllSettled(t *testing.T) {
	p1 := New(func(a ...any) any {
		return "hello"
	})
	p2 := New(func(a ...any) any {
		panic(errors.New("Promise rejected"))
	})
	p3 := New(func(a ...any) any {
		return nil
	})
	promises := []*Promise{p1, p2, p3}
	results := AllSettled(promises)
	if len(results) != 3 {
		t.Fatalf("Promise should have 3 results. Got %d", len(results))
	}
	if results[0] != "fulfilled" {
		t.Fatalf("Promise result should be fulfilled. Got %v", results[0])
	}
	if results[1] != "rejected" {
		t.Fatalf("Promise result should be rejected. Got %v", results[1])
	}
	if results[2] != "fulfilled" {
		t.Fatalf("Promise result of nil should be fulfilled. Got %v", results[2])
	}
}

func TestPromiseAny(t *testing.T) {
	p1 := New(func(a ...any) any {
		panic(errors.New("Promise rejected"))
	})
	p2 := New(func(a ...any) any {
		return "world"
	})
	promises := []*Promise{p1, p2}
	result, err := Any(promises)
	if err != nil {
		t.Fatal("Promise should not have an error")
	}
	if result != "world" {
		t.Fatalf("Promise result should be world. Got %v", result)
	}
}

func TestPromiseAnyReject(t *testing.T) {
	p1 := New(func(a ...any) any {
		panic(errors.New("Promise rejected"))
	})
	p2 := New(func(a ...any) any {
		panic(errors.New("Promise rejected"))
	})
	promises := []*Promise{p1, p2}
	_, err := Any(promises)
	if err == nil {
		t.Fatal("Promise should have an error")
	}
}

func TestPromiseAnyAllResolved(t *testing.T) {
	p1 := New(func(a ...any) any {
		time.Sleep(100 * time.Millisecond)
		return "hello"
	})
	p2 := New(func(a ...any) any {
		time.Sleep(50 * time.Millisecond)
		return "world"
	})
	promises := []*Promise{p1, p2}
	result, err := Race(promises)
	if err != nil {
		t.Fatal("Promise should not have an error")
	}
	if result != "world" {
		t.Fatal("Promise result should be world")
	}
}

func TestPromiseRace(t *testing.T) {
	p1 := New(func(a ...any) any {
		time.Sleep(100 * time.Millisecond)
		return "hello"
	})
	p2 := New(func(a ...any) any {
		time.Sleep(50 * time.Millisecond)
		return "world"
	})
	promises := []*Promise{p1, p2}
	result, err := Race(promises)
	if err != nil {
		t.Fatal("Promise should not have an error")
	}
	if result != "world" {
		t.Fatal("Promise result should be world")
	}
}

func TestPromiseRaceRejectFirst(t *testing.T) {
	p1 := New(func(a ...any) any {
		time.Sleep(50 * time.Millisecond)
		panic("Promise rejected")
	})
	p2 := New(func(a ...any) any {
		time.Sleep(150 * time.Millisecond)
		return "world"
	})
	promises := []*Promise{p1, p2}
	_, err := Race(promises)
	if err == nil {
		t.Fatal("Promise should have rejected")
	}
}

func TestPromiseRaceRejectAfterFullfill(t *testing.T) {
	p1 := New(func(a ...any) any {
		time.Sleep(150 * time.Millisecond)
		panic("Promise rejected")
	})
	p2 := New(func(a ...any) any {
		time.Sleep(50 * time.Millisecond)
		return "world"
	})
	promises := []*Promise{p1, p2}
	result, err := Race(promises)
	if err != nil {
		t.Fatal("Promise should have been fullfilled")
	}
	if result != "world" {
		t.Fatal("Promise result should be world")
	}
}
