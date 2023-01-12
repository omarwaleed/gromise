package gromise

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Promise struct {
	Fn     (func(...any) any)
	Result *any
	Error  error
}

// Creates a new promise. Returns a value if resolved or rejects if function panics. Does NOT run the promise.
func New(fn func(...any) any) *Promise {
	p := &Promise{
		Fn: fn,
	}
	return p
}

// Runs the promise given. If the promise panics, it will be rejected. If the promise is rejected, the error will be set on the promise.
func (p *Promise) Run(wg *sync.WaitGroup) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				p.Error = fmt.Errorf("Promise rejected: %v", err)
			}
			if wg != nil {
				wg.Done()
			}
		}()
		result := p.Fn()
		p.Result = &result
	}()
}

// Given a list of promises, runs them all and returns the results. Rejects if any promise panics.
func All(promises []*Promise) ([]any, error) {
	results := []any{}
	wg := sync.WaitGroup{}
	wg.Add(len(promises))
	for _, promise := range promises {
		promise.Run(&wg)
	}
	wg.Wait()
	for _, promise := range promises {
		if promise.Error != nil {
			return nil, promise.Error
		}
		results = append(results, *promise.Result)
	}
	return results, nil
}

// Given a list of promises, runs them all and returns the status of each promise fullfilled/rejected.
func AllSettled(promises []*Promise) []any {
	results := make([]any, len(promises))
	wg := sync.WaitGroup{}
	wg.Add(len(promises))
	for _, promise := range promises {
		promise.Run(&wg)
	}
	wg.Wait()
	for index, promise := range promises {
		if promise.Error != nil {
			results[index] = "rejected"
		} else {
			results[index] = "fulfilled"
		}
	}
	return results
}

// Given a list of promises, runs them all and returns the first non rejected result. Rejects if all promises reject.
func Any(promises []*Promise) (any, error) {
	for _, promise := range promises {
		promise.Run(nil)
	}
	ticker := time.NewTicker(time.Millisecond * 5)
	defer ticker.Stop()
	for {
		<-ticker.C
		allReject := true
		for _, promise := range promises {
			if promise.Result != nil {
				return *promise.Result, nil
			}
			if promise.Error == nil {
				allReject = false
			}
		}
		if allReject {
			break
		}
	}
	return nil, errors.New("All promises rejected")
}

// Given a list of promises, runs them all and returns the first non rejected result. Rejects if any promise rejects first.
func Race(promises []*Promise) (any, error) {
	for _, promise := range promises {
		promise.Run(nil)
	}
	ticker := time.NewTicker(time.Millisecond * 5)
	defer ticker.Stop()
	for {
		<-ticker.C
		for _, promise := range promises {
			if promise.Result != nil {
				return *promise.Result, nil
			}
			if promise.Error != nil {
				return nil, promise.Error
			}
		}
	}
}
