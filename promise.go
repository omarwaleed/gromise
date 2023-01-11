package gromise

import "sync"

type Promise struct {
	Fn     (func(...any) any)
	Result *any
	Error  error
}

func New(fn func(...any) any) *Promise {
	p := &Promise{
		Fn: fn,
	}
	return p
}

func (p *Promise) Run(wg *sync.WaitGroup) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				p.Error = err.(error)
			}
			wg.Done()
		}()
		result := p.Fn()
		p.Result = &result
	}()
}

func All(promises []*Promise) ([]any, error) {
	results := []any{}
	wg := sync.WaitGroup{}
	wg.Add(len(promises))
	for _, promise := range promises {
		promise.Run(&wg)
		if promise.Error != nil {
			return nil, promise.Error
		}
		results = append(results, *promise.Result)
	}
	wg.Wait()
	return results, nil
}

func AllSettled(promises []*Promise) []any {
	results := make([]any, len(promises))
	wg := sync.WaitGroup{}
	for _, promise := range promises {
		wg.Add(1)
		promise.Run(&wg)
		if promise.Error != nil {
			results = append(results, promise.Error)
		} else {
			results = append(results, *promise.Result)
		}
	}
	wg.Wait()
	return results
}
