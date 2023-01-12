package gromise

import (
	"errors"
	"fmt"
	"sync"
)

type Promise struct {
	Fn     (func(...any) any)
	Result *any
	Error  error
}

type promiseArrayWithMutex struct {
	promises []*Promise
	mu       sync.Mutex
}

// Creates a new promise. Returns a value if resolved or rejects if function panics. Does NOT run the promise.
func New(fn func(...any) any) *Promise {
	p := &Promise{
		Fn: fn,
	}
	return p
}

// Runs the promise given
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
	wg := sync.WaitGroup{}
	wg.Add(len(promises))
	for _, promise := range promises {
		promise.Run(&wg)
	}
	wg.Wait()
	for _, promise := range promises {
		if promise.Error == nil {
			return promise.Result, nil
		}
	}
	return nil, errors.New("All promises rejected")
}

// Given a list of promises, runs them all and returns the first non rejected result. Rejects if all promises reject.
func Race(promises []*Promise) (any, error) {
	wg := sync.WaitGroup{}
	wg.Add(len(promises))
	pa := promiseArrayWithMutex{
		mu:       sync.Mutex{},
		promises: promises,
	}
	for _, promise := range promises {
		promise.Run(&wg)
	}
	wg.Wait()
	pa.mu.Lock()
	defer pa.mu.Unlock()
	for _, promise := range pa.promises {
		if promise.Error == nil {
			return promise.Result, nil
		}
	}
	return nil, errors.New("All promises rejected")
}
