package pool

import (
	"sync"
	"sync/atomic"

	"github.com/fzzy/radix/extra/pool"

	"umbrella-go/umbrella-common/redis"
)

type (
	DialFunc = pool.DialFunc
	RawPool  = pool.Pool
)

type Pool struct {
	pool *pool.Pool
	opts []redis.Option
}

func WrapPool(p *pool.Pool, opts ...redis.Option) *Pool {
	return &Pool{
		pool: p,
		opts: opts,
	}
}

func UnwrapPool(p *Pool) *pool.Pool {
	return p.pool
}

func (p *Pool) CarefullyPut(conn *redis.Client, potentialErr *error) {
	p.pool.CarefullyPut(redis.UnwrapClient(conn), potentialErr)
}

func (p *Pool) Get() (*redis.Client, error) {
	c, err := p.pool.Get()
	if err != nil {
		return nil, err
	}

	return redis.WrapClient(c, p.opts...), nil
}

func (p *Pool) Put(conn *redis.Client) {
	p.pool.Put(redis.UnwrapClient(conn))
}

func (p *Pool) Empty() {
	p.pool.Empty()
}

func NewCustomPool(network, addr string, size int, df pool.DialFunc, opts ...redis.Option) (*Pool, error) {
	p, err := pool.NewCustomPool(network, addr, size, df)
	if err != nil {
		return nil, err
	}
	return WrapPool(p, opts...), nil
}

func NewPool(network, addr string, size int, opts ...redis.Option) (*Pool, error) {
	p, err := pool.NewPool(network, addr, size)
	if err != nil {
		return nil, err
	}
	return WrapPool(p, opts...), nil
}

func NewOrEmptyPool(network, addr string, size int, opts ...redis.Option) *Pool {
	p := pool.NewOrEmptyPool(network, addr, size)
	return WrapPool(p, opts...)
}

type LazyPool struct {
	*Pool
	done uint32
	m    sync.Mutex

	newPool func() (*Pool, error)
}

func NewLazyPool(newPool func() (*Pool, error)) *LazyPool {
	return &LazyPool{
		newPool: newPool,
	}
}

func NewCustomLazyPool(network, addr string, size int, df pool.DialFunc, opts ...redis.Option) *LazyPool {
	return NewLazyPool(func() (*Pool, error) {
		return NewCustomPool(network, addr, size, df, opts...)
	})
}

func (lr *LazyPool) GetPool() (*Pool, error) {
	if atomic.LoadUint32(&lr.done) == 1 {
		return lr.Pool, nil
	}

	lr.m.Lock()
	defer lr.m.Unlock()
	if lr.done == 0 {
		p, err := lr.newPool()
		if err != nil {
			return nil, err
		}

		lr.Pool = p
		lr.newPool = nil
		atomic.StoreUint32(&lr.done, 1)
	}

	return lr.Pool, nil
}

func (lr *LazyPool) Empty() {
	if atomic.LoadUint32(&lr.done) == 1 {
		lr.Pool.Empty()
	}
}
