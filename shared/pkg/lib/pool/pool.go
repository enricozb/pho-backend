package pool

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Pool interface {
	Go(func() error)
	Wait() error
}

type pool struct {
	ctx     context.Context
	g       *errgroup.Group
	sem     *semaphore.Weighted
	workers int
}

var _ Pool = &pool{}

func NewPool(workers int) *pool {
	sem := semaphore.NewWeighted(int64(workers))
	g, ctx := errgroup.WithContext(context.Background())

	return &pool{ctx, g, sem, workers}
}

func (p *pool) Go(f func() error) {
	p.g.Go(func() error {
		if err := p.sem.Acquire(p.ctx, 1); err != nil {
			return fmt.Errorf("acquire: %v", err)
		}
		defer p.sem.Release(1)

		return f()
	})
}

func (p *pool) Wait() error {
	return p.g.Wait()
}
